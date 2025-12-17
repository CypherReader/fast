package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	userRepo       ports.UserRepository
	paymentGateway ports.PaymentGateway
}

func NewSubscriptionService(userRepo ports.UserRepository, paymentGateway ports.PaymentGateway) *SubscriptionService {
	return &SubscriptionService{
		userRepo:       userRepo,
		paymentGateway: paymentGateway,
	}
}

func (s *SubscriptionService) UpgradeToVault(ctx context.Context, userID uuid.UUID, paymentMethodID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsVaultMember() {
		return errors.New("user is already a vault member")
	}

	// 1. Create/Get Stripe Customer
	// For MVP, we assume we might need to create it or store it on the user
	// TODO: Add StripeCustomerID to User struct if not present
	customerID, err := s.paymentGateway.CreateCustomer(user.Email, user.Email)
	if err != nil {
		return err
	}

	// 2. Create Subscription
	// We need a Price ID for the Vault subscription. Ideally this is config.
	const VaultPriceID = "price_H5ggYJ..." // Replace with env var or constant
	subID, err := s.paymentGateway.CreateSubscription(customerID, VaultPriceID)
	if err != nil {
		return err
	}

	// 3. Update User
	user.SubscriptionTier = domain.TierVault
	user.SubscriptionStatus = domain.SubStatusActive
	user.SubscriptionID = subID
	user.VaultDeposit = 4.99 // UPDATED FOR V2: Standard subscription price

	return s.userRepo.Save(ctx, user)
}

func (s *SubscriptionService) DowngradeToFree(ctx context.Context, userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// In a real app, we would cancel the Stripe subscription here
	// s.paymentGateway.CancelSubscription(user.SubscriptionID)

	user.SubscriptionTier = domain.TierFree
	user.SubscriptionStatus = domain.SubStatusCanceled
	user.SubscriptionID = ""

	return s.userRepo.Save(ctx, user)
}
