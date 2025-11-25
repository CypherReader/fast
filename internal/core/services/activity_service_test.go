package services

import (
	"context"
	"fastinghero/internal/adapters/repository/memory"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestActivityService_SyncAndGet(t *testing.T) {
	repo := memory.NewActivityRepository()
	service := NewActivityService(repo)
	ctx := context.Background()
	userID := uuid.New()

	activity := domain.Activity{
		ID:        "act-123",
		Type:      domain.ActivityTypeWalk,
		StartTime: time.Now(),
		Steps:     5000,
		Distance:  3.5,
		Source:    "GARMIN",
	}

	// Test Sync
	err := service.SyncActivity(ctx, userID, activity)
	assert.NoError(t, err)

	// Test GetActivities
	activities, err := service.GetActivities(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, activities, 1)
	assert.Equal(t, 5000, activities[0].Steps)
	assert.Equal(t, "GARMIN", activities[0].Source)

	// Test GetActivity
	fetched, err := service.GetActivity(ctx, "act-123")
	assert.NoError(t, err)
	assert.Equal(t, "act-123", fetched.ID)
}
