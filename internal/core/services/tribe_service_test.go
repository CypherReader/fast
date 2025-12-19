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

// MockTribeRepository is a mock implementation of ports.TribeRepository
type MockTribeRepository struct {
	mock.Mock
}

func (m *MockTribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func (m *MockTribeRepository) Update(ctx context.Context, tribe *domain.Tribe) error {
	args := m.Called(ctx, tribe)
	return args.Error(0)
}

func (m *MockTribeRepository) Delete(ctx context.Context, tribeID string) error {
	args := m.Called(ctx, tribeID)
	return args.Error(0)
}

func (m *MockTribeRepository) FindByID(ctx context.Context, id string) (*domain.Tribe, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tribe, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tribe), args.Error(1)
}

func (m *MockTribeRepository) List(ctx context.Context, query domain.ListTribesQuery) ([]domain.Tribe, int, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.Tribe), args.Int(1), args.Error(2)
}

func (m *MockTribeRepository) SaveMembership(ctx context.Context, membership *domain.TribeMembership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}

func (m *MockTribeRepository) UpdateMembership(ctx context.Context, membership *domain.TribeMembership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}

func (m *MockTribeRepository) DeleteMembership(ctx context.Context, tribeID, userID string) error {
	args := m.Called(ctx, tribeID, userID)
	return args.Error(0)
}

func (m *MockTribeRepository) FindMembership(ctx context.Context, tribeID, userID string) (*domain.TribeMembership, error) {
	args := m.Called(ctx, tribeID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TribeMembership), args.Error(1)
}

func (m *MockTribeRepository) GetMembersByTribeID(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	args := m.Called(ctx, tribeID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.TribeMember), args.Error(1)
}

func (m *MockTribeRepository) GetUserTribes(ctx context.Context, userID, status string) ([]domain.Tribe, error) {
	args := m.Called(ctx, userID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tribe), args.Error(1)
}

func (m *MockTribeRepository) UpdateMemberCounts(ctx context.Context, tribeID string) error {
	args := m.Called(ctx, tribeID)
	return args.Error(0)
}

func (m *MockTribeRepository) GetMembershipCount(ctx context.Context, tribeID string) (int, error) {
	args := m.Called(ctx, tribeID)
	return args.Int(0), args.Error(1)
}

func (m *MockTribeRepository) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	args := m.Called(ctx, tribeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TribeStats), args.Error(1)
}

// ============== CREATE TRIBE TESTS ==============

func TestTribeService_CreateTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	userID := uuid.New().String()

	req := domain.CreateTribeRequest{
		Name:            "Fasting Warriors",
		Description:     "A tribe for dedicated fasters",
		FastingSchedule: "16:8",
		PrimaryGoal:     "Weight Loss",
		Privacy:         "public",
	}

	// Slug doesn't exist yet
	mockRepo.On("FindBySlug", ctx, "fasting-warriors").Return(nil, errors.New("not found"))
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)
	mockRepo.On("SaveMembership", ctx, mock.AnythingOfType("*domain.TribeMembership")).Return(nil)

	tribe, err := service.CreateTribe(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	assert.Equal(t, "Fasting Warriors", tribe.Name)
	assert.Equal(t, "fasting-warriors", tribe.Slug)
	assert.Equal(t, userID, tribe.CreatorID)
	assert.Equal(t, 1, tribe.MemberCount)
	mockRepo.AssertExpectations(t)
}

func TestTribeService_CreateTribe_DuplicateSlug(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	userID := uuid.New().String()

	req := domain.CreateTribeRequest{
		Name:        "Fasting Warriors",
		Description: "A tribe for dedicated fasters",
		Privacy:     "public",
	}

	existingTribe := &domain.Tribe{ID: "existing-id", Slug: "fasting-warriors"}
	mockRepo.On("FindBySlug", ctx, "fasting-warriors").Return(existingTribe, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)
	mockRepo.On("SaveMembership", ctx, mock.AnythingOfType("*domain.TribeMembership")).Return(nil)

	tribe, err := service.CreateTribe(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, tribe)
	// Slug should have random suffix
	assert.Contains(t, tribe.Slug, "fasting-warriors-")
	assert.NotEqual(t, "fasting-warriors", tribe.Slug)
}

// ============== JOIN TRIBE TESTS ==============

func TestTribeService_JoinTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:      tribeID,
		Name:    "Test Tribe",
		Privacy: "public",
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(nil, errors.New("not found"))
	mockRepo.On("SaveMembership", ctx, mock.AnythingOfType("*domain.TribeMembership")).Return(nil)
	mockRepo.On("UpdateMemberCounts", ctx, tribeID).Return(nil)

	err := service.JoinTribe(ctx, tribeID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTribeService_JoinTribe_AlreadyMember(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{ID: tribeID, Name: "Test Tribe", Privacy: "public"}
	existingMembership := &domain.TribeMembership{
		TribeID: tribeID,
		UserID:  userID,
		Status:  "active",
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(existingMembership, nil)

	err := service.JoinTribe(ctx, tribeID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already a member")
}

func TestTribeService_JoinTribe_PrivatePending(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:      tribeID,
		Name:    "Private Tribe",
		Privacy: "private",
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(nil, errors.New("not found"))
	mockRepo.On("SaveMembership", ctx, mock.MatchedBy(func(m *domain.TribeMembership) bool {
		return m.Status == "pending" // Should be pending for private tribes
	})).Return(nil)
	mockRepo.On("UpdateMemberCounts", ctx, tribeID).Return(nil)

	err := service.JoinTribe(ctx, tribeID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTribeService_JoinTribe_Reactivate(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{ID: tribeID, Name: "Test Tribe", Privacy: "public"}
	leftMembership := &domain.TribeMembership{
		TribeID: tribeID,
		UserID:  userID,
		Status:  "left",
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(leftMembership, nil)
	mockRepo.On("UpdateMembership", ctx, mock.MatchedBy(func(m *domain.TribeMembership) bool {
		return m.Status == "active"
	})).Return(nil)

	err := service.JoinTribe(ctx, tribeID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============== LEAVE TRIBE TESTS ==============

func TestTribeService_LeaveTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	membership := &domain.TribeMembership{
		TribeID: tribeID,
		UserID:  userID,
		Role:    "member",
		Status:  "active",
	}

	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(membership, nil)
	mockRepo.On("DeleteMembership", ctx, tribeID, userID).Return(nil)
	mockRepo.On("UpdateMemberCounts", ctx, tribeID).Return(nil)

	err := service.LeaveTribe(ctx, tribeID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTribeService_LeaveTribe_NotMember(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(nil, nil)

	err := service.LeaveTribe(ctx, tribeID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a member")
}

func TestTribeService_LeaveTribe_CreatorCannotLeave(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	membership := &domain.TribeMembership{
		TribeID: tribeID,
		UserID:  userID,
		Role:    "creator",
		Status:  "active",
	}

	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(membership, nil)

	err := service.LeaveTribe(ctx, tribeID, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creator cannot leave")
}

// ============== GET TRIBE TESTS ==============

func TestTribeService_GetTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:          tribeID,
		Name:        "Test Tribe",
		MemberCount: 10,
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)

	result, err := service.GetTribe(ctx, tribeID, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Tribe", result.Name)
}

func TestTribeService_GetTribe_WithMembership(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{ID: tribeID, Name: "Test Tribe"}
	membership := &domain.TribeMembership{
		TribeID: tribeID,
		UserID:  userID,
		Role:    "member",
		Status:  "active",
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("FindMembership", ctx, tribeID, userID).Return(membership, nil)

	result, err := service.GetTribe(ctx, tribeID, &userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsJoined)
	assert.Equal(t, "member", result.UserRole)
}

func TestTribeService_GetTribe_NotFound(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()

	mockRepo.On("FindByID", ctx, tribeID).Return(nil, errors.New("not found"))

	result, err := service.GetTribe(ctx, tribeID, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============== UPDATE TRIBE TESTS ==============

func TestTribeService_UpdateTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:          tribeID,
		Name:        "Test Tribe",
		Description: "Old description",
		CreatorID:   userID,
	}

	newDesc := "New description"
	req := domain.UpdateTribeRequest{
		Description: &newDesc,
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Tribe")).Return(nil)

	result, err := service.UpdateTribe(ctx, tribeID, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New description", result.Description)
}

func TestTribeService_UpdateTribe_Unauthorized(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	creatorID := uuid.New().String()
	otherUserID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:        tribeID,
		Name:      "Test Tribe",
		CreatorID: creatorID, // Different from requesting user
	}

	newDesc := "New description"
	req := domain.UpdateTribeRequest{
		Description: &newDesc,
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)

	result, err := service.UpdateTribe(ctx, tribeID, otherUserID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
}

// ============== DELETE TRIBE TESTS ==============

func TestTribeService_DeleteTribe_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	userID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:        tribeID,
		CreatorID: userID,
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)
	mockRepo.On("Delete", ctx, tribeID).Return(nil)

	err := service.DeleteTribe(ctx, tribeID, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTribeService_DeleteTribe_Unauthorized(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()
	creatorID := uuid.New().String()
	otherUserID := uuid.New().String()

	tribe := &domain.Tribe{
		ID:        tribeID,
		CreatorID: creatorID,
	}

	mockRepo.On("FindByID", ctx, tribeID).Return(tribe, nil)

	err := service.DeleteTribe(ctx, tribeID, otherUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

// ============== GET TRIBE MEMBERS TESTS ==============

func TestTribeService_GetTribeMembers_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()

	members := []domain.TribeMember{
		{TribeMembership: domain.TribeMembership{UserID: uuid.New().String(), Role: "creator"}, UserName: "User 1"},
		{TribeMembership: domain.TribeMembership{UserID: uuid.New().String(), Role: "member"}, UserName: "User 2"},
	}

	mockRepo.On("GetMembersByTribeID", ctx, tribeID, 20, 0).Return(members, nil)

	result, err := service.GetTribeMembers(ctx, tribeID, 0, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestTribeService_GetTribeMembers_LimitEnforced(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	tribeID := uuid.New().String()

	// Request 500, should be capped at 100
	mockRepo.On("GetMembersByTribeID", ctx, tribeID, 100, 0).Return([]domain.TribeMember{}, nil)

	_, err := service.GetTribeMembers(ctx, tribeID, 500, 0)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============== GET MY TRIBES TESTS ==============

func TestTribeService_GetMyTribes_Success(t *testing.T) {
	mockRepo := new(MockTribeRepository)
	service := NewTribeService(mockRepo)
	ctx := context.Background()
	userID := uuid.New().String()

	tribes := []domain.Tribe{
		{ID: uuid.New().String(), Name: "Tribe 1"},
		{ID: uuid.New().String(), Name: "Tribe 2"},
	}

	mockRepo.On("GetUserTribes", ctx, userID, "active").Return(tribes, nil)

	result, err := service.GetMyTribes(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============== SLUG GENERATION TESTS ==============

func TestGenerateSlug(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Fasting Warriors", "fasting-warriors"},
		{"OMAD Crew!!!", "omad-crew"},
		{"The 16:8 Club", "the-168-club"},
		{"Low---Carb   Life", "low-carb-life"},
		{"  Spaces Around  ", "spaces-around"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := generateSlug(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGenerateSlug_LongName(t *testing.T) {
	longName := "This is a very long tribe name that should be truncated to sixty characters maximum because we have a limit"
	result := generateSlug(longName)
	assert.LessOrEqual(t, len(result), 60)
}
