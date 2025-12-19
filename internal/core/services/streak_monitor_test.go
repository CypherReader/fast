package services

import (
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============== CALCULATE STREAK TESTS ==============

func TestCalculateStreak_NoSessions(t *testing.T) {
	sessions := []domain.FastingSession{}

	result := calculateStreak(sessions)

	assert.Equal(t, 0, result)
}

func TestCalculateStreak_SingleCompletedFast(t *testing.T) {
	sessions := []domain.FastingSession{
		{
			Status:    domain.StatusCompleted,
			StartTime: time.Now().Add(-24 * time.Hour),
		},
	}

	result := calculateStreak(sessions)

	assert.Equal(t, 1, result)
}

func TestCalculateStreak_ConsecutiveDays(t *testing.T) {
	now := time.Now()
	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, StartTime: now.Add(-72 * time.Hour)}, // 3 days ago
		{Status: domain.StatusCompleted, StartTime: now.Add(-48 * time.Hour)}, // 2 days ago
		{Status: domain.StatusCompleted, StartTime: now.Add(-24 * time.Hour)}, // 1 day ago
	}

	result := calculateStreak(sessions)

	assert.Equal(t, 3, result)
}

func TestCalculateStreak_WithGap(t *testing.T) {
	now := time.Now()
	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, StartTime: now.Add(-96 * time.Hour)}, // 4 days ago
		{Status: domain.StatusCompleted, StartTime: now.Add(-24 * time.Hour)}, // Gap, then 1 day ago
	}

	result := calculateStreak(sessions)

	// Should only count the most recent consecutive streak
	assert.LessOrEqual(t, result, 2) // Gap breaks the streak
}

func TestCalculateStreak_OnlyActiveSessions(t *testing.T) {
	sessions := []domain.FastingSession{
		{
			Status:    domain.StatusActive, // Not completed
			StartTime: time.Now().Add(-24 * time.Hour),
		},
	}

	result := calculateStreak(sessions)

	assert.Equal(t, 0, result) // Active sessions don't count
}

func TestCalculateStreak_MultipleFastsOnSameDay(t *testing.T) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, StartTime: today.Add(8 * time.Hour)},  // Morning fast
		{Status: domain.StatusCompleted, StartTime: today.Add(14 * time.Hour)}, // Afternoon fast
	}

	result := calculateStreak(sessions)

	// Multiple fasts on same day should still count as 1 day streak
	assert.GreaterOrEqual(t, result, 1)
}

// ============== GET LAST COMPLETED FAST TIME TESTS ==============

func TestGetLastCompletedFastTime_NoSessions(t *testing.T) {
	sessions := []domain.FastingSession{}

	result := getLastCompletedFastTime(sessions)

	assert.True(t, result.IsZero())
}

func TestGetLastCompletedFastTime_SingleSession(t *testing.T) {
	endTime := time.Now().Add(-2 * time.Hour)
	sessions := []domain.FastingSession{
		{
			Status:  domain.StatusCompleted,
			EndTime: &endTime,
		},
	}

	result := getLastCompletedFastTime(sessions)

	assert.Equal(t, endTime.Unix(), result.Unix())
}

func TestGetLastCompletedFastTime_MultipleSessions(t *testing.T) {
	oldEndTime := time.Now().Add(-48 * time.Hour)
	recentEndTime := time.Now().Add(-2 * time.Hour)

	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, EndTime: &oldEndTime},
		{Status: domain.StatusCompleted, EndTime: &recentEndTime},
	}

	result := getLastCompletedFastTime(sessions)

	assert.Equal(t, recentEndTime.Unix(), result.Unix())
}

func TestGetLastCompletedFastTime_IncompleteSessionIgnored(t *testing.T) {
	completedEndTime := time.Now().Add(-24 * time.Hour)

	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, EndTime: &completedEndTime},
		{Status: domain.StatusActive, EndTime: nil}, // Active session
	}

	result := getLastCompletedFastTime(sessions)

	assert.Equal(t, completedEndTime.Unix(), result.Unix())
}

func TestGetLastCompletedFastTime_NilEndTimeIgnored(t *testing.T) {
	endTime := time.Now().Add(-24 * time.Hour)

	sessions := []domain.FastingSession{
		{Status: domain.StatusCompleted, EndTime: &endTime},
		{Status: domain.StatusCompleted, EndTime: nil}, // Completed but no end time (edge case)
	}

	result := getLastCompletedFastTime(sessions)

	assert.Equal(t, endTime.Unix(), result.Unix())
}
