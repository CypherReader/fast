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

// MockGamificationRepository is a mock of ports.GamificationRepository
type MockGamificationRepository struct {
	mock.Mock
}

func (m *MockGamificationRepository) GetUserStreak(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserStreak), args.Error(1)
}

func (m *MockGamificationRepository) UpdateUserStreak(ctx context.Context, streak *domain.UserStreak) error {
	args := m.Called(ctx, streak)
	return args.Error(0)
}

func (m *MockGamificationRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]domain.UserBadge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.UserBadge), args.Error(1)
}

func (m *MockGamificationRepository) SaveUserBadge(ctx context.Context, badge *domain.UserBadge) error {
	args := m.Called(ctx, badge)
	return args.Error(0)
}

// ============== GET USER PROFILE TESTS ==============

func TestGamificationService_GetUserProfile_Success(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	streak := &domain.UserStreak{
		UserID:        userID,
		CurrentStreak: 5,
		LongestStreak: 10,
	}

	badges := []domain.UserBadge{
		{UserID: userID, BadgeID: domain.BadgeFirstFast, EarnedAt: time.Now()},
	}

	mockRepo.On("GetUserStreak", ctx, userID).Return(streak, nil)
	mockRepo.On("GetUserBadges", ctx, userID).Return(badges, nil)

	resultStreak, resultBadges, err := service.GetUserGamificationProfile(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resultStreak)
	assert.Equal(t, 5, resultStreak.CurrentStreak)
	assert.Len(t, resultBadges, 1)
}

func TestGamificationService_GetUserProfile_NoStreak(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	// No streak exists yet
	mockRepo.On("GetUserStreak", ctx, userID).Return(nil, nil)
	mockRepo.On("GetUserBadges", ctx, userID).Return([]domain.UserBadge{}, nil)

	resultStreak, resultBadges, err := service.GetUserGamificationProfile(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, resultStreak) // Should return empty streak
	assert.Equal(t, 0, resultStreak.CurrentStreak)
	assert.Len(t, resultBadges, 0)
}

func TestGamificationService_GetUserProfile_Error(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetUserStreak", ctx, userID).Return(nil, errors.New("db error"))

	resultStreak, resultBadges, err := service.GetUserGamificationProfile(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, resultStreak)
	assert.Nil(t, resultBadges)
}

// ============== UPDATE STREAK TESTS ==============

func TestGamificationService_UpdateStreak_FirstDay(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	// No existing streak
	mockRepo.On("GetUserStreak", ctx, userID).Return(nil, nil)
	mockRepo.On("UpdateUserStreak", ctx, mock.MatchedBy(func(s *domain.UserStreak) bool {
		return s.CurrentStreak == 1
	})).Return(nil)

	err := service.UpdateStreak(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGamificationService_UpdateStreak_ConsecutiveDay(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	yesterday := time.Now().Add(-24 * time.Hour)
	existingStreak := &domain.UserStreak{
		UserID:           userID,
		CurrentStreak:    5,
		LongestStreak:    10,
		LastActivityDate: yesterday,
	}

	mockRepo.On("GetUserStreak", ctx, userID).Return(existingStreak, nil)
	mockRepo.On("UpdateUserStreak", ctx, mock.MatchedBy(func(s *domain.UserStreak) bool {
		return s.CurrentStreak == 6
	})).Return(nil)

	err := service.UpdateStreak(ctx, userID)

	assert.NoError(t, err)
}

func TestGamificationService_UpdateStreak_StreakBroken(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	twoDaysAgo := time.Now().Add(-48 * time.Hour)
	existingStreak := &domain.UserStreak{
		UserID:           userID,
		CurrentStreak:    5,
		LongestStreak:    10,
		LastActivityDate: twoDaysAgo,
	}

	mockRepo.On("GetUserStreak", ctx, userID).Return(existingStreak, nil)
	mockRepo.On("UpdateUserStreak", ctx, mock.MatchedBy(func(s *domain.UserStreak) bool {
		return s.CurrentStreak == 1 // Streak reset
	})).Return(nil)

	err := service.UpdateStreak(ctx, userID)

	assert.NoError(t, err)
}

func TestGamificationService_UpdateStreak_SameDay(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	today := time.Now()
	existingStreak := &domain.UserStreak{
		UserID:           userID,
		CurrentStreak:    5,
		LastActivityDate: today,
	}

	mockRepo.On("GetUserStreak", ctx, userID).Return(existingStreak, nil)
	// Should not call UpdateUserStreak since it's the same day

	err := service.UpdateStreak(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "UpdateUserStreak")
}

func TestGamificationService_UpdateStreak_NewLongest(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	yesterday := time.Now().Add(-24 * time.Hour)
	existingStreak := &domain.UserStreak{
		UserID:           userID,
		CurrentStreak:    10,
		LongestStreak:    10,
		LastActivityDate: yesterday,
	}

	mockRepo.On("GetUserStreak", ctx, userID).Return(existingStreak, nil)
	mockRepo.On("UpdateUserStreak", ctx, mock.MatchedBy(func(s *domain.UserStreak) bool {
		return s.CurrentStreak == 11 && s.LongestStreak == 11
	})).Return(nil)

	err := service.UpdateStreak(ctx, userID)

	assert.NoError(t, err)
}

// ============== CHECK AND AWARD BADGES TESTS ==============

func TestGamificationService_CheckAndAwardBadges_FirstFast(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetUserBadges", ctx, userID).Return([]domain.UserBadge{}, nil)
	mockRepo.On("GetUserStreak", ctx, userID).Return(nil, nil)
	mockRepo.On("SaveUserBadge", ctx, mock.MatchedBy(func(b *domain.UserBadge) bool {
		return b.BadgeID == domain.BadgeFirstFast
	})).Return(nil)

	err := service.CheckAndAwardBadges(ctx, userID, "fast_completed", nil)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGamificationService_CheckAndAwardBadges_AlreadyHasBadge(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	existingBadges := []domain.UserBadge{
		{UserID: userID, BadgeID: domain.BadgeFirstFast},
	}

	mockRepo.On("GetUserBadges", ctx, userID).Return(existingBadges, nil)
	mockRepo.On("GetUserStreak", ctx, userID).Return(nil, nil)

	err := service.CheckAndAwardBadges(ctx, userID, "fast_completed", nil)

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "SaveUserBadge") // Should not award duplicate
}

func TestGamificationService_CheckAndAwardBadges_Streak3(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	streak := &domain.UserStreak{
		UserID:        userID,
		CurrentStreak: 3,
	}

	mockRepo.On("GetUserBadges", ctx, userID).Return([]domain.UserBadge{
		{BadgeID: domain.BadgeFirstFast}, // Already has first fast
	}, nil)
	mockRepo.On("GetUserStreak", ctx, userID).Return(streak, nil)
	mockRepo.On("SaveUserBadge", ctx, mock.MatchedBy(func(b *domain.UserBadge) bool {
		return b.BadgeID == domain.BadgeStreak3
	})).Return(nil)

	err := service.CheckAndAwardBadges(ctx, userID, "fast_completed", nil)

	assert.NoError(t, err)
}

func TestGamificationService_CheckAndAwardBadges_Streak7(t *testing.T) {
	mockRepo := new(MockGamificationRepository)
	mockFastingRepo := new(MockFastingRepository)
	service := NewGamificationService(mockRepo, mockFastingRepo)
	ctx := context.Background()
	userID := uuid.New()

	streak := &domain.UserStreak{
		UserID:        userID,
		CurrentStreak: 7,
	}

	mockRepo.On("GetUserBadges", ctx, userID).Return([]domain.UserBadge{
		{BadgeID: domain.BadgeFirstFast},
		{BadgeID: domain.BadgeStreak3},
	}, nil)
	mockRepo.On("GetUserStreak", ctx, userID).Return(streak, nil)
	mockRepo.On("SaveUserBadge", ctx, mock.MatchedBy(func(b *domain.UserBadge) bool {
		return b.BadgeID == domain.BadgeStreak7
	})).Return(nil)

	err := service.CheckAndAwardBadges(ctx, userID, "fast_completed", nil)

	assert.NoError(t, err)
}
