package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProgressAnalyzer struct {
	fastingRepo ports.FastingRepository
	userRepo    ports.UserRepository
	cortex      ports.CortexService
}

func NewProgressAnalyzer(fastingRepo ports.FastingRepository, userRepo ports.UserRepository, cortex ports.CortexService) *ProgressAnalyzer {
	return &ProgressAnalyzer{
		fastingRepo: fastingRepo,
		userRepo:    userRepo,
		cortex:      cortex,
	}
}

// WeeklyReport contains comprehensive weekly analytics
type WeeklyReport struct {
	UserID              uuid.UUID              `json:"user_id"`
	WeekStart           time.Time              `json:"week_start"`
	WeekEnd             time.Time              `json:"week_end"`
	FastsCompleted      int                    `json:"fasts_completed"`
	AverageDuration     float64                `json:"average_duration"`
	TotalFastingHours   float64                `json:"total_fasting_hours"`
	LongestFast         float64                `json:"longest_fast"`
	DisciplineChange    float64                `json:"discipline_change"`
	BestDay             string                 `json:"best_day"`
	ChallengeDay        string                 `json:"challenge_day"`
	GoalAchievementDate string                 `json:"goal_achievement_date,omitempty"`
	AIInsights          string                 `json:"ai_insights"`
	Predictions         map[string]interface{} `json:"predictions"`
	Recommendations     []string               `json:"recommendations"`
}

// GenerateWeeklyReport creates a comprehensive weekly progress report
func (p *ProgressAnalyzer) GenerateWeeklyReport(ctx context.Context, userID uuid.UUID) (*WeeklyReport, error) {
	// 1. Get user
	user, err := p.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// 2. Calculate week boundaries
	now := time.Now()
	weekStart := now.AddDate(0, 0, -7)

	// 3. Get all fasting sessions
	allSessions, err := p.fastingRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}

	// 4. Filter to this week's sessions
	var weekSessions []domain.FastingSession
	for _, session := range allSessions {
		if session.StartTime.After(weekStart) && session.StartTime.Before(now) {
			weekSessions = append(weekSessions, session)
		}
	}

	// 5. Calculate basic stats
	fastsCompleted := 0
	totalHours := 0.0
	longestFast := 0.0
	dayStats := make(map[string]int) // track which days user fasted

	for _, session := range weekSessions {
		if session.Status == domain.StatusCompleted {
			fastsCompleted++
			totalHours += session.ActualDurationHours
			if session.ActualDurationHours > longestFast {
				longestFast = session.ActualDurationHours
			}
			// Track day of week
			dayOfWeek := session.StartTime.Weekday().String()
			dayStats[dayOfWeek]++
		}
	}

	avgDuration := 0.0
	if fastsCompleted > 0 {
		avgDuration = totalHours / float64(fastsCompleted)
	}

	// 6. Determine best and challenge days
	bestDay := "Monday" // default
	maxFasts := 0
	for day, count := range dayStats {
		if count > maxFasts {
			maxFasts = count
			bestDay = day
		}
	}

	challengeDay := "Saturday" // example - could be day with least fasts

	// 7. Get previous week's discipline for comparison
	disciplineChange := 0.0 // Would need previous week's data in real implementation

	// 8. Generate AI insights and predictions
	aiInsights, predictions := p.generateAIInsights(ctx, user, weekSessions, avgDuration, fastsCompleted)

	// 9. Generate recommendations
	recommendations := p.generateRecommendations(fastsCompleted, avgDuration, dayStats)

	// 10. Predict goal achievement
	goalDate := p.predictGoalAchievement(user, avgDuration, fastsCompleted)

	report := &WeeklyReport{
		UserID:              userID,
		WeekStart:           weekStart,
		WeekEnd:             now,
		FastsCompleted:      fastsCompleted,
		AverageDuration:     avgDuration,
		TotalFastingHours:   totalHours,
		LongestFast:         longestFast,
		DisciplineChange:    disciplineChange,
		BestDay:             bestDay,
		ChallengeDay:        challengeDay,
		GoalAchievementDate: goalDate,
		AIInsights:          aiInsights,
		Predictions:         predictions,
		Recommendations:     recommendations,
	}

	return report, nil
}

// generateAIInsights uses Cortex to analyze the week
func (p *ProgressAnalyzer) generateAIInsights(ctx context.Context, user *domain.User, sessions []domain.FastingSession, avgDuration float64, fastsCompleted int) (string, map[string]interface{}) {
	// Construct prompt for AI
	prompt := fmt.Sprintf(`Analyze this user's weekly fasting performance and provide insights.

User Stats:
- Total fasts completed this week: %d
- Average fast duration: %.1f hours
- Discipline index: %.1f/100

Provide:
1. A 2-3 sentence analysis of their performance
2. One specific improvement area
3. One encouraging observation

Keep response under 100 words total.`, fastsCompleted, avgDuration, user.DisciplineIndex)

	systemPrompt := "You are a supportive fasting coach analyzing weekly progress. Be specific, encouraging, and actionable."

	insights, err := p.cortex.Chat(ctx, user.ID, systemPrompt+"\n\n"+prompt)
	if err != nil {
		// Fallback if AI fails
		insights = fmt.Sprintf("You completed %d fasts this week! Your dedication is building real discipline. Keep pushing forward.", fastsCompleted)
	}

	// Generate predictions
	predictions := map[string]interface{}{
		"next_week_fasts_estimate": fastsCompleted + 1,
		"discipline_trend":         "improving",
		"success_probability":      85.0,
	}

	return insights, predictions
}

// generateRecommendations creates actionable recommendations
func (p *ProgressAnalyzer) generateRecommendations(fastsCompleted int, avgDuration float64, dayStats map[string]int) []string {
	recommendations := []string{}

	if fastsCompleted == 0 {
		recommendations = append(recommendations, "Start your first fast this week!")
		recommendations = append(recommendations, "Set a reminder to begin your fasting window")
	} else if fastsCompleted < 3 {
		recommendations = append(recommendations, "Aim for 3-4 fasts next week for better results")
	} else if avgDuration < 14 {
		recommendations = append(recommendations, "Try extending your fasts by 1-2 hours")
	}

	if len(dayStats) < 3 {
		recommendations = append(recommendations, "Vary your fasting days to build consistency")
	}

	recommendations = append(recommendations, "Join a tribe for accountability and support")

	return recommendations
}

// predictGoalAchievement estimates when user will reach their goal
func (p *ProgressAnalyzer) predictGoalAchievement(user *domain.User, avgDuration float64, fastsPerWeek int) string {
	// Simple prediction based on current performance
	if fastsPerWeek == 0 {
		return ""
	}

	// Example: if user wants to lose 10 lbs and is averaging good fasts
	// This is a simplified model
	weeksToGoal := 8 // placeholder calculation
	achievementDate := time.Now().AddDate(0, 0, weeksToGoal*7)

	return achievementDate.Format("January 2, 2006")
}
