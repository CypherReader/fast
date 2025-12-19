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

// MockSOSRepository is a mock implementation of ports.SOSRepository
type MockSOSRepository struct {
	mock.Mock
}

func (m *MockSOSRepository) Save(ctx context.Context, sos *domain.SOSFlare) error {
	args := m.Called(ctx, sos)
	return args.Error(0)
}

func (m *MockSOSRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SOSFlare, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SOSFlare), args.Error(1)
}

func (m *MockSOSRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.SOSFlare, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SOSFlare), args.Error(1)
}

func (m *MockSOSRepository) FindAllActive(ctx context.Context) ([]*domain.SOSFlare, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SOSFlare), args.Error(1)
}

func (m *MockSOSRepository) UpdateStatus(ctx context.Context, sosID uuid.UUID, status domain.SOSStatus) error {
	args := m.Called(ctx, sosID, status)
	return args.Error(0)
}

func (m *MockSOSRepository) UpdateCortexResponse(ctx context.Context, sosID uuid.UUID) error {
	args := m.Called(ctx, sosID)
	return args.Error(0)
}

func (m *MockSOSRepository) SaveHypeResponse(ctx context.Context, hype *domain.HypeResponse) error {
	args := m.Called(ctx, hype)
	return args.Error(0)
}

func (m *MockSOSRepository) GetHypeResponses(ctx context.Context, sosID uuid.UUID) ([]domain.HypeResponse, error) {
	args := m.Called(ctx, sosID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.HypeResponse), args.Error(1)
}

func (m *MockSOSRepository) IncrementHypeCount(ctx context.Context, sosID uuid.UUID) error {
	args := m.Called(ctx, sosID)
	return args.Error(0)
}

func (m *MockSOSRepository) GetUserHypeCount(ctx context.Context, userID uuid.UUID, date time.Time) (int, error) {
	args := m.Called(ctx, userID, date)
	return args.Int(0), args.Error(1)
}

// MockTribeService is a mock of ports.TribeService
type MockTribeService struct {
	mock.Mock
}

func (m *MockTribeService) GetMyTribes(ctx context.Context, userID string) ([]domain.Tribe, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tribe), args.Error(1)
}

func (m *MockTribeService) GetTribe(ctx context.Context, tribeID string, currentUserID *string) (*domain.Tribe, error) {
	args := m.Called(ctx, tribeID, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeService) GetTribeMembers(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	args := m.Called(ctx, tribeID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TribeMember), args.Error(1)
}

func (m *MockTribeService) CreateTribe(ctx context.Context, userID string, req domain.CreateTribeRequest) (*domain.Tribe, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeService) JoinTribe(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeService) LeaveTribe(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeService) UpdateTribe(ctx context.Context, tribeID, userID string, req domain.UpdateTribeRequest) (*domain.Tribe, error) {
	args := m.Called(ctx, tribeID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeService) DeleteTribe(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeService) ListTribes(ctx context.Context, query domain.ListTribesQuery, currentUserID *string) ([]domain.Tribe, int, error) {
	args := m.Called(ctx, query, currentUserID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tribe), args.Int(1), args.Error(2)
}

func (m *MockTribeService) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	args := m.Called(ctx, tribeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TribeStats), args.Error(1)
}

// MockNotificationService is a mock of ports.NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendNotification(ctx context.Context, userID uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	args := m.Called(ctx, userID, title, body, notifType, data)
	return args.Error(0)
}

func (m *MockNotificationService) SendBatchNotification(ctx context.Context, userIDs []uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error {
	args := m.Called(ctx, userIDs, title, body, notifType, data)
	return args.Error(0)
}

func (m *MockNotificationService) RegisterFCMToken(ctx context.Context, userID uuid.UUID, token, deviceType string) error {
	args := m.Called(ctx, userID, token, deviceType)
	return args.Error(0)
}

func (m *MockNotificationService) UnregisterFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockNotificationService) GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Notification), args.Error(1)
}

// MockCortexService is a mock of ports.CortexService
type MockCortexService struct {
	mock.Mock
}

func (m *MockCortexService) Chat(ctx context.Context, userID uuid.UUID, message string) (string, error) {
	args := m.Called(ctx, userID, message)
	return args.String(0), args.Error(1)
}

func (m *MockCortexService) GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error) {
	args := m.Called(ctx, userID, fastingHours)
	return args.String(0), args.Error(1)
}

func (m *MockCortexService) AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error) {
	args := m.Called(ctx, imageBase64, description)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *MockCortexService) GetCravingHelp(ctx context.Context, userID uuid.UUID, cravingDescription string) (interface{}, error) {
	args := m.Called(ctx, userID, cravingDescription)
	return args.Get(0), args.Error(1)
}

// ============== SEND SOS TESTS ==============

func TestSOSService_SendSOSFlare_Success(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:                       userID,
		Name:                     "Test User",
		PushNotificationsEnabled: true,
	}

	activeFast := &domain.FastingSession{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    domain.StatusActive,
		StartTime: time.Now().Add(-10 * time.Hour),
	}

	// Setup expectations
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSOSRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(activeFast, nil)
	mockSOSRepo.On("Save", ctx, mock.AnythingOfType("*domain.SOSFlare")).Return(nil)
	mockCortexService.On("GetCravingHelp", ctx, userID, "I want pizza").Return(map[string]interface{}{"help": "Drink water"}, nil)
	mockTribeService.On("GetMyTribes", ctx, userID.String()).Return([]domain.Tribe{}, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	sos, aiResponse, err := service.SendSOSFlare(ctx, userID, "I want pizza")

	assert.NoError(t, err)
	assert.NotNil(t, sos)
	assert.NotNil(t, aiResponse)
	assert.Equal(t, domain.SOSStatusActive, sos.Status)
	assert.True(t, sos.HoursFasted >= 10)
}

func TestSOSService_SendSOSFlare_Cooldown(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	userID := uuid.New()

	recentTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	user := &domain.User{
		ID:                       userID,
		PushNotificationsEnabled: true,
	}

	recentSOS := &domain.SOSFlare{
		ID:        uuid.New(),
		UserID:    userID,
		CreatedAt: recentTime,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSOSRepo.On("FindActiveByUserID", ctx, userID).Return(recentSOS, nil)

	sos, aiResponse, err := service.SendSOSFlare(ctx, userID, "I want pizza")

	assert.Error(t, err)
	assert.Nil(t, sos)
	assert.Nil(t, aiResponse)
	assert.Contains(t, err.Error(), "cooldown")
}

func TestSOSService_SendSOSFlare_NoActiveFast(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID, PushNotificationsEnabled: true}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockSOSRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)
	mockFastingRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))

	sos, aiResponse, err := service.SendSOSFlare(ctx, userID, "I want pizza")

	assert.Error(t, err)
	assert.Nil(t, sos)
	assert.Nil(t, aiResponse)
	assert.Contains(t, err.Error(), "no active fast")
}

// ============== SEND HYPE TESTS ==============

func TestSOSService_SendHype_Success(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()
	fromUserID := uuid.New()
	toUserID := uuid.New()

	sos := &domain.SOSFlare{
		ID:     sosID,
		UserID: toUserID,
		Status: domain.SOSStatusActive,
	}

	sender := &domain.User{
		ID:   fromUserID,
		Name: "Helper",
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)
	mockUserRepo.On("FindByID", ctx, fromUserID).Return(sender, nil)
	mockSOSRepo.On("SaveHypeResponse", ctx, mock.AnythingOfType("*domain.HypeResponse")).Return(nil)
	mockSOSRepo.On("IncrementHypeCount", ctx, sosID).Return(nil)
	mockNotificationService.On("SendNotification", ctx, toUserID, mock.Anything, mock.Anything, domain.NotificationTypeHypeReceived, mock.Anything).Return(nil)

	err := service.SendHype(ctx, sosID, fromUserID, "💪", "You got this!")

	assert.NoError(t, err)
	mockSOSRepo.AssertExpectations(t)
}

func TestSOSService_SendHype_SOSNotActive(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()
	fromUserID := uuid.New()

	sos := &domain.SOSFlare{
		ID:     sosID,
		Status: domain.SOSStatusRescued, // Already resolved
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)

	err := service.SendHype(ctx, sosID, fromUserID, "💪", "You got this!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no longer active")
}

func TestSOSService_SendHype_SOSNotFound(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()
	fromUserID := uuid.New()

	mockSOSRepo.On("FindByID", ctx, sosID).Return(nil, errors.New("not found"))

	err := service.SendHype(ctx, sosID, fromUserID, "💪", "You got this!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// ============== RESOLVE SOS TESTS ==============

func TestSOSService_ResolveSOS_Survived(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()
	userID := uuid.New()
	helperID := uuid.New()

	sos := &domain.SOSFlare{
		ID:     sosID,
		UserID: userID,
		Status: domain.SOSStatusActive,
	}

	user := &domain.User{
		ID:   userID,
		Name: "Survivor",
	}

	hypes := []domain.HypeResponse{
		{ID: uuid.New(), FromUserID: helperID},
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)
	mockSOSRepo.On("UpdateStatus", ctx, sosID, domain.SOSStatusRescued).Return(nil)
	mockSOSRepo.On("GetHypeResponses", ctx, sosID).Return(hypes, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockNotificationService.On("SendBatchNotification", ctx, mock.Anything, mock.Anything, mock.Anything, domain.NotificationTypeSOSResolved, mock.Anything).Return(nil)

	err := service.ResolveSOS(ctx, sosID, true)

	assert.NoError(t, err)
	mockSOSRepo.AssertExpectations(t)
}

func TestSOSService_ResolveSOS_Failed(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()

	sos := &domain.SOSFlare{
		ID:     sosID,
		Status: domain.SOSStatusActive,
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)
	mockSOSRepo.On("UpdateStatus", ctx, sosID, domain.SOSStatusFailed).Return(nil)

	err := service.ResolveSOS(ctx, sosID, false)

	assert.NoError(t, err)
	// No notification should be sent for failed SOS
	mockNotificationService.AssertNotCalled(t, "SendBatchNotification")
}

// ============== CORTEX BACKUP TESTS ==============

func TestSOSService_CheckAndSendCortexBackup_NoResponseAfter10Min(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()
	userID := uuid.New()

	sos := &domain.SOSFlare{
		ID:              sosID,
		UserID:          userID,
		HoursFasted:     14.5,
		Status:          domain.SOSStatusActive,
		CortexResponded: false,
		CreatedAt:       time.Now().Add(-15 * time.Minute), // 15 minutes ago
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)
	mockSOSRepo.On("GetHypeResponses", ctx, sosID).Return([]domain.HypeResponse{}, nil) // No responses
	mockNotificationService.On("SendNotification", ctx, userID, mock.Anything, mock.Anything, domain.NotificationTypeCortexBackup, mock.Anything).Return(nil)
	mockSOSRepo.On("UpdateCortexResponse", ctx, sosID).Return(nil)

	err := service.CheckAndSendCortexBackup(ctx, sosID)

	assert.NoError(t, err)
	mockSOSRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestSOSService_CheckAndSendCortexBackup_HasHypeResponses(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()

	sos := &domain.SOSFlare{
		ID:              sosID,
		Status:          domain.SOSStatusActive,
		CortexResponded: false,
		CreatedAt:       time.Now().Add(-15 * time.Minute),
	}

	hypes := []domain.HypeResponse{
		{ID: uuid.New()},
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)
	mockSOSRepo.On("GetHypeResponses", ctx, sosID).Return(hypes, nil) // Has responses

	err := service.CheckAndSendCortexBackup(ctx, sosID)

	assert.NoError(t, err)
	// Cortex backup should NOT be sent since there are hype responses
	mockNotificationService.AssertNotCalled(t, "SendNotification")
}

func TestSOSService_CheckAndSendCortexBackup_TooEarly(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()

	sos := &domain.SOSFlare{
		ID:              sosID,
		Status:          domain.SOSStatusActive,
		CortexResponded: false,
		CreatedAt:       time.Now().Add(-5 * time.Minute), // Only 5 minutes ago
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)

	err := service.CheckAndSendCortexBackup(ctx, sosID)

	assert.NoError(t, err)
	// Should not check hypes or send notification since it's too early
	mockSOSRepo.AssertNotCalled(t, "GetHypeResponses")
	mockNotificationService.AssertNotCalled(t, "SendNotification")
}

func TestSOSService_CheckAndSendCortexBackup_AlreadyResponded(t *testing.T) {
	mockSOSRepo := new(MockSOSRepository)
	mockUserRepo := new(MockUserRepository)
	mockTribeService := new(MockTribeService)
	mockNotificationService := new(MockNotificationService)
	mockCortexService := new(MockCortexService)
	mockFastingRepo := new(MockFastingRepository)

	service := NewSOSService(
		mockSOSRepo,
		mockUserRepo,
		mockTribeService,
		mockNotificationService,
		mockCortexService,
		mockFastingRepo,
	)

	ctx := context.Background()
	sosID := uuid.New()

	sos := &domain.SOSFlare{
		ID:              sosID,
		Status:          domain.SOSStatusActive,
		CortexResponded: true, // Already responded
		CreatedAt:       time.Now().Add(-15 * time.Minute),
	}

	mockSOSRepo.On("FindByID", ctx, sosID).Return(sos, nil)

	err := service.CheckAndSendCortexBackup(ctx, sosID)

	assert.NoError(t, err)
	mockNotificationService.AssertNotCalled(t, "SendNotification")
}
