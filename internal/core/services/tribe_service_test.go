package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTribeRepo
type MockTribeRepo struct {
	mock.Mock
}

func (m *MockTribeRepo) Save(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func (m *MockTribeRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeRepo) FindAll(ctx context.Context) ([]domain.Tribe, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Tribe), args.Error(1)
}

func (m *MockTribeRepo) AddMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeRepo) RemoveMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

// MockUserRepo
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Save(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) FindByReferralCode(ctx context.Context, code string) (*domain.User, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestCreateTribe(t *testing.T) {
	mockRepo := new(MockTribeRepo)
	mockUserRepo := new(MockUserRepo)
	service := NewTribeService(mockRepo, mockUserRepo)
	ctx := context.Background()

	creatorID := uuid.New()

	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)
	mockRepo.On("AddMember", ctx, mock.AnythingOfType("uuid.UUID"), creatorID).Return(nil)

	tribe, err := service.CreateTribe(ctx, "Spartans", "The elite warriors", creatorID)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	assert.Equal(t, "Spartans", tribe.Name)
	assert.Equal(t, "The elite warriors", tribe.Description)

	mockRepo.AssertExpectations(t)
}

func TestJoinTribe(t *testing.T) {
	mockRepo := new(MockTribeRepo)
	mockUserRepo := new(MockUserRepo)
	service := NewTribeService(mockRepo, mockUserRepo)
	ctx := context.Background()

	tribeID := uuid.New()
	newMemberID := uuid.New()

	// Mock user check
	user := &domain.User{ID: newMemberID} // User not in a tribe
	mockUserRepo.On("FindByID", ctx, newMemberID).Return(user, nil)

	// Mock AddMember
	mockRepo.On("AddMember", ctx, tribeID, newMemberID).Return(nil)

	err := service.JoinTribe(ctx, tribeID, newMemberID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
