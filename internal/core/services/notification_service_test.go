package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============== MOCK IMPLEMENTATIONS ==============

// MockNotificationRepository is a mock implementation of ports.NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) SaveToken(ctx context.Context, token *domain.FCMToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FCMToken), args.Error(1)
}

func (m *MockNotificationRepository) DeleteToken(ctx context.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}

func (m *MockNotificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	args := m.Called(ctx, notificationID)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

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

// ============== NotificationRepository Mock Tests ==============

func TestMockNotificationRepository_SaveToken(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	token := &domain.FCMToken{
		ID:         uuid.New(),
		UserID:     userID,
		Token:      "fcm_token_123",
		DeviceType: "ios",
		CreatedAt:  time.Now(),
	}

	mockRepo.On("SaveToken", ctx, token).Return(nil)

	err := mockRepo.SaveToken(ctx, token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_GetUserTokens(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	tokens := []domain.FCMToken{
		{ID: uuid.New(), UserID: userID, Token: "token1", DeviceType: "ios"},
		{ID: uuid.New(), UserID: userID, Token: "token2", DeviceType: "android"},
	}

	mockRepo.On("GetUserTokens", ctx, userID).Return(tokens, nil)

	result, err := mockRepo.GetUserTokens(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_GetUserTokens_NoTokens(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetUserTokens", ctx, userID).Return([]domain.FCMToken{}, nil)

	result, err := mockRepo.GetUserTokens(ctx, userID)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_DeleteToken(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()

	mockRepo.On("DeleteToken", ctx, "fcm_token_123").Return(nil)

	err := mockRepo.DeleteToken(ctx, "fcm_token_123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_SaveNotification(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	notification := &domain.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Test Notification",
		Message:   "This is a test message",
		Type:      "fasting_complete",
		Read:      false,
		CreatedAt: time.Now(),
	}

	mockRepo.On("Save", ctx, notification).Return(nil)

	err := mockRepo.Save(ctx, notification)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_FindByUserID(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	notifications := []domain.Notification{
		{ID: uuid.New(), UserID: userID, Title: "Notif 1", Message: "Message 1"},
		{ID: uuid.New(), UserID: userID, Title: "Notif 2", Message: "Message 2"},
		{ID: uuid.New(), UserID: userID, Title: "Notif 3", Message: "Message 3"},
	}

	mockRepo.On("FindByUserID", ctx, userID, 10).Return(notifications, nil)

	result, err := mockRepo.FindByUserID(ctx, userID, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_MarkAsRead(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	notificationID := uuid.New()

	mockRepo.On("MarkAsRead", ctx, notificationID).Return(nil)

	err := mockRepo.MarkAsRead(ctx, notificationID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockNotificationRepository_MarkAllAsRead(t *testing.T) {
	mockRepo := new(MockNotificationRepository)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("MarkAllAsRead", ctx, userID).Return(nil)

	err := mockRepo.MarkAllAsRead(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============== Integration Tests with Notification Types ==============

func TestNotificationService_NotificationTypes(t *testing.T) {
	// Test that all notification types are properly defined
	types := []domain.NotificationType{
		domain.NotificationTypeFastComplete,
		domain.NotificationTypeFastStartReminder,
		domain.NotificationTypeFastEndReminder,
		domain.NotificationTypeHydrationReminder,
		domain.NotificationTypeWeeklyCheckIn,
		domain.NotificationTypeFriendInvite,
		domain.NotificationTypeSOSFlare,
		domain.NotificationTypeHypeReceived,
		domain.NotificationTypeCortexBackup,
	}

	for _, notifType := range types {
		assert.NotEmpty(t, string(notifType))
	}
}

func TestNoOpNotificationService_SendNotification_WithData(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userID := uuid.New()

	data := map[string]string{
		"fast_id": uuid.New().String(),
		"action":  "complete",
	}

	err := service.SendNotification(ctx, userID, "Fast Complete!", "You did it!", domain.NotificationTypeFastComplete, data)

	assert.NoError(t, err)
}

func TestNoOpNotificationService_SendBatchNotification_AllTypes(t *testing.T) {
	service := NewNoOpNotificationService()
	ctx := context.Background()
	userIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	testCases := []struct {
		title     string
		body      string
		notifType domain.NotificationType
	}{
		{"Fast Complete!", "Great job!", domain.NotificationTypeFastComplete},
		{"Friend Invite", "You've been invited", domain.NotificationTypeFriendInvite},
		{"SOS Alert", "Someone needs help", domain.NotificationTypeSOSFlare},
		{"Hype!", "💪 Keep going!", domain.NotificationTypeHypeReceived},
	}

	for _, tc := range testCases {
		t.Run(string(tc.notifType), func(t *testing.T) {
			err := service.SendBatchNotification(ctx, userIDs, tc.title, tc.body, tc.notifType, nil)
			assert.NoError(t, err)
		})
	}
}
