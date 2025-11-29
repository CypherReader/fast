package memory

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetWeeklyStats_WeightAggregation(t *testing.T) {
	repo := NewTelemetryRepository()
	ctx := context.Background()
	userID := uuid.New()

	// 1. Add multiple weight entries for the same day
	now := time.Now()

	// Entry 1: Morning weight
	err := repo.SaveData(ctx, &domain.TelemetryData{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domain.MetricWeight,
		Value:     150.0,
		Timestamp: now.Add(-2 * time.Hour),
	})
	assert.NoError(t, err)

	// Entry 2: Evening weight (correction)
	err = repo.SaveData(ctx, &domain.TelemetryData{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domain.MetricWeight,
		Value:     152.0,
		Timestamp: now,
	})
	assert.NoError(t, err)

	// 2. Get Weekly Stats
	stats, err := repo.GetWeeklyStats(ctx, userID, domain.MetricWeight)
	assert.NoError(t, err)

	// 3. Verify that the value for today is 152.0 (latest), NOT 302.0 (sum)
	todayStr := now.Format("2006-01-02")
	var todayStat domain.DailyStat
	found := false
	for _, s := range stats {
		if s.Date == todayStr {
			todayStat = s
			found = true
			break
		}
	}

	assert.True(t, found, "Should have stats for today")
	assert.Equal(t, 152.0, todayStat.Value, "Weight should be the latest value, not sum")
}

func TestGetWeeklyStats_StepsAggregation(t *testing.T) {
	repo := NewTelemetryRepository()
	ctx := context.Background()
	userID := uuid.New()

	// 1. Add multiple step entries for the same day
	now := time.Now()

	// Entry 1: Morning walk
	err := repo.SaveData(ctx, &domain.TelemetryData{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domain.MetricSteps,
		Value:     5000,
		Timestamp: now.Add(-2 * time.Hour),
	})
	assert.NoError(t, err)

	// Entry 2: Evening walk
	err = repo.SaveData(ctx, &domain.TelemetryData{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      domain.MetricSteps,
		Value:     3000,
		Timestamp: now,
	})
	assert.NoError(t, err)

	// 2. Get Weekly Stats
	stats, err := repo.GetWeeklyStats(ctx, userID, domain.MetricSteps)
	assert.NoError(t, err)

	// 3. Verify that the value for today is 8000 (sum)
	todayStr := now.Format("2006-01-02")
	var todayStat domain.DailyStat
	found := false
	for _, s := range stats {
		if s.Date == todayStr {
			todayStat = s
			found = true
			break
		}
	}

	assert.True(t, found, "Should have stats for today")
	assert.Equal(t, 8000.0, todayStat.Value, "Steps should be summed")
}
