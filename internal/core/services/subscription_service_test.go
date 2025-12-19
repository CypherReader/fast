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

// --- Mock PaymentGateway for SubscriptionService tests ---

type MockPaymentGatewayForSubscription struct {
	mock.Mock
}

func (m *MockPaymentGatewayForSubscription) CreateCustomer(email, name string) (string, error) {
	args := m.Called(email, name)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForSubscription) CreateSubscription(customerID, priceID string) (string, error) {
	args := m.Called(customerID, priceID)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForSubscription) CreatePayout(amount float64, currency, destination string) (string, error) {
	args := m.Called(amount, currency, destination)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForSubscription) ConstructEvent(payload []byte, header string) (interface{}, error) {
	args := m.Called(payload, header)
	return args.Get(0), args.Error(1)
}

// --- Tests for SubscriptionService ---

func TestNewSubscriptionService(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)

	svc := NewSubscriptionService(userRepo, pg)
	assert.NotNil(t, svc)
}

func TestSubscriptionService_UpgradeToVault_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierFree,
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateCustomer", user.Email, user.Email).Return("cus_new123", nil)
	pg.On("CreateSubscription", "cus_new123", mock.AnythingOfType("string")).Return("sub_new123", nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.UpgradeToVault(ctx, userID, "pm_test123")
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	pg.AssertExpectations(t)
}

func TestSubscriptionService_UpgradeToVault_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("user not found"))

	err := svc.UpgradeToVault(ctx, uuid.New(), "pm_test123")
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpgradeToVault_AlreadyVaultMember(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:                 userID,
		Email:              "test@example.com",
		SubscriptionTier:   domain.TierVault,
		SubscriptionStatus: domain.SubStatusActive, // IsVaultMember() requires both TierVault AND SubStatusActive
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)

	err := svc.UpgradeToVault(ctx, userID, "pm_test123")
	assert.Error(t, err)
	assert.Equal(t, "user is already a vault member", err.Error())
	// Only assert userRepo expectations as pg is not called when user is already vault member
	userRepo.AssertExpectations(t)
}

func TestSubscriptionService_UpgradeToVault_CreateCustomerError(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierFree,
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateCustomer", user.Email, user.Email).Return("", errors.New("payment gateway error"))

	err := svc.UpgradeToVault(ctx, userID, "pm_test123")
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	pg.AssertExpectations(t)
}

func TestSubscriptionService_UpgradeToVault_CreateSubscriptionError(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierFree,
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateCustomer", user.Email, user.Email).Return("cus_new123", nil)
	pg.On("CreateSubscription", "cus_new123", mock.AnythingOfType("string")).Return("", errors.New("subscription creation failed"))

	err := svc.UpgradeToVault(ctx, userID, "pm_test123")
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	pg.AssertExpectations(t)
}

func TestSubscriptionService_UpgradeToVault_SaveUserError(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierFree,
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateCustomer", user.Email, user.Email).Return("cus_new123", nil)
	pg.On("CreateSubscription", "cus_new123", mock.AnythingOfType("string")).Return("sub_new123", nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

	err := svc.UpgradeToVault(ctx, userID, "pm_test123")
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
	pg.AssertExpectations(t)
}

func TestSubscriptionService_DowngradeToFree_Success(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierVault,
		SubscriptionID:   "sub_existing",
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Run(func(args mock.Arguments) {
		savedUser := args.Get(1).(*domain.User)
		assert.Equal(t, domain.TierFree, savedUser.SubscriptionTier)
		assert.Empty(t, savedUser.SubscriptionID)
	}).Return(nil)

	err := svc.DowngradeToFree(ctx, userID)
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestSubscriptionService_DowngradeToFree_UserNotFound(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("user not found"))

	err := svc.DowngradeToFree(ctx, uuid.New())
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSubscriptionService_DowngradeToFree_SaveError(t *testing.T) {
	userRepo := new(MockUserRepository)
	pg := new(MockPaymentGatewayForSubscription)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		SubscriptionTier: domain.TierVault,
	}

	svc := NewSubscriptionService(userRepo, pg)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

	err := svc.DowngradeToFree(ctx, userID)
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}
