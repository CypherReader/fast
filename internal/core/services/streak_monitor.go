package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StreakMonitor struct {
	userRepo     ports.UserRepository
	fastingRepo  ports.FastingRepository
	cortex       ports.CortexService
	notifService ports.NotificationService
}

func NewStreakMonitor(userRepo ports.UserRepository, fastingRepo ports.FastingRepository, cortex ports.CortexService, notifService ports.NotificationService) *StreakMonitor {
	return &StreakMonitor{
		userRepo:     userRepo,
		fastingRepo:  fastingRepo,
		cortex:       cortex,
		notifService: notifService,
	}
}

// StreakRiskResponse contains streak status and intervention
type StreakRiskResponse struct {
	IsAtRisk          bool    `json:"is_at_risk"`
	CurrentStreak     int     `json:"current_streak"`
	DaysSinceLastFast int     `json:"days_since_last_fast"`
	HoursUntilLoss    float64 `json:"hours_until_loss"`
	UrgencyLevel      string  `json:"urgency_level"` // "none", "warning", "critical"
	AIMessage         string  `json:"ai_message"`
	SuggestedAction   string  `json:"suggested_action"`
	MotivationalFact  string  `json:"motivational_fact"`
}

// CheckStreakRisk monitors if user's streak is at risk
func (s *StreakMonitor) CheckStreakRisk(ctx context.Context, userID uuid.UUID) (*StreakRiskResponse, error) {
	// 1. Get user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Get fasting history
	sessions, err := s.fastingRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}

	// 3. Calculate current streak and last fast time
	currentStreak := calculateStreak(sessions)
	lastFastTime := getLastCompletedFastTime(sessions)

	response := &StreakRiskResponse{
		CurrentStreak: currentStreak,
		UrgencyLevel:  "none",
	}

	if lastFastTime.IsZero() {
		// No fasts yet
		response.SuggestedAction = "Start your first fast to begin your streak!"
		return response, nil
	}

	// 4. Calculate time since last fast
	daysSinceLastFast := int(time.Since(lastFastTime).Hours() / 24)
	response.DaysSinceLastFast = daysSinceLastFast

	// 5. Determine risk level
	// Streak breaks if no fast in 48 hours (2 days)
	hoursUntilLoss := 48.0 - time.Since(lastFastTime).Hours()
	response.HoursUntilLoss = hoursUntilLoss

	if hoursUntilLoss <= 0 {
		// Streak already lost
		response.IsAtRisk = false
		response.UrgencyLevel = "none"
		response.SuggestedAction = "Your streak has reset. Start fresh today!"
		return response, nil
	}

	if hoursUntilLoss <= 4 {
		// CRITICAL - 4 hours or less
		response.IsAtRisk = true
		response.UrgencyLevel = "critical"
	} else if hoursUntilLoss <= 12 {
		// WARNING - 12 hours or less
		response.IsAtRisk = true
		response.UrgencyLevel = "warning"
	}

	// 6. If at risk, generate AI intervention
	if response.IsAtRisk {
		aiMessage, suggestedAction, motivationalFact := s.generateStreakIntervention(ctx, user, currentStreak, hoursUntilLoss)
		response.AIMessage = aiMessage
		response.SuggestedAction = suggestedAction
		response.MotivationalFact = motivationalFact
	}

	return response, nil
}

// generateStreakIntervention creates personalized streak protection message
func (s *StreakMonitor) generateStreakIntervention(ctx context.Context, user *domain.User, streak int, hoursLeft float64) (string, string, string) {
	prompt := fmt.Sprintf(`URGENT: User's %d-day fasting streak is at risk! Only %.1f hours left.
Discipline: %.1f/100

Create emergency intervention:
1. Urgent motivational message (20 words max)
2. Immediate action they should take
3. One powerful fact about why their streak matters

Be direct, urgent, and motivating. This is critical.`, streak, hoursLeft, user.DisciplineIndex)

	systemPrompt := "You are an emergency streak coach. Your job is to save this user's streak. Be urgent but supportive."

	aiResponse, err := s.cortex.Chat(ctx, user.ID, systemPrompt+"\n\n"+prompt)

	// Fallback if AI fails
	aiMessage := fmt.Sprintf("Your %d-day streak is in danger! Don't let %d days of discipline vanish.", streak, streak)
	suggestedAction := "Start an emergency fast RIGHT NOW"
	motivationalFact := fmt.Sprintf("You've fasted %d times successfully. You've got this!", streak)

	if err == nil && aiResponse != "" {
		// Parse AI response (simplified for MVP)
		aiMessage = aiResponse
	}

	return aiMessage, suggestedAction, motivationalFact
}

// TriggerProactiveAlert sends notification when streak is at risk
func (s *StreakMonitor) TriggerProactiveAlert(ctx context.Context, userID uuid.UUID) error {
	riskStatus, err := s.CheckStreakRisk(ctx, userID)
	if err != nil {
		return err
	}

	if !riskStatus.IsAtRisk {
		return nil // No alert needed
	}

	// Send push notification
	title := "🔥 Streak Alert!"
	body := riskStatus.AIMessage
	if body == "" {
		body = fmt.Sprintf("Your %d-day streak expires in %.0f hours. Start a fast now!",
			riskStatus.CurrentStreak, riskStatus.HoursUntilLoss)
	}

	return s.notifService.SendNotification(
		ctx,
		userID,
		title,
		body,
		"streak_alert", // generic type since NotificationTypeStreak doesn't exist yet
		map[string]string{
			"urgency":    riskStatus.UrgencyLevel,
			"hours_left": fmt.Sprintf("%.1f", riskStatus.HoursUntilLoss),
		},
	)
}

// Helper functions

func calculateStreak(sessions []domain.FastingSession) int {
	if len(sessions) == 0 {
		return 0
	}

	// Sort sessions by date (most recent first)
	// Count consecutive days with completed fasts
	streak := 0
	lastDate := time.Time{}

	for i := len(sessions) - 1; i >= 0; i-- {
		session := sessions[i]
		if session.Status != domain.StatusCompleted {
			continue
		}

		currentDate := session.StartTime.Truncate(24 * time.Hour)

		if lastDate.IsZero() {
			// First fast
			lastDate = currentDate
			streak = 1
		} else {
			daysDiff := int(lastDate.Sub(currentDate).Hours() / 24)
			if daysDiff == 1 {
				// Consecutive day
				streak++
				lastDate = currentDate
			} else if daysDiff > 1 {
				// Gap in streak, break
				break
			}
			// Same day fasts don't break streak but don't increment
		}
	}

	return streak
}

func getLastCompletedFastTime(sessions []domain.FastingSession) time.Time {
	var latestTime time.Time

	for _, session := range sessions {
		if session.Status == domain.StatusCompleted && session.EndTime != nil {
			if session.EndTime.After(latestTime) {
				latestTime = *session.EndTime
			}
		}
	}

	return latestTime
}
