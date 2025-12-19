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
	"github.com/stripe/stripe-go/v74"
)

// --- Mock PaymentGateway for StripeService tests ---

type MockPaymentGatewayForStripe struct {
	mock.Mock
}

func (m *MockPaymentGatewayForStripe) CreateCustomer(email, name string) (string, error) {
	args := m.Called(email, name)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForStripe) CreateSubscription(customerID, priceID string) (string, error) {
	args := m.Called(customerID, priceID)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForStripe) CreatePayout(amount float64, currency, destination string) (string, error) {
	args := m.Called(amount, currency, destination)
	return args.String(0), args.Error(1)
}

func (m *MockPaymentGatewayForStripe) ConstructEvent(payload []byte, header string) (interface{}, error) {
	args := m.Called(payload, header)
	return args.Get(0), args.Error(1)
}

// --- Mock SubscriptionRepository ---

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Save(ctx context.Context, sub *domain.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) FindByStripeSubscriptionID(ctx context.Context, stripeSubID string) (*domain.Subscription, error) {
	args := m.Called(ctx, stripeSubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Subscription), args.Error(1)
}

// --- Tests for StripeService ---

func TestNewStripeService(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)

	svc := NewStripeService(pg, subRepo, userRepo)
	assert.NotNil(t, svc)
}

func TestStripeService_CreateCustomer_ExistingCustomer(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)

	svc := NewStripeService(pg, subRepo, userRepo)

	user := &domain.User{
		ID:               uuid.New(),
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: "cus_existing123",
	}

	customerID, err := svc.CreateCustomer(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, "cus_existing123", customerID)
}

func TestStripeService_CreateCustomer_NewCustomer(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)

	svc := NewStripeService(pg, subRepo, userRepo)
	ctx := context.Background()

	user := &domain.User{
		ID:               uuid.New(),
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: "", // No existing customer
	}

	pg.On("CreateCustomer", user.Email, user.Name).Return("cus_new123", nil)
	userRepo.On("Save", ctx, user).Return(nil)

	customerID, err := svc.CreateCustomer(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, "cus_new123", customerID)
	pg.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestStripeService_CreateCustomer_PaymentGatewayError(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)

	svc := NewStripeService(pg, subRepo, userRepo)
	ctx := context.Background()

	user := &domain.User{
		ID:               uuid.New(),
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: "",
	}

	pg.On("CreateCustomer", user.Email, user.Name).Return("", errors.New("payment gateway error"))

	_, err := svc.CreateCustomer(ctx, user)
	assert.Error(t, err)
	pg.AssertExpectations(t)
}

func TestStripeService_CreateSubscription_Success(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: "cus_existing",
	}

	svc := NewStripeService(pg, subRepo, userRepo)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateSubscription", "cus_existing", "price_test").Return("sub_new123", nil)
	subRepo.On("Save", ctx, mock.AnythingOfType("*domain.Subscription")).Return(nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	sub, err := svc.CreateSubscription(ctx, userID, "price_test")
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, "sub_new123", sub.StripeSubscriptionID)
	userRepo.AssertExpectations(t)
	pg.AssertExpectations(t)
	subRepo.AssertExpectations(t)
}

func TestStripeService_CreateSubscription_UserNotFound(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	svc := NewStripeService(pg, subRepo, userRepo)

	userRepo.On("FindByID", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("user not found"))

	_, err := svc.CreateSubscription(ctx, uuid.New(), "price_test")
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestStripeService_CreateSubscription_CreatesCustomerIfMissing(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	userID := uuid.New()
	user := &domain.User{
		ID:               userID,
		Email:            "test@example.com",
		Name:             "Test User",
		StripeCustomerID: "", // No existing customer
	}

	svc := NewStripeService(pg, subRepo, userRepo)

	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	pg.On("CreateCustomer", user.Email, user.Name).Return("cus_new123", nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)
	pg.On("CreateSubscription", "cus_new123", "price_test").Return("sub_new123", nil)
	subRepo.On("Save", ctx, mock.AnythingOfType("*domain.Subscription")).Return(nil)

	sub, err := svc.CreateSubscription(ctx, userID, "price_test")
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	pg.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_ConstructEventError(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	svc := NewStripeService(pg, subRepo, userRepo)

	pg.On("ConstructEvent", []byte("{}"), "bad_sig").Return(nil, errors.New("invalid signature"))

	err := svc.HandleWebhook(ctx, []byte("{}"), "bad_sig")
	assert.Error(t, err)
	pg.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_InvalidEventType(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	svc := NewStripeService(pg, subRepo, userRepo)

	pg.On("ConstructEvent", []byte("{}"), "sig").Return("not an event", nil) // Return wrong type

	err := svc.HandleWebhook(ctx, []byte("{}"), "sig")
	assert.Error(t, err)
	assert.Equal(t, "invalid event type", err.Error())
	pg.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_SubscriptionNotFound(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	svc := NewStripeService(pg, subRepo, userRepo)

	event := stripe.Event{
		Type: "customer.subscription.updated",
		Data: &stripe.EventData{
			Raw: []byte(`{"id":"sub_test123","status":"active","current_period_start":1609459200,"current_period_end":1612137600}`),
		},
	}

	pg.On("ConstructEvent", []byte("{}"), "sig").Return(event, nil)
	subRepo.On("FindByStripeSubscriptionID", ctx, "sub_test123").Return(nil, nil) // Subscription not found

	// Should complete without error when sub not found
	err := svc.HandleWebhook(ctx, []byte("{}"), "sig")
	assert.NoError(t, err)
	pg.AssertExpectations(t)
	subRepo.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_SubscriptionUpdate(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	userID := uuid.New()
	subID := uuid.New()

	svc := NewStripeService(pg, subRepo, userRepo)

	event := stripe.Event{
		Type: "customer.subscription.updated",
		Data: &stripe.EventData{
			Raw: []byte(`{"id":"sub_test123","status":"active","current_period_start":1609459200,"current_period_end":1612137600}`),
		},
	}

	existingSub := &domain.Subscription{
		ID:                   subID,
		UserID:               userID,
		StripeSubscriptionID: "sub_test123",
		Status:               domain.SubStatusActive,
	}

	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
	}

	pg.On("ConstructEvent", []byte("{}"), "sig").Return(event, nil)
	subRepo.On("FindByStripeSubscriptionID", ctx, "sub_test123").Return(existingSub, nil)
	subRepo.On("Save", ctx, mock.AnythingOfType("*domain.Subscription")).Return(nil)
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.HandleWebhook(ctx, []byte("{}"), "sig")
	assert.NoError(t, err)
	pg.AssertExpectations(t)
	subRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_SubscriptionDeleted(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	userID := uuid.New()
	subID := uuid.New()

	svc := NewStripeService(pg, subRepo, userRepo)

	event := stripe.Event{
		Type: "customer.subscription.deleted",
		Data: &stripe.EventData{
			Raw: []byte(`{"id":"sub_test123","status":"canceled","current_period_start":1609459200,"current_period_end":1612137600}`),
		},
	}

	now := time.Now()
	existingSub := &domain.Subscription{
		ID:                   subID,
		UserID:               userID,
		StripeSubscriptionID: "sub_test123",
		Status:               domain.SubStatusActive,
		CurrentPeriodStart:   &now,
		CurrentPeriodEnd:     &now,
	}

	user := &domain.User{
		ID:    userID,
		Email: "test@example.com",
	}

	pg.On("ConstructEvent", []byte("{}"), "sig").Return(event, nil)
	subRepo.On("FindByStripeSubscriptionID", ctx, "sub_test123").Return(existingSub, nil)
	subRepo.On("Save", ctx, mock.AnythingOfType("*domain.Subscription")).Return(nil)
	userRepo.On("FindByID", ctx, userID).Return(user, nil)
	userRepo.On("Save", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.HandleWebhook(ctx, []byte("{}"), "sig")
	assert.NoError(t, err)
	pg.AssertExpectations(t)
	subRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestStripeService_HandleWebhook_FindByStripeSubError(t *testing.T) {
	pg := new(MockPaymentGatewayForStripe)
	subRepo := new(MockSubscriptionRepository)
	userRepo := new(MockUserRepository)
	ctx := context.Background()

	svc := NewStripeService(pg, subRepo, userRepo)

	event := stripe.Event{
		Type: "customer.subscription.updated",
		Data: &stripe.EventData{
			Raw: []byte(`{"id":"sub_test123","status":"active","current_period_start":1609459200,"current_period_end":1612137600}`),
		},
	}

	pg.On("ConstructEvent", []byte("{}"), "sig").Return(event, nil)
	subRepo.On("FindByStripeSubscriptionID", ctx, "sub_test123").Return(nil, errors.New("database error"))

	err := svc.HandleWebhook(ctx, []byte("{}"), "sig")
	assert.Error(t, err)
	pg.AssertExpectations(t)
	subRepo.AssertExpectations(t)
}
