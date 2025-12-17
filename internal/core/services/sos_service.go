package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type SOSService struct {
	sosRepo         ports.SOSRepository
	userRepo        ports.UserRepository
	tribeService    ports.TribeService
	notificationSvc ports.NotificationService
	cortexService   ports.CortexService
	fastingRepo     ports.FastingRepository
}

func NewSOSService(
	sosRepo ports.SOSRepository,
	userRepo ports.UserRepository,
	tribeService ports.TribeService,
	notificationSvc ports.NotificationService,
	cortexService ports.CortexService,
	fastingRepo ports.FastingRepository,
) *SOSService {
	return &SOSService{
		sosRepo:         sosRepo,
		userRepo:        userRepo,
		tribeService:    tribeService,
		notificationSvc: notificationSvc,
		cortexService:   cortexService,
		fastingRepo:     fastingRepo,
	}
}

// SendSOSFlare triggers both AI help and tribe notifications
func (s *SOSService) SendSOSFlare(ctx context.Context, userID uuid.UUID, cravingDescription string) (*domain.SOSFlare, interface{}, error) {
	// 1. Get SOS settings
	settings, err := s.GetSOSSettings(ctx, userID)
	if err != nil {
		settings = &domain.SOSSettings{
			NotifyTribeOnSOS: true,
			AnonymousMode:    false,
		}
	}

	// 2. Rate Limit Check (24-hour cooldown)
	if settings.LastSOSAt != nil {
		timeSince := time.Since(*settings.LastSOSAt)
		if timeSince < 24*time.Hour {
			remaining := 24*time.Hour - timeSince
			return nil, nil, fmt.Errorf("cooldown active: %v remaining", remaining.Round(time.Minute))
		}
	}

	// 3. Get current fast
	activeFast, err := s.fastingRepo.FindActiveByUserID(ctx, userID)
	if err != nil || activeFast == nil {
		return nil, nil, fmt.Errorf("no active fast found")
	}

	hoursFasted := time.Since(activeFast.StartTime).Hours()

	// 4. Create SOS record
	sos := &domain.SOSFlare{
		ID:              uuid.New(),
		UserID:          userID,
		FastingID:       activeFast.ID,
		Description:     cravingDescription,
		HoursFasted:     hoursFasted,
		Status:          domain.SOSStatusActive,
		IsAnonymous:     settings.AnonymousMode,
		CortexResponded: false,
		CreatedAt:       time.Now(),
	}

	if err := s.sosRepo.Save(ctx, sos); err != nil {
		return nil, nil, fmt.Errorf("failed to save SOS: %w", err)
	}

	// 5. Get AI Response (existing Cortex service)
	aiResponse, err := s.cortexService.GetCravingHelp(ctx, userID, cravingDescription)
	if err != nil {
		log.Printf("Cortex error for SOS %s: %v", sos.ID, err)
		aiResponse = nil // Continue even if AI fails
	}

	// 6. Get user's tribes and broadcast
	tribes, err := s.tribeService.GetMyTribes(ctx, userID.String())
	if err == nil && len(tribes) > 0 && settings.NotifyTribeOnSOS {
		// Fan out asynchronously
		go s.broadcastToTribes(context.Background(), sos, tribes, userID, settings.AnonymousMode)
	}

	// 7. Update last SOS timestamp
	now := time.Now()
	settings.LastSOSAt = &now
	_ = s.UpdateSOSSettings(ctx, userID, settings)

	return sos, aiResponse, nil
}

// broadcastToTribes sends SOS notifications to all tribe members
func (s *SOSService) broadcastToTribes(ctx context.Context, sos *domain.SOSFlare, tribes []domain.Tribe, requesterID uuid.UUID, isAnonymous bool) {
	// Get user info
	user, err := s.userRepo.FindByID(ctx, requesterID)
	if err != nil {
		return
	}

	userName := user.Name
	if userName == "" {
		userName = "A tribe member"
	}

	// Use anonymous name if setting enabled
	displayName := userName
	if isAnonymous {
		displayName = "A tribe member"
	}

	for _, tribe := range tribes {
		// Get tribe members
		members, err := s.tribeService.GetTribeMembers(ctx, tribe.ID, 1000, 0)
		if err != nil {
			continue
		}

		// Filter out the requester and collect member IDs
		var memberIDs []uuid.UUID
		for _, member := range members {
			memberUUID, err := uuid.Parse(member.UserID)
			if err == nil && memberUUID != requesterID {
				memberIDs = append(memberIDs, memberUUID)
			}
		}

		if len(memberIDs) == 0 {
			continue
		}

		// Send batch notification
		title := fmt.Sprintf("🆘 Tribe Alert: %s needs backup!", displayName)
		body := fmt.Sprintf("Struggling at Hour %.1f. Send Hype to keep them in the fight.", sos.HoursFasted)

		data := map[string]string{
			"sos_id":      sos.ID.String(),
			"user_id":     requesterID.String(),
			"action_type": "send_hype",
			"deep_link":   fmt.Sprintf("app://fastinghero/sos/%s", sos.ID.String()),
		}

		err = s.notificationSvc.SendBatchNotification(
			ctx,
			memberIDs,
			title,
			body,
			domain.NotificationTypeSOSFlare,
			data,
		)
		if err != nil {
			log.Printf("Failed to send SOS notification for tribe %s: %v", tribe.ID, err)
		}
	}
}

// SendHype allows a tribe member to respond with encouragement
func (s *SOSService) SendHype(ctx context.Context, sosID, fromUserID uuid.UUID, emoji, message string) error {
	// 1. Get SOS
	sos, err := s.sosRepo.FindByID(ctx, sosID)
	if err != nil {
		return fmt.Errorf("SOS not found: %w", err)
	}

	if sos.Status != domain.SOSStatusActive {
		return fmt.Errorf("SOS is no longer active")
	}

	// 2. Check daily hype limit based on tribe size
	// First, get the tribe to check member count
	if sos.TribeID != nil {
		tribe, err := s.tribeService.GetTribe(ctx, sos.TribeID.String(), nil)
		if err == nil {
			// Get user's hype count for today
			today := time.Now().Truncate(24 * time.Hour)
			hypeCount, err := s.sosRepo.GetUserHypeCount(ctx, fromUserID, today)
			if err == nil {
				limit := domain.GetHypeLimit(tribe.MemberCount)
				if hypeCount >= limit {
					return fmt.Errorf("daily hype limit reached (%d/%d)", hypeCount, limit)
				}
			}
		}
	}

	// 3. Get sender info
	sender, err := s.userRepo.FindByID(ctx, fromUserID)
	senderName := "A friend"
	if err == nil && sender != nil && sender.Name != "" {
		senderName = sender.Name
	}

	// 4. Save hype response
	hype := &domain.HypeResponse{
		ID:         uuid.New(),
		SOSID:      sosID,
		FromUserID: fromUserID,
		FromName:   senderName,
		Message:    message,
		Emoji:      emoji,
		CreatedAt:  time.Now(),
	}

	if err := s.sosRepo.SaveHypeResponse(ctx, hype); err != nil {
		return fmt.Errorf("failed to save hype: %w", err)
	}

	// 5. Increment hype count
	if err := s.sosRepo.IncrementHypeCount(ctx, sosID); err != nil {
		return fmt.Errorf("failed to increment hype count: %w", err)
	}

	// 6. Notify the struggling user
	title := fmt.Sprintf("%s sent reinforcements! %s", senderName, emoji)
	body := message
	if body == "" {
		body = "You've got this! Keep pushing!"
	}

	data := map[string]string{
		"sos_id":       sosID.String(),
		"from_user_id": fromUserID.String(),
		"type":         "hype_received",
	}

	err = s.notificationSvc.SendNotification(
		ctx,
		sos.UserID,
		title,
		body,
		domain.NotificationTypeHypeReceived,
		data,
	)
	if err != nil {
		log.Printf("Failed to send hype notification: %v", err)
	}

	return nil
}

// ResolveSOS resolves an SOS as rescued or failed
func (s *SOSService) ResolveSOS(ctx context.Context, sosID uuid.UUID, survived bool) error {
	// 1. Get SOS
	sos, err := s.sosRepo.FindByID(ctx, sosID)
	if err != nil {
		return fmt.Errorf("SOS not found: %w", err)
	}

	// 2. Update status
	status := domain.SOSStatusFailed
	if survived {
		status = domain.SOSStatusRescued
	}

	if err := s.sosRepo.UpdateStatus(ctx, sosID, status); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// 3. Only send "rescued" notification if user survived
	if !survived {
		return nil
	}

	// 4. Get all hype responders
	hypes, err := s.sosRepo.GetHypeResponses(ctx, sosID)
	if err != nil || len(hypes) == 0 {
		return nil // No one to thank
	}

	// Get user info
	user, err := s.userRepo.FindByID(ctx, sos.UserID)
	userName := "A friend"
	if err == nil && user != nil && user.Name != "" {
		userName = user.Name
	}

	// Send thank you to all supporters
	var responderIDs []uuid.UUID
	for _, hype := range hypes {
		responderIDs = append(responderIDs, hype.FromUserID)
	}

	title := fmt.Sprintf("🎉 %s held the line!", userName)
	body := fmt.Sprintf("False alarm. %s survived the urge. Thanks for the support!", userName)

	err = s.notificationSvc.SendBatchNotification(
		ctx,
		responderIDs,
		title,
		body,
		domain.NotificationTypeSOSResolved,
		nil,
	)
	if err != nil {
		log.Printf("Failed to send rescue notification: %v", err)
	}

	return nil
}

// CheckAndSendCortexBackup checks if an SOS needs Cortex auto-response (after 10 minutes)
// This should be called by a cron job every minute
func (s *SOSService) CheckAndSendCortexBackup(ctx context.Context, sosID uuid.UUID) error {
	// 1. Get SOS
	sos, err := s.sosRepo.FindByID(ctx, sosID)
	if err != nil {
		return err
	}

	// 2. Check if already responded or resolved
	if sos.CortexResponded || sos.Status != domain.SOSStatusActive {
		return nil
	}

	// 3. Check if 10 minutes have passed
	timeSince := time.Since(sos.CreatedAt)
	if timeSince < 10*time.Minute {
		return nil // Too early
	}

	// 4. Get hype responses
	hypes, err := s.sosRepo.GetHypeResponses(ctx, sosID)
	if err != nil {
		return err
	}

	// 5. If anyone responded, don't send Cortex backup
	if len(hypes) > 0 {
		return nil
	}

	// 6. Generate Cortex backup message
	message := fmt.Sprintf("I see you're at %.1f hours and struggling. The tribe hasn't responded yet, but I'm here. "+
		"This craving will pass in 15-20 minutes. Drink water, move your body, or call someone. You're stronger than this moment.", sos.HoursFasted)

	// 7. Send notification
	title := "💪 Cortex Backup"
	data := map[string]string{
		"sos_id": sosID.String(),
		"type":   "cortex_backup",
	}

	err = s.notificationSvc.SendNotification(
		ctx,
		sos.UserID,
		title,
		message,
		domain.NotificationTypeCortexBackup,
		data,
	)
	if err != nil {
		return err
	}

	// 8. Mark as responded
	return s.sosRepo.UpdateCortexResponse(ctx, sosID)
}

// GetSOSSettings retrieves user's SOS settings
func (s *SOSService) GetSOSSettings(ctx context.Context, userID uuid.UUID) (*domain.SOSSettings, error) {
	// Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// For now, we'll store settings in the user object
	// In a real implementation, you might have a separate settings table
	settings := &domain.SOSSettings{
		NotifyTribeOnSOS: user.PushNotificationsEnabled, // Use existing field
		AnonymousMode:    false,                         // Default for now
		LastSOSAt:        nil,
	}

	// Try to get last SOS from active SOS
	activeSOS, _ := s.sosRepo.FindActiveByUserID(ctx, userID)
	if activeSOS != nil {
		settings.LastSOSAt = &activeSOS.CreatedAt
	}

	return settings, nil
}

// UpdateSOSSettings updates user's SOS settings
func (s *SOSService) UpdateSOSSettings(ctx context.Context, userID uuid.UUID, settings *domain.SOSSettings) error {
	// For now, we store minimal settings in memory
	// In production, you'd update a settings table
	// This is a placeholder that updates user's notification preference
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.PushNotificationsEnabled = settings.NotifyTribeOnSOS
	return s.userRepo.Save(ctx, user)
}
