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

// MockSocialRepository is a mock of ports.SocialRepository
type MockSocialRepository struct {
	mock.Mock
}

func (m *MockSocialRepository) SaveFriendNetwork(ctx context.Context, fn *domain.FriendNetwork) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func (m *MockSocialRepository) FindFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendNetwork, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FriendNetwork), args.Error(1)
}

func (m *MockSocialRepository) SaveTribe(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func (m *MockSocialRepository) FindTribeByID(ctx context.Context, tribeID uuid.UUID) (*domain.Tribe, error) {
	args := m.Called(ctx, tribeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockSocialRepository) FindAllTribes(ctx context.Context, limit, offset int) ([]domain.Tribe, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tribe), args.Error(1)
}

func (m *MockSocialRepository) SaveChallenge(ctx context.Context, challenge *domain.FriendChallenge) error {
	args := m.Called(ctx, challenge)
	return args.Error(0)
}

func (m *MockSocialRepository) FindChallengesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FriendChallenge, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.FriendChallenge), args.Error(1)
}

func (m *MockSocialRepository) SaveEvent(ctx context.Context, event *domain.SocialEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockSocialRepository) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.SocialEvent, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SocialEvent), args.Error(1)
}

// ============== ADD FRIEND TESTS ==============

func TestSocialService_AddFriend_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	friendID := uuid.New()

	mockRepo.On("FindFriends", ctx, userID).Return([]domain.FriendNetwork{}, nil)
	mockRepo.On("SaveFriendNetwork", ctx, mock.AnythingOfType("*domain.FriendNetwork")).Return(nil)

	err := service.AddFriend(ctx, userID, friendID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSocialService_AddFriend_CannotAddSelf(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	err := service.AddFriend(ctx, userID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot add self as friend")
}

func TestSocialService_AddFriend_AlreadyFriends(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	friendID := uuid.New()

	existingFriends := []domain.FriendNetwork{
		{ID: uuid.New(), UserID: userID, FriendID: friendID, Status: "accepted"},
	}

	mockRepo.On("FindFriends", ctx, userID).Return(existingFriends, nil)

	err := service.AddFriend(ctx, userID, friendID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already friends")
}

// ============== CREATE TRIBE TESTS ==============

func TestSocialService_CreateTribe_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveTribe", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)
	mockRepo.On("SaveEvent", ctx, mock.AnythingOfType("*domain.SocialEvent")).Return(nil)

	tribe, err := service.CreateTribe(ctx, userID, "Fasting Warriors", "A tribe for fasters", true)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	assert.Equal(t, "Fasting Warriors", tribe.Name)
	assert.Equal(t, "public", tribe.Privacy)
	assert.Equal(t, 1, tribe.MemberCount)
}

func TestSocialService_CreateTribe_Private(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveTribe", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)
	mockRepo.On("SaveEvent", ctx, mock.AnythingOfType("*domain.SocialEvent")).Return(nil)

	tribe, err := service.CreateTribe(ctx, userID, "Secret Fasters", "Private tribe", false)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	assert.Equal(t, "private", tribe.Privacy)
}

func TestSocialService_CreateTribe_SaveError(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("SaveTribe", ctx, mock.AnythingOfType("*domain.Tribe")).Return(errors.New("db error"))

	tribe, err := service.CreateTribe(ctx, userID, "Fasting Warriors", "A tribe", true)

	assert.Error(t, err)
	assert.Nil(t, tribe)
}

// ============== CREATE CHALLENGE TESTS ==============

func TestSocialService_CreateChallenge_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	startDate := time.Now()
	endDate := time.Now().Add(7 * 24 * time.Hour)

	mockRepo.On("SaveChallenge", ctx, mock.AnythingOfType("*domain.FriendChallenge")).Return(nil)

	challenge, err := service.CreateChallenge(ctx, userID, "7 Day Challenge", domain.ChallengeTypeFasting, 7, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, challenge)
	assert.Equal(t, "7 Day Challenge", challenge.Name)
	assert.Equal(t, domain.ChallengeTypeFasting, challenge.ChallengeType)
	assert.Equal(t, 7, challenge.Goal)
	assert.Equal(t, "active", challenge.Status)
}

// ============== GET FRIENDS TESTS ==============

func TestSocialService_GetFriends_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	friends := []domain.FriendNetwork{
		{ID: uuid.New(), UserID: userID, FriendID: uuid.New(), Status: "accepted"},
		{ID: uuid.New(), UserID: userID, FriendID: uuid.New(), Status: "pending"},
	}

	mockRepo.On("FindFriends", ctx, userID).Return(friends, nil)

	result, err := service.GetFriends(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestSocialService_GetFriends_Empty(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindFriends", ctx, userID).Return([]domain.FriendNetwork{}, nil)

	result, err := service.GetFriends(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

// ============== GET TRIBE TESTS ==============

func TestSocialService_GetTribe_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New()

	tribe := &domain.Tribe{
		ID:   tribeID.String(),
		Name: "Test Tribe",
	}

	mockRepo.On("FindTribeByID", ctx, tribeID).Return(tribe, nil)

	result, err := service.GetTribe(ctx, tribeID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Tribe", result.Name)
}

func TestSocialService_GetTribe_NotFound(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New()

	mockRepo.On("FindTribeByID", ctx, tribeID).Return(nil, errors.New("not found"))

	result, err := service.GetTribe(ctx, tribeID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============== LIST TRIBES TESTS ==============

func TestSocialService_ListTribes_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()

	tribes := []domain.Tribe{
		{ID: uuid.New().String(), Name: "Tribe 1"},
		{ID: uuid.New().String(), Name: "Tribe 2"},
	}

	mockRepo.On("FindAllTribes", ctx, 10, 0).Return(tribes, nil)

	result, err := service.ListTribes(ctx, 10, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============== GET CHALLENGES TESTS ==============

func TestSocialService_GetChallenges_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	challenges := []domain.FriendChallenge{
		{ID: uuid.New(), Name: "Challenge 1"},
		{ID: uuid.New(), Name: "Challenge 2"},
	}

	mockRepo.On("FindChallengesByUserID", ctx, userID).Return(challenges, nil)

	result, err := service.GetChallenges(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============== GET FEED TESTS ==============

func TestSocialService_GetFeed_Success(t *testing.T) {
	mockRepo := new(MockSocialRepository)
	service := NewSocialService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	events := []domain.SocialEvent{
		{ID: uuid.New(), EventType: domain.EventFastCompleted},
		{ID: uuid.New(), EventType: domain.EventTribeJoined},
	}

	mockRepo.On("GetFeed", ctx, userID, 20, 0).Return(events, nil)

	result, err := service.GetFeed(ctx, userID, 20, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
