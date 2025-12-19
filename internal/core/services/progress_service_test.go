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

// MockProgressRepository is a mock implementation of ports.ProgressRepository
type MockProgressRepository struct {
	mock.Mock
}

func (m *MockProgressRepository) SaveWeightLog(ctx context.Context, log *domain.WeightLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockProgressRepository) GetWeightHistory(ctx context.Context, userID uuid.UUID, days int) ([]domain.WeightLog, error) {
	args := m.Called(ctx, userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.WeightLog), args.Error(1)
}

func (m *MockProgressRepository) SaveHydrationLog(ctx context.Context, log *domain.HydrationLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockProgressRepository) GetHydrationLog(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.HydrationLog, error) {
	args := m.Called(ctx, userID, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.HydrationLog), args.Error(1)
}

func (m *MockProgressRepository) SaveActivityLog(ctx context.Context, log *domain.ActivityLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockProgressRepository) GetActivityStats(ctx context.Context, userID uuid.UUID, days int) ([]domain.ActivityLog, error) {
	args := m.Called(ctx, userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.ActivityLog), args.Error(1)
}

// ============== LOG WEIGHT TESTS ==============

func TestProgressService_LogWeight_SuccessLbs(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveWeightLog", ctx, mock.AnythingOfType("*domain.WeightLog")).Return(nil)

	log, err := service.LogWeight(ctx, userID, 175.5, "lbs")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 175.5, log.WeightLbs)
	assert.InDelta(t, 79.6, log.WeightKg, 0.2)
	mockRepo.AssertExpectations(t)
}

func TestProgressService_LogWeight_SuccessKg(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveWeightLog", ctx, mock.AnythingOfType("*domain.WeightLog")).Return(nil)

	log, err := service.LogWeight(ctx, userID, 80.0, "kg")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 80.0, log.WeightKg)
	assert.InDelta(t, 176.4, log.WeightLbs, 0.2)
}

func TestProgressService_LogWeight_InvalidUnit(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	log, err := service.LogWeight(ctx, userID, 175.5, "stones")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "invalid unit")
}

func TestProgressService_LogWeight_ZeroWeight(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	log, err := service.LogWeight(ctx, userID, 0, "lbs")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "positive number")
}

func TestProgressService_LogWeight_NegativeWeight(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	log, err := service.LogWeight(ctx, userID, -50, "lbs")

	assert.Error(t, err)
	assert.Nil(t, log)
}

func TestProgressService_LogWeight_TooHigh(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Over 500kg limit
	log, err := service.LogWeight(ctx, userID, 600.0, "kg")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "must be between")
}

func TestProgressService_LogWeight_TooLow(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Under 20kg limit
	log, err := service.LogWeight(ctx, userID, 10.0, "kg")

	assert.Error(t, err)
	assert.Nil(t, log)
}

// ============== GET WEIGHT HISTORY TESTS ==============

func TestProgressService_GetWeightHistory_Success(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	history := []domain.WeightLog{
		{ID: uuid.New(), UserID: userID, WeightLbs: 175.0, LoggedAt: time.Now()},
		{ID: uuid.New(), UserID: userID, WeightLbs: 174.5, LoggedAt: time.Now().Add(-24 * time.Hour)},
	}

	mockRepo.On("GetWeightHistory", ctx, userID, 30).Return(history, nil)

	result, err := service.GetWeightHistory(ctx, userID, 30)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestProgressService_GetWeightHistory_Empty(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetWeightHistory", ctx, userID, 30).Return([]domain.WeightLog{}, nil)

	result, err := service.GetWeightHistory(ctx, userID, 30)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

// ============== LOG HYDRATION TESTS ==============

func TestProgressService_LogHydration_NewDay(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	// No existing log for today
	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(nil, nil)
	mockRepo.On("SaveHydrationLog", ctx, mock.AnythingOfType("*domain.HydrationLog")).Return(nil)

	log, err := service.LogHydration(ctx, userID, 1, "glasses")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 1, log.GlassesCount)
	mockRepo.AssertExpectations(t)
}

func TestProgressService_LogHydration_AddToExisting(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	existingLog := &domain.HydrationLog{
		ID:           uuid.New(),
		UserID:       userID,
		GlassesCount: 3,
		LoggedDate:   time.Now().Truncate(24 * time.Hour),
	}

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(existingLog, nil)
	mockRepo.On("SaveHydrationLog", ctx, existingLog).Return(nil)

	log, err := service.LogHydration(ctx, userID, 2, "glasses")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 5, log.GlassesCount) // 3 + 2 = 5
}

func TestProgressService_LogHydration_MlConversion(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(nil, nil)
	mockRepo.On("SaveHydrationLog", ctx, mock.AnythingOfType("*domain.HydrationLog")).Return(nil)

	// 500ml = 2 glasses (250ml per glass)
	log, err := service.LogHydration(ctx, userID, 500, "ml")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 2, log.GlassesCount)
}

func TestProgressService_LogHydration_SmallMlAmount(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(nil, nil)
	mockRepo.On("SaveHydrationLog", ctx, mock.AnythingOfType("*domain.HydrationLog")).Return(nil)

	// 100ml should still count as 1 glass minimum
	log, err := service.LogHydration(ctx, userID, 100, "ml")

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 1, log.GlassesCount)
}

func TestProgressService_LogHydration_InvalidUnit(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	log, err := service.LogHydration(ctx, userID, 1, "liters")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "invalid unit")
}

func TestProgressService_LogHydration_NegativeAmount(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	log, err := service.LogHydration(ctx, userID, -1, "glasses")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "positive")
}

func TestProgressService_LogHydration_TooMuchMl(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Over 10000ml limit
	log, err := service.LogHydration(ctx, userID, 15000, "ml")

	assert.Error(t, err)
	assert.Nil(t, log)
	assert.Contains(t, err.Error(), "too large")
}

// ============== GET DAILY HYDRATION TESTS ==============

func TestProgressService_GetDailyHydration_Found(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	existingLog := &domain.HydrationLog{
		ID:           uuid.New(),
		UserID:       userID,
		GlassesCount: 5,
		LoggedDate:   time.Now().Truncate(24 * time.Hour),
	}

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(existingLog, nil)

	log, err := service.GetDailyHydration(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, log)
	assert.Equal(t, 5, log.GlassesCount)
}

func TestProgressService_GetDailyHydration_NotFound(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(nil, nil)

	log, err := service.GetDailyHydration(ctx, userID)

	assert.NoError(t, err)
	assert.Nil(t, log)
}

func TestProgressService_GetDailyHydration_Error(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := NewProgressService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetHydrationLog", ctx, userID, mock.Anything).Return(nil, errors.New("db error"))

	log, err := service.GetDailyHydration(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, log)
}
