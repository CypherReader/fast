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

// MockReferralRepository is a mock of ports.ReferralRepository
type MockReferralRepository struct {
	mock.Mock
}

func (m *MockReferralRepository) Save(ctx context.Context, referral *domain.Referral) error {
	args := m.Called(ctx, referral)
	return args.Error(0)
}

func (m *MockReferralRepository) Update(ctx context.Context, referral *domain.Referral) error {
	args := m.Called(ctx, referral)
	return args.Error(0)
}

func (m *MockReferralRepository) FindByReferrerID(ctx context.Context, referrerID uuid.UUID) ([]domain.Referral, error) {
	args := m.Called(ctx, referrerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Referral), args.Error(1)
}

func (m *MockReferralRepository) FindByRefereeID(ctx context.Context, refereeID uuid.UUID) (*domain.Referral, error) {
	args := m.Called(ctx, refereeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Referral), args.Error(1)
}

// ============== GENERATE REFERRAL CODE TESTS ==============

func TestReferralService_GenerateReferralCode_Success(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:           userID,
		ReferralCode: "",
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	code, err := service.GenerateReferralCode(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Len(t, code, 6)
}

func TestReferralService_GenerateReferralCode_AlreadyExists(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:           userID,
		ReferralCode: "EXIST1",
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	code, err := service.GenerateReferralCode(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, "EXIST1", code) // Returns existing code
}

func TestReferralService_GenerateReferralCode_UserNotFound(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("not found"))

	code, err := service.GenerateReferralCode(ctx, userID)

	assert.Error(t, err)
	assert.Empty(t, code)
}

// ============== GET REFERRAL CODE TESTS ==============

func TestReferralService_GetReferralCode_Existing(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:           userID,
		ReferralCode: "MYCODE",
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)

	code, err := service.GetReferralCode(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, "MYCODE", code)
}

func TestReferralService_GetReferralCode_GenerateNew(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:           userID,
		ReferralCode: "",
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	code, err := service.GetReferralCode(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, code)
}

// ============== TRACK REFERRAL TESTS ==============

func TestReferralService_TrackReferral_Success(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	referrerID := uuid.New()
	refereeID := uuid.New()

	referrer := &domain.User{
		ID:           referrerID,
		ReferralCode: "REFER1",
	}

	mockUserRepo.On("FindByReferralCode", ctx, "REFER1").Return(referrer, nil)
	mockReferralRepo.On("FindByRefereeID", ctx, refereeID).Return(nil, errors.New("not found"))
	mockReferralRepo.On("Save", ctx, mock.AnythingOfType("*domain.Referral")).Return(nil)

	err := service.TrackReferral(ctx, "REFER1", refereeID)

	assert.NoError(t, err)
}

func TestReferralService_TrackReferral_EmptyCode(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	refereeID := uuid.New()

	err := service.TrackReferral(ctx, "", refereeID)

	assert.NoError(t, err) // Empty code is ignored, not an error
}

func TestReferralService_TrackReferral_SelfReferral(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	referrer := &domain.User{
		ID:           userID,
		ReferralCode: "MYCODE",
	}

	mockUserRepo.On("FindByReferralCode", ctx, "MYCODE").Return(referrer, nil)

	err := service.TrackReferral(ctx, "MYCODE", userID) // Same user

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot refer yourself")
}

func TestReferralService_TrackReferral_AlreadyReferred(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	referrerID := uuid.New()
	refereeID := uuid.New()

	referrer := &domain.User{ID: referrerID, ReferralCode: "REFER1"}
	existingReferral := &domain.Referral{ID: uuid.New()}

	mockUserRepo.On("FindByReferralCode", ctx, "REFER1").Return(referrer, nil)
	mockReferralRepo.On("FindByRefereeID", ctx, refereeID).Return(existingReferral, nil)

	err := service.TrackReferral(ctx, "REFER1", refereeID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already referred")
}

// ============== GET REFERRAL STATS TESTS ==============

func TestReferralService_GetReferralStats_WithReferrals(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	referrals := []domain.Referral{
		{Status: domain.ReferralStatusCompleted, RewardValue: 5.0},
		{Status: domain.ReferralStatusCompleted, RewardValue: 5.0},
		{Status: domain.ReferralStatusPending, RewardValue: 5.0}, // Not counted
	}

	mockReferralRepo.On("FindByReferrerID", ctx, userID).Return(referrals, nil)

	earned, count, err := service.GetReferralStats(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 10.0, earned) // 2 completed x $5
	assert.Equal(t, 2, count)
}

func TestReferralService_GetReferralStats_NoReferrals(t *testing.T) {
	mockReferralRepo := new(MockReferralRepository)
	mockUserRepo := new(MockUserRepository)
	mockVaultService := new(MockVaultService)

	service := NewReferralService(mockReferralRepo, mockUserRepo, mockVaultService)
	ctx := context.Background()
	userID := uuid.New()

	mockReferralRepo.On("FindByReferrerID", ctx, userID).Return([]domain.Referral{}, nil)

	earned, count, err := service.GetReferralStats(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 0.0, earned)
	assert.Equal(t, 0, count)
}

// ============== GENERATE RANDOM CODE TEST ==============

func TestGenerateRandomCode(t *testing.T) {
	code1, err1 := generateRandomCode(6)
	code2, err2 := generateRandomCode(6)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Len(t, code1, 6)
	assert.Len(t, code2, 6)
	assert.NotEqual(t, code1, code2) // Should be random
}
