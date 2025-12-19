package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============== MOCK IMPLEMENTATIONS ==============

// MockReminderRepository is a mock implementation of ports.ReminderRepository
type MockReminderRepository struct {
	mock.Mock
}

func (m *MockReminderRepository) Save(ctx context.Context, reminder *domain.ScheduledReminder) error {
	args := m.Called(ctx, reminder)
	return args.Error(0)
}

func (m *MockReminderRepository) FindPending(ctx context.Context, before time.Time) ([]domain.ScheduledReminder, error) {
	args := m.Called(ctx, before)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ScheduledReminder), args.Error(1)
}

func (m *MockReminderRepository) MarkSent(ctx context.Context, reminderID uuid.UUID) error {
	args := m.Called(ctx, reminderID)
	return args.Error(0)
}

func (m *MockReminderRepository) DeleteByUserAndType(ctx context.Context, userID uuid.UUID, reminderType domain.ReminderType) error {
	args := m.Called(ctx, userID, reminderType)
	return args.Error(0)
}

func (m *MockReminderRepository) GetUserSettings(ctx context.Context, userID uuid.UUID) (*domain.ReminderSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReminderSettings), args.Error(1)
}

func (m *MockReminderRepository) SaveUserSettings(ctx context.Context, settings *domain.ReminderSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

// MockNotificationService is a mock implementation of ports.NotificationService
type MockNotificationServiceForReminder struct {
	mock.Mock
}

func (m *MockNotificationServiceForReminder) SendNotification(ctx context.Context, userID uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	args := m.Called(ctx, userID, title, body, notifType, data)
	return args.Error(0)
}

func (m *MockNotificationServiceForReminder) SendBatchNotification(ctx context.Context, userIDs []uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	args := m.Called(ctx, userIDs, title, body, notifType, data)
	return args.Error(0)
}

func (m *MockNotificationServiceForReminder) RegisterFCMToken(ctx context.Context, userID uuid.UUID, token, deviceType string) error {
	args := m.Called(ctx, userID, token, deviceType)
	return args.Error(0)
}

func (m *MockNotificationServiceForReminder) UnregisterFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockNotificationServiceForReminder) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

// MockCortexServiceForReminder is a mock implementation of ports.CortexService
type MockCortexServiceForReminder struct {
	mock.Mock
}

func (m *MockCortexServiceForReminder) Chat(ctx context.Context, userID uuid.UUID, message string) (string, error) {
	args := m.Called(ctx, userID, message)
	return args.String(0), args.Error(1)
}

func (m *MockCortexServiceForReminder) GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error) {
	args := m.Called(ctx, userID, fastingHours)
	return args.String(0), args.Error(1)
}

func (m *MockCortexServiceForReminder) AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error) {
	args := m.Called(ctx, imageBase64, description)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *MockCortexServiceForReminder) GetCravingHelp(ctx context.Context, userID uuid.UUID, cravingDescription string) (interface{}, error) {
	args := m.Called(ctx, userID, cravingDescription)
	return args.Get(0), args.Error(1)
}

// ============== TESTS ==============

func TestNewSmartReminderService(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)

	assert.NotNil(t, service)
}

func TestSmartReminderService_ScheduleFastStartReminder_Success(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:                 userID,
		ReminderFastStart:      true,
		PreferredFastStartHour: 20, // 8 PM
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)
	mockReminderRepo.On("DeleteByUserAndType", ctx, userID, domain.ReminderTypeFastStart).Return(nil)
	mockReminderRepo.On("Save", ctx, mock.AnythingOfType("*domain.ScheduledReminder")).Return(nil)

	err := service.ScheduleFastStartReminder(ctx, userID)

	assert.NoError(t, err)
	mockReminderRepo.AssertExpectations(t)
}

func TestSmartReminderService_ScheduleFastStartReminder_Disabled(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:            userID,
		ReminderFastStart: false, // Disabled
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)

	err := service.ScheduleFastStartReminder(ctx, userID)

	assert.NoError(t, err)
	// Save should NOT be called
	mockReminderRepo.AssertNotCalled(t, "Save")
}

func TestSmartReminderService_ScheduleFastStartReminder_SettingsError(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(nil, errors.New("not found"))

	err := service.ScheduleFastStartReminder(ctx, userID)

	// Should return nil (reminders disabled behavior)
	assert.NoError(t, err)
}

func TestSmartReminderService_ScheduleFastEndReminder_Success(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:          userID,
		ReminderFastEnd: true,
	}

	// Fast ends in 2 hours
	fastEndTime := time.Now().Add(2 * time.Hour)

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)
	mockReminderRepo.On("DeleteByUserAndType", ctx, userID, domain.ReminderTypeFastEnd).Return(nil)
	mockReminderRepo.On("Save", ctx, mock.AnythingOfType("*domain.ScheduledReminder")).Return(nil)

	err := service.ScheduleFastEndReminder(ctx, userID, fastEndTime)

	assert.NoError(t, err)
	mockReminderRepo.AssertExpectations(t)
}

func TestSmartReminderService_ScheduleFastEndReminder_Disabled(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:          userID,
		ReminderFastEnd: false, // Disabled
	}

	fastEndTime := time.Now().Add(2 * time.Hour)

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)

	err := service.ScheduleFastEndReminder(ctx, userID, fastEndTime)

	assert.NoError(t, err)
	mockReminderRepo.AssertNotCalled(t, "Save")
}

func TestSmartReminderService_ScheduleFastEndReminder_TimePassed(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:          userID,
		ReminderFastEnd: true,
	}

	// Fast ended 10 minutes ago (so 30 min before is already past)
	fastEndTime := time.Now().Add(-10 * time.Minute)

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)

	err := service.ScheduleFastEndReminder(ctx, userID, fastEndTime)

	assert.NoError(t, err)
	mockReminderRepo.AssertNotCalled(t, "Save")
}

func TestSmartReminderService_ScheduleHydrationReminder_Success(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:                   userID,
		ReminderHydration:        true,
		HydrationIntervalMinutes: 60,
	}

	activeFast := &domain.FastingSession{
		ID:     uuid.New(),
		UserID: userID,
		Status: domain.StatusActive,
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(activeFast, nil)
	mockReminderRepo.On("Save", ctx, mock.AnythingOfType("*domain.ScheduledReminder")).Return(nil)

	err := service.ScheduleHydrationReminder(ctx, userID)

	assert.NoError(t, err)
	mockReminderRepo.AssertExpectations(t)
}

func TestSmartReminderService_ScheduleHydrationReminder_Disabled(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:            userID,
		ReminderHydration: false, // Disabled
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)

	err := service.ScheduleHydrationReminder(ctx, userID)

	assert.NoError(t, err)
	mockReminderRepo.AssertNotCalled(t, "Save")
}

func TestSmartReminderService_ScheduleHydrationReminder_NotFasting(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	settings := &domain.ReminderSettings{
		UserID:            userID,
		ReminderHydration: true,
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil) // Not fasting

	err := service.ScheduleHydrationReminder(ctx, userID)

	assert.NoError(t, err)
	mockReminderRepo.AssertNotCalled(t, "Save")
}

func TestSmartReminderService_ProcessPendingReminders_Success(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()
	reminderID := uuid.New()

	pendingReminders := []domain.ScheduledReminder{
		{
			ID:           reminderID,
			UserID:       userID,
			ReminderType: domain.ReminderTypeFastStart,
			ScheduledAt:  time.Now().Add(-5 * time.Minute),
			Sent:         false,
			Message:      "Time to start your fast! 🌙",
		},
	}

	mockReminderRepo.On("FindPending", ctx, mock.AnythingOfType("time.Time")).Return(pendingReminders, nil)
	mockNotifService.On("SendNotification", ctx, userID, "⏰ Time to Fast!", "Time to start your fast! 🌙", domain.NotificationTypeFastStartReminder, mock.Anything).Return(nil)
	mockReminderRepo.On("MarkSent", ctx, reminderID).Return(nil)

	err := service.ProcessPendingReminders(ctx)

	assert.NoError(t, err)
	mockReminderRepo.AssertExpectations(t)
	mockNotifService.AssertExpectations(t)
}

func TestSmartReminderService_ProcessPendingReminders_NoPending(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()

	mockReminderRepo.On("FindPending", ctx, mock.AnythingOfType("time.Time")).Return([]domain.ScheduledReminder{}, nil)

	err := service.ProcessPendingReminders(ctx)

	assert.NoError(t, err)
	mockNotifService.AssertNotCalled(t, "SendNotification")
}

func TestSmartReminderService_ProcessPendingReminders_HydrationRescheduled(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()
	reminderID := uuid.New()

	// Hydration reminder
	pendingReminders := []domain.ScheduledReminder{
		{
			ID:           reminderID,
			UserID:       userID,
			ReminderType: domain.ReminderTypeHydration,
			ScheduledAt:  time.Now().Add(-1 * time.Minute),
			Sent:         false,
			Message:      "💧 Stay hydrated!",
		},
	}

	settings := &domain.ReminderSettings{
		UserID:                   userID,
		ReminderHydration:        true,
		HydrationIntervalMinutes: 60,
	}

	activeFast := &domain.FastingSession{
		ID:     uuid.New(),
		UserID: userID,
		Status: domain.StatusActive,
	}

	mockReminderRepo.On("FindPending", ctx, mock.AnythingOfType("time.Time")).Return(pendingReminders, nil)
	mockNotifService.On("SendNotification", ctx, userID, "💧 Hydration Check", "💧 Stay hydrated!", domain.NotificationTypeHydrationReminder, mock.Anything).Return(nil)
	mockReminderRepo.On("MarkSent", ctx, reminderID).Return(nil)
	// For rescheduling hydration reminder
	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(settings, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(activeFast, nil)
	mockReminderRepo.On("Save", ctx, mock.AnythingOfType("*domain.ScheduledReminder")).Return(nil)

	err := service.ProcessPendingReminders(ctx)

	assert.NoError(t, err)
	mockReminderRepo.AssertExpectations(t)
}

func TestSmartReminderService_GetReminderSettings_Existing(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	existingSettings := &domain.ReminderSettings{
		UserID:                   userID,
		ReminderFastStart:        true,
		ReminderFastEnd:          true,
		ReminderHydration:        true,
		PreferredFastStartHour:   19,
		HydrationIntervalMinutes: 45,
	}

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(existingSettings, nil)

	settings, err := service.GetReminderSettings(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, 19, settings.PreferredFastStartHour)
	assert.Equal(t, 45, settings.HydrationIntervalMinutes)
}

func TestSmartReminderService_GetReminderSettings_Defaults(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(nil, errors.New("not found"))

	settings, err := service.GetReminderSettings(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, settings)
	// Check defaults
	assert.Equal(t, userID, settings.UserID)
	assert.True(t, settings.ReminderFastStart)
	assert.True(t, settings.ReminderFastEnd)
	assert.False(t, settings.ReminderHydration)
	assert.Equal(t, 20, settings.PreferredFastStartHour) // Default 8 PM
	assert.Equal(t, 60, settings.HydrationIntervalMinutes)
}

func TestSmartReminderService_UpdateReminderSettings_Success(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	newSettings := &domain.ReminderSettings{
		ReminderFastStart:        true,
		ReminderFastEnd:          false,
		ReminderHydration:        true,
		PreferredFastStartHour:   21,
		HydrationIntervalMinutes: 30,
	}

	mockReminderRepo.On("SaveUserSettings", ctx, mock.AnythingOfType("*domain.ReminderSettings")).Return(nil)
	// Since ReminderFastStart is true, ScheduleFastStartReminder will be called
	existingSettings := &domain.ReminderSettings{
		UserID:                 userID,
		ReminderFastStart:      true,
		PreferredFastStartHour: 21,
	}
	mockReminderRepo.On("GetUserSettings", ctx, userID).Return(existingSettings, nil)
	mockReminderRepo.On("DeleteByUserAndType", ctx, userID, domain.ReminderTypeFastStart).Return(nil)
	mockReminderRepo.On("Save", ctx, mock.AnythingOfType("*domain.ScheduledReminder")).Return(nil)

	err := service.UpdateReminderSettings(ctx, userID, newSettings)

	assert.NoError(t, err)
	assert.Equal(t, userID, newSettings.UserID)
	mockReminderRepo.AssertExpectations(t)
}

func TestSmartReminderService_UpdateReminderSettings_DisabledFastStart(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	newSettings := &domain.ReminderSettings{
		ReminderFastStart: false, // Disabled
	}

	mockReminderRepo.On("SaveUserSettings", ctx, mock.AnythingOfType("*domain.ReminderSettings")).Return(nil)

	err := service.UpdateReminderSettings(ctx, userID, newSettings)

	assert.NoError(t, err)
	// ScheduleFastStartReminder should NOT be called when disabled
	mockReminderRepo.AssertNotCalled(t, "DeleteByUserAndType")
}

func TestSmartReminderService_AnalyzeOptimalFastingWindow_NewUser(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Name:            "Test User",
		DisciplineIndex: 0,
	}

	// New user with no fasting history
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return([]domain.FastingSession{}, nil)

	window, err := service.AnalyzeOptimalFastingWindow(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, window)
	assert.Equal(t, 16, window.SuggestedDuration) // Default 16 hours
	assert.Equal(t, 0.5, window.ConfidenceScore)  // Low confidence for new user
	assert.Contains(t, window.Reasoning, "Based on popular fasting patterns")
}

func TestSmartReminderService_AnalyzeOptimalFastingWindow_ExperiencedUser(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:              userID,
		Name:            "Experienced User",
		DisciplineIndex: 80,
	}

	// Experienced user with 5 completed fasts
	now := time.Now()
	fastingHistory := []domain.FastingSession{
		{ID: uuid.New(), UserID: userID, StartTime: time.Date(now.Year(), now.Month(), now.Day()-5, 19, 0, 0, 0, now.Location()), EndTime: ptrTime(time.Date(now.Year(), now.Month(), now.Day()-4, 13, 0, 0, 0, now.Location()))}, // 18h
		{ID: uuid.New(), UserID: userID, StartTime: time.Date(now.Year(), now.Month(), now.Day()-4, 20, 0, 0, 0, now.Location()), EndTime: ptrTime(time.Date(now.Year(), now.Month(), now.Day()-3, 12, 0, 0, 0, now.Location()))}, // 16h
		{ID: uuid.New(), UserID: userID, StartTime: time.Date(now.Year(), now.Month(), now.Day()-3, 19, 0, 0, 0, now.Location()), EndTime: ptrTime(time.Date(now.Year(), now.Month(), now.Day()-2, 11, 0, 0, 0, now.Location()))}, // 16h
		{ID: uuid.New(), UserID: userID, StartTime: time.Date(now.Year(), now.Month(), now.Day()-2, 20, 0, 0, 0, now.Location()), EndTime: ptrTime(time.Date(now.Year(), now.Month(), now.Day()-1, 14, 0, 0, 0, now.Location()))}, // 18h
		{ID: uuid.New(), UserID: userID, StartTime: time.Date(now.Year(), now.Month(), now.Day()-1, 19, 0, 0, 0, now.Location()), EndTime: ptrTime(now)},                                                                          // 16h or so
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockFastingRepo.On("FindByUserID", ctx, userID).Return(fastingHistory, nil)
	mockCortexService.On("Chat", ctx, userID, mock.Anything).Return("Your 7 PM start time aligns perfectly with your body's natural hunger patterns!", nil)

	window, err := service.AnalyzeOptimalFastingWindow(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, window)
	assert.Equal(t, 0.8, window.ConfidenceScore) // Higher confidence
	assert.Contains(t, window.Reasoning, "7 PM")
}

func TestSmartReminderService_AnalyzeOptimalFastingWindow_UserNotFound(t *testing.T) {
	mockReminderRepo := new(MockReminderRepository)
	mockUserRepo := new(MockUserRepository)
	mockFastingRepo := new(MockFastingRepository)
	mockNotifService := new(MockNotificationServiceForReminder)
	mockCortexService := new(MockCortexServiceForReminder)

	service := NewSmartReminderService(mockReminderRepo, mockUserRepo, mockFastingRepo, mockNotifService, mockCortexService)
	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found"))

	window, err := service.AnalyzeOptimalFastingWindow(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, window)
	assert.Contains(t, err.Error(), "failed to fetch user")
}

// Helper function
func ptrTime(t time.Time) *time.Time {
	return &t
}
