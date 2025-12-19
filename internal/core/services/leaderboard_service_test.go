package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLeaderboardRepository is a mock of ports.LeaderboardRepository
type MockLeaderboardRepository struct {
	mock.Mock
}

func (m *MockLeaderboardRepository) GetGlobalLeaderboard(ctx context.Context, limit int) ([]domain.LeaderboardEntry, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LeaderboardEntry), args.Error(1)
}

func (m *MockLeaderboardRepository) GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.LeaderboardEntry, error) {
	args := m.Called(ctx, tribeID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.LeaderboardEntry), args.Error(1)
}

// ============== GET GLOBAL LEADERBOARD TESTS ==============

func TestLeaderboardService_GetGlobalLeaderboard_Success(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()

	entries := []domain.LeaderboardEntry{
		{UserID: uuid.New(), Rank: 1, DisciplineScore: 100},
		{UserID: uuid.New(), Rank: 2, DisciplineScore: 95},
		{UserID: uuid.New(), Rank: 3, DisciplineScore: 90},
	}

	mockRepo.On("GetGlobalLeaderboard", ctx, 50).Return(entries, nil)

	result, err := service.GetGlobalLeaderboard(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].Rank)
}

func TestLeaderboardService_GetGlobalLeaderboard_Empty(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetGlobalLeaderboard", ctx, 50).Return([]domain.LeaderboardEntry{}, nil)

	result, err := service.GetGlobalLeaderboard(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestLeaderboardService_GetGlobalLeaderboard_Error(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetGlobalLeaderboard", ctx, 50).Return(nil, errors.New("db error"))

	result, err := service.GetGlobalLeaderboard(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============== GET TRIBE LEADERBOARD TESTS ==============

func TestLeaderboardService_GetTribeLeaderboard_Success(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New()

	entries := []domain.LeaderboardEntry{
		{UserID: uuid.New(), Rank: 1, DisciplineScore: 50},
		{UserID: uuid.New(), Rank: 2, DisciplineScore: 45},
	}

	mockRepo.On("GetTribeLeaderboard", ctx, tribeID, 50).Return(entries, nil)

	result, err := service.GetTribeLeaderboard(ctx, tribeID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLeaderboardService_GetTribeLeaderboard_Empty(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New()

	mockRepo.On("GetTribeLeaderboard", ctx, tribeID, 50).Return([]domain.LeaderboardEntry{}, nil)

	result, err := service.GetTribeLeaderboard(ctx, tribeID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestLeaderboardService_GetTribeLeaderboard_Error(t *testing.T) {
	mockRepo := new(MockLeaderboardRepository)
	service := NewLeaderboardService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New()

	mockRepo.On("GetTribeLeaderboard", ctx, tribeID, 50).Return(nil, errors.New("db error"))

	result, err := service.GetTribeLeaderboard(ctx, tribeID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
