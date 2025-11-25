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

func (m *MockTribeRepo) Create(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func (m *MockTribeRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeRepo) Update(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func TestCreateTribe(t *testing.T) {
	mockRepo := new(MockTribeRepo)
	service := NewTribeService(mockRepo)
	ctx := context.Background()

	creatorID := uuid.New()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)

	tribe, err := service.CreateTribe(ctx, "Spartans", creatorID)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	assert.Equal(t, "Spartans", tribe.Name)
	assert.Contains(t, tribe.MemberIDs, creatorID)
	mockRepo.AssertExpectations(t)
}

func TestJoinTribe(t *testing.T) {
	mockRepo := new(MockTribeRepo)
	service := NewTribeService(mockRepo)
	ctx := context.Background()

	tribeID := uuid.New()
	existingMemberID := uuid.New()
	newMemberID := uuid.New()

	existingTribe := &domain.Tribe{
		ID:        tribeID,
		Name:      "Spartans",
		MemberIDs: []uuid.UUID{existingMemberID},
	}

	mockRepo.On("GetByID", ctx, tribeID).Return(existingTribe, nil)
	mockRepo.On("Update", ctx, existingTribe).Return(nil)

	err := service.JoinTribe(ctx, tribeID, newMemberID)

	assert.NoError(t, err)
	assert.Contains(t, existingTribe.MemberIDs, newMemberID)
	assert.Equal(t, 2, len(existingTribe.MemberIDs))
	mockRepo.AssertExpectations(t)
}
