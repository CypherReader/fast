package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// SmartReminderService handles AI-powered reminder scheduling and delivery
type SmartReminderService struct {
	reminderRepo        ports.ReminderRepository
	userRepo            ports.UserRepository
	fastingRepo         ports.FastingRepository
	notificationService ports.NotificationService
	cortexService       ports.CortexService
}

// NewSmartReminderService creates a new smart reminder service
func NewSmartReminderService(
	reminderRepo ports.ReminderRepository,
	userRepo ports.UserRepository,
	fastingRepo ports.FastingRepository,
	notificationService ports.NotificationService,
	cortexService ports.CortexService,
) *SmartReminderService {
	return &SmartReminderService{
		reminderRepo:        reminderRepo,
		userRepo:            userRepo,
		fastingRepo:         fastingRepo,
		notificationService: notificationService,
		cortexService:       cortexService,
	}
}

// ScheduleFastStartReminder schedules a reminder to start fasting based on user preferences
func (s *SmartReminderService) ScheduleFastStartReminder(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.reminderRepo.GetUserSettings(ctx, userID)
	if err != nil || !settings.ReminderFastStart {
		return nil // Reminders disabled
	}

	// Calculate next scheduled time based on preferred hour
	now := time.Now()
	scheduledAt := time.Date(now.Year(), now.Month(), now.Day(), settings.PreferredFastStartHour, 0, 0, 0, now.Location())

	// If the time has already passed today, schedule for tomorrow
	if scheduledAt.Before(now) {
		scheduledAt = scheduledAt.Add(24 * time.Hour)
	}

	// Delete any existing fast start reminders for this user
	_ = s.reminderRepo.DeleteByUserAndType(ctx, userID, domain.ReminderTypeFastStart)

	reminder := &domain.ScheduledReminder{
		ID:           uuid.New(),
		UserID:       userID,
		ReminderType: domain.ReminderTypeFastStart,
		ScheduledAt:  scheduledAt,
		Sent:         false,
		Message:      "Time to start your fast! 🌙",
		CreatedAt:    time.Now(),
	}

	return s.reminderRepo.Save(ctx, reminder)
}

// ScheduleFastEndReminder schedules a reminder for when a fast is about to complete
func (s *SmartReminderService) ScheduleFastEndReminder(ctx context.Context, userID uuid.UUID, fastEndTime time.Time) error {
	settings, err := s.reminderRepo.GetUserSettings(ctx, userID)
	if err != nil || !settings.ReminderFastEnd {
		return nil // Reminders disabled
	}

	// Schedule reminder 30 minutes before fast ends
	scheduledAt := fastEndTime.Add(-30 * time.Minute)

	// Don't schedule if it's already past
	if scheduledAt.Before(time.Now()) {
		return nil
	}

	// Delete any existing fast end reminders for this user
	_ = s.reminderRepo.DeleteByUserAndType(ctx, userID, domain.ReminderTypeFastEnd)

	reminder := &domain.ScheduledReminder{
		ID:           uuid.New(),
		UserID:       userID,
		ReminderType: domain.ReminderTypeFastEnd,
		ScheduledAt:  scheduledAt,
		Sent:         false,
		Message:      "Your fast is almost complete! 🎉 Just 30 minutes to go!",
		CreatedAt:    time.Now(),
	}

	return s.reminderRepo.Save(ctx, reminder)
}

// ScheduleHydrationReminder schedules a hydration reminder during fasting
func (s *SmartReminderService) ScheduleHydrationReminder(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.reminderRepo.GetUserSettings(ctx, userID)
	if err != nil || !settings.ReminderHydration {
		return nil // Reminders disabled
	}

	// Check if user is currently fasting
	activeFast, err := s.fastingRepo.FindActiveByUserID(ctx, userID)
	if err != nil || activeFast == nil {
		return nil // Not fasting, no need for hydration reminders
	}

	// Schedule next hydration reminder
	interval := time.Duration(settings.HydrationIntervalMinutes) * time.Minute
	scheduledAt := time.Now().Add(interval)

	reminder := &domain.ScheduledReminder{
		ID:           uuid.New(),
		UserID:       userID,
		ReminderType: domain.ReminderTypeHydration,
		ScheduledAt:  scheduledAt,
		Sent:         false,
		Message:      "💧 Stay hydrated! Drink a glass of water.",
		CreatedAt:    time.Now(),
	}

	return s.reminderRepo.Save(ctx, reminder)
}

// ProcessPendingReminders processes and sends all due reminders
func (s *SmartReminderService) ProcessPendingReminders(ctx context.Context) error {
	pendingReminders, err := s.reminderRepo.FindPending(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to find pending reminders: %w", err)
	}

	for _, reminder := range pendingReminders {
		if err := s.sendReminder(ctx, &reminder); err != nil {
			log.Printf("Failed to send reminder %s: %v", reminder.ID, err)
			continue
		}

		if err := s.reminderRepo.MarkSent(ctx, reminder.ID); err != nil {
			log.Printf("Failed to mark reminder %s as sent: %v", reminder.ID, err)
		}

		// If it's a hydration reminder and user is still fasting, schedule the next one
		if reminder.ReminderType == domain.ReminderTypeHydration {
			_ = s.ScheduleHydrationReminder(ctx, reminder.UserID)
		}
	}

	return nil
}

// sendReminder sends a single reminder via push notification
func (s *SmartReminderService) sendReminder(ctx context.Context, reminder *domain.ScheduledReminder) error {
	var notifType domain.NotificationType
	var title string

	switch reminder.ReminderType {
	case domain.ReminderTypeFastStart:
		notifType = domain.NotificationTypeFastStartReminder
		title = "⏰ Time to Fast!"
	case domain.ReminderTypeFastEnd:
		notifType = domain.NotificationTypeFastEndReminder
		title = "🎉 Almost There!"
	case domain.ReminderTypeHydration:
		notifType = domain.NotificationTypeHydrationReminder
		title = "💧 Hydration Check"
	case domain.ReminderTypeWeekly:
		notifType = domain.NotificationTypeWeeklyCheckIn
		title = "📊 Weekly Summary"
	default:
		return fmt.Errorf("unknown reminder type: %s", reminder.ReminderType)
	}

	return s.notificationService.SendNotification(
		ctx,
		reminder.UserID,
		title,
		reminder.Message,
		notifType,
		map[string]string{
			"reminder_id":   reminder.ID.String(),
			"reminder_type": string(reminder.ReminderType),
		},
	)
}

// GetReminderSettings retrieves the user's reminder preferences
func (s *SmartReminderService) GetReminderSettings(ctx context.Context, userID uuid.UUID) (*domain.ReminderSettings, error) {
	settings, err := s.reminderRepo.GetUserSettings(ctx, userID)
	if err != nil {
		// Return default settings if none found
		return &domain.ReminderSettings{
			UserID:                   userID,
			ReminderFastStart:        true,
			ReminderFastEnd:          true,
			ReminderHydration:        false,
			PreferredFastStartHour:   20, // 8 PM default
			HydrationIntervalMinutes: 60,
		}, nil
	}
	return settings, nil
}

// UpdateReminderSettings updates the user's reminder preferences
func (s *SmartReminderService) UpdateReminderSettings(ctx context.Context, userID uuid.UUID, settings *domain.ReminderSettings) error {
	settings.UserID = userID
	if err := s.reminderRepo.SaveUserSettings(ctx, settings); err != nil {
		return err
	}

	// Reschedule fast start reminder if enabled
	if settings.ReminderFastStart {
		return s.ScheduleFastStartReminder(ctx, userID)
	}

	return nil
}

// AnalyzeOptimalFastingWindow uses AI to suggest the best fasting times
func (s *SmartReminderService) AnalyzeOptimalFastingWindow(ctx context.Context, userID uuid.UUID) (*domain.OptimalFastingWindow, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Get user's fasting history
	history, _ := s.fastingRepo.FindByUserID(ctx, userID)

	// Analyze patterns
	var avgStartHour float64
	var avgDuration float64
	completedFasts := 0

	for _, fast := range history {
		if fast.EndTime != nil {
			completedFasts++
			avgStartHour += float64(fast.StartTime.Hour())
			avgDuration += fast.EndTime.Sub(fast.StartTime).Hours()
		}
	}

	// Default suggestion based on popular fasting window (8 PM - 12 PM)
	suggestedStartHour := 20
	suggestedDuration := 16
	confidence := 0.5

	if completedFasts > 3 {
		avgStartHour = avgStartHour / float64(completedFasts)
		avgDuration = avgDuration / float64(completedFasts)
		suggestedStartHour = int(avgStartHour)
		suggestedDuration = int(avgDuration)
		confidence = 0.8
	}

	now := time.Now()
	suggestedStart := time.Date(now.Year(), now.Month(), now.Day(), suggestedStartHour, 0, 0, 0, now.Location())
	if suggestedStart.Before(now) {
		suggestedStart = suggestedStart.Add(24 * time.Hour)
	}
	suggestedEnd := suggestedStart.Add(time.Duration(suggestedDuration) * time.Hour)

	// Get AI reasoning
	reasoning := s.generateOptimalWindowReasoning(ctx, user, completedFasts, suggestedStartHour, suggestedDuration)

	return &domain.OptimalFastingWindow{
		SuggestedStartTime: suggestedStart,
		SuggestedEndTime:   suggestedEnd,
		SuggestedDuration:  suggestedDuration,
		Reasoning:          reasoning,
		ConfidenceScore:    confidence,
	}, nil
}

// generateOptimalWindowReasoning creates AI-powered reasoning for the suggestion
func (s *SmartReminderService) generateOptimalWindowReasoning(ctx context.Context, user *domain.User, historyCount int, startHour, duration int) string {
	if historyCount < 3 {
		return fmt.Sprintf("Based on popular fasting patterns, we recommend starting at %d:00. As you complete more fasts, we'll personalize this to your schedule.", startHour)
	}

	prompt := fmt.Sprintf(`Provide a brief (1-2 sentences) personalized explanation for why %d:00 is their optimal fasting start time with a %d hour window. Their discipline score is %.1f. Be encouraging and specific.`, startHour, duration, user.DisciplineIndex)

	response, err := s.cortexService.Chat(ctx, user.ID, prompt)
	if err != nil || response == "" {
		return fmt.Sprintf("Based on your %d successful fasts, starting at %d:00 with a %d-hour window has worked best for you.", historyCount, startHour, duration)
	}

	return response
}
