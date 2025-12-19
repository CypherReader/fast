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

// MockFastingRepository is a mock implementation of ports.FastingRepository
type MockFastingRepository struct {
	mock.Mock
}

func (m *MockFastingRepository) Save(ctx context.Context, session *domain.FastingSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockFastingRepository) Update(ctx context.Context, session *domain.FastingSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockFastingRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingSession), args.Error(1)
}

func (m *MockFastingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FastingSession), args.Error(1)
}

func (m *MockFastingRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.FastingSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FastingSession), args.Error(1)
}

// MockVaultService is a mock of ports.VaultService
type MockVaultService struct {
	mock.Mock
}

func (m *MockVaultService) GetCurrentParticipation(ctx context.Context, userID uuid.UUID) (*domain.VaultParticipation, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.VaultParticipation), args.Error(1)
}

func (m *MockVaultService) UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast, verifiedKetosis bool) {
	m.Called(ctx, user, completedFast, verifiedKetosis)
}

func (m *MockVaultService) CalculatePrice(ctx context.Context, user *domain.User) float64 {
	args := m.Called(ctx, user)
	return args.Get(0).(float64)
}

func (m *MockVaultService) ProcessDailyEarnings(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockVaultService) AddDailyEarnings(ctx context.Context, user *domain.User, amount float64) {
	m.Called(ctx, user, amount)
}

func (m *MockVaultService) CalculateVaultStatus(user *domain.User) (deposit float64, earned float64, potentialRefund float64) {
	args := m.Called(user)
	return args.Get(0).(float64), args.Get(1).(float64), args.Get(2).(float64)
}

func (m *MockVaultService) CalculateDailyEarning(disciplineIndex int) float64 {
	args := m.Called(disciplineIndex)
	return args.Get(0).(float64)
}

// ============== TESTS ==============

func TestFastingService_StartFast_Success(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	// No active session exists
	mockRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))
	mockVault.On("GetCurrentParticipation", ctx, userID).Return(nil, errors.New("not found"))
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)

	session, err := service.StartFast(ctx, userID, domain.Plan168, 16, nil)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, domain.StatusActive, session.Status)
	assert.Equal(t, 16, session.GoalHours)
	assert.Equal(t, domain.Plan168, session.PlanType)
	mockRepo.AssertExpectations(t)
}

func TestFastingService_StartFast_AlreadyActive(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	activeSession := &domain.FastingSession{
		ID:     uuid.New(),
		UserID: userID,
		Status: domain.StatusActive,
	}

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(activeSession, nil)

	session, err := service.StartFast(ctx, userID, domain.Plan168, 16, nil)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "active fasting session already exists")
}

func TestFastingService_StartFast_AllPlanTypes(t *testing.T) {
	planTypes := []struct {
		plan      domain.FastingPlanType
		goalHours int
	}{
		{domain.Plan168, 16},
		{domain.Plan186, 18},
		{domain.PlanOMAD, 23},
		{domain.Plan24h, 24},
		{domain.Plan36h, 36},
	}

	for _, tc := range planTypes {
		t.Run(string(tc.plan), func(t *testing.T) {
			mockRepo := new(MockFastingRepository)
			mockVault := new(MockVaultService)
			mockUserRepo := new(MockUserRepository)

			service := NewFastingService(mockRepo, mockVault, mockUserRepo)
			ctx := context.Background()
			userID := uuid.New()

			mockRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))
			mockVault.On("GetCurrentParticipation", ctx, userID).Return(nil, errors.New("not found"))
			mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)

			session, err := service.StartFast(ctx, userID, tc.plan, tc.goalHours, nil)

			assert.NoError(t, err)
			assert.NotNil(t, session)
			assert.Equal(t, tc.plan, session.PlanType)
			assert.Equal(t, tc.goalHours, session.GoalHours)
		})
	}
}

func TestFastingService_StartFast_CustomStartTime(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Start time 2 hours ago
	customStart := time.Now().Add(-2 * time.Hour)

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))
	mockVault.On("GetCurrentParticipation", ctx, userID).Return(nil, errors.New("not found"))
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)

	session, err := service.StartFast(ctx, userID, domain.Plan168, 16, &customStart)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, customStart.Unix(), session.StartTime.Unix())
}

func TestFastingService_StopFast_Success(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	startTime := time.Now().Add(-17 * time.Hour) // 17 hours ago
	activeSession := &domain.FastingSession{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    domain.StatusActive,
		StartTime: startTime,
		GoalHours: 16,
		PlanType:  domain.Plan168,
	}

	user := &domain.User{
		ID:              userID,
		DisciplineIndex: 50,
	}

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(activeSession, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockVault.On("UpdateDisciplineIndex", ctx, user, true, false).Return()
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	session, err := service.StopFast(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, domain.StatusCompleted, session.Status)
	assert.NotNil(t, session.EndTime)
	assert.True(t, session.ActualDurationHours >= 17)
}

func TestFastingService_StopFast_NoActiveSession(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(nil, errors.New("not found"))

	session, err := service.StopFast(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestFastingService_StopFast_EarlyEnd(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	startTime := time.Now().Add(-10 * time.Hour) // Only 10 hours, goal was 16
	activeSession := &domain.FastingSession{
		ID:        uuid.New(),
		UserID:    userID,
		Status:    domain.StatusActive,
		StartTime: startTime,
		GoalHours: 16,
		PlanType:  domain.Plan168,
	}

	user := &domain.User{
		ID:              userID,
		DisciplineIndex: 50,
	}

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(activeSession, nil)
	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockVault.On("CalculatePrice", ctx, user).Return(10.0)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	session, err := service.StopFast(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, domain.StatusCompleted, session.Status)
	// Discipline should be reduced for early end
	assert.True(t, user.DisciplineIndex < 50)
}

func TestFastingService_StopFast_PhaseCalculation(t *testing.T) {
	testCases := []struct {
		hours         float64
		expectedPhase string
	}{
		{6, "Anabolic"},
		{14, "Catabolic"},
		{19, "Ketosis"},
		{25, "Autophagy"},
		{50, "Deep Autophagy"},
		{75, "Immune Regeneration"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedPhase, func(t *testing.T) {
			mockRepo := new(MockFastingRepository)
			mockVault := new(MockVaultService)
			mockUserRepo := new(MockUserRepository)

			service := NewFastingService(mockRepo, mockVault, mockUserRepo)
			ctx := context.Background()
			userID := uuid.New()

			startTime := time.Now().Add(-time.Duration(tc.hours) * time.Hour)
			activeSession := &domain.FastingSession{
				ID:        uuid.New(),
				UserID:    userID,
				Status:    domain.StatusActive,
				StartTime: startTime,
				GoalHours: 16,
				PlanType:  domain.Plan168,
			}

			user := &domain.User{
				ID:              userID,
				DisciplineIndex: 50,
			}

			mockRepo.On("FindActiveByUserID", ctx, userID).Return(activeSession, nil)
			mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
			mockVault.On("UpdateDisciplineIndex", ctx, user, mock.Anything, false).Return()
			mockVault.On("CalculatePrice", ctx, user).Return(10.0)
			mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.FastingSession")).Return(nil)
			mockUserRepo.On("Save", ctx, user).Return(nil)

			session, err := service.StopFast(ctx, userID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPhase, session.PhaseReached)
		})
	}
}

func TestFastingService_GetCurrentFast_Active(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	activeSession := &domain.FastingSession{
		ID:     uuid.New(),
		UserID: userID,
		Status: domain.StatusActive,
	}

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(activeSession, nil)

	session, err := service.GetCurrentFast(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, domain.StatusActive, session.Status)
}

func TestFastingService_GetCurrentFast_None(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindActiveByUserID", ctx, userID).Return(nil, nil)

	session, err := service.GetCurrentFast(ctx, userID)

	assert.NoError(t, err)
	assert.Nil(t, session)
}

func TestFastingService_GetFastingHistory(t *testing.T) {
	mockRepo := new(MockFastingRepository)
	mockVault := new(MockVaultService)
	mockUserRepo := new(MockUserRepository)

	service := NewFastingService(mockRepo, mockVault, mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	history := []domain.FastingSession{
		{ID: uuid.New(), UserID: userID, Status: domain.StatusCompleted},
		{ID: uuid.New(), UserID: userID, Status: domain.StatusCompleted},
	}

	mockRepo.On("FindByUserID", ctx, userID).Return(history, nil)

	sessions, err := service.GetFastingHistory(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, sessions, 2)
}
