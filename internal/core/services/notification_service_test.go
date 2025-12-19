package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ============== NoOpNotificationService TESTS ==============

func TestNoOpNotificationService_SendNotification(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userID := uuid.New()

	err := service.SendNotification(ctx, userID, "Test Title", "Test Body", "test", nil)

	assert.NoError(t, err)
}

func TestNoOpNotificationService_SendBatchNotification(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userIDs := []uuid.UUID{uuid.New(), uuid.New()}

	err := service.SendBatchNotification(ctx, userIDs, "Title", "Body", "test", nil)

	assert.NoError(t, err)
}

func TestNoOpNotificationService_RegisterFCMToken(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userID := uuid.New()

	err := service.RegisterFCMToken(ctx, userID, "token123", "ios")

	assert.NoError(t, err)
}

func TestNoOpNotificationService_UnregisterFCMToken(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userID := uuid.New()

	err := service.UnregisterFCMToken(ctx, userID, "token123")

	assert.NoError(t, err)
}

func TestNoOpNotificationService_GetHistory(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userID := uuid.New()

	history, err := service.GetHistory(ctx, userID, 10)

	assert.NoError(t, err)
	assert.Empty(t, history)
}
