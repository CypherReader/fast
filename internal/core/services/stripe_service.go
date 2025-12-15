package services

import (
	"context"
	"encoding/json"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v74"
)

type StripeService struct {
	paymentGateway ports.PaymentGateway
	subRepo        ports.SubscriptionRepository
	userRepo       ports.UserRepository
}

func NewStripeService(pg ports.PaymentGateway, subRepo ports.SubscriptionRepository, userRepo ports.UserRepository) *StripeService {
	return &StripeService{
		paymentGateway: pg,
		subRepo:        subRepo,
		userRepo:       userRepo,
	}
}

func (s *StripeService) CreateCustomer(ctx context.Context, user *domain.User) (string, error) {
	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	customerID, err := s.paymentGateway.CreateCustomer(user.Email, user.Name)
	if err != nil {
		return "", err
	}

	user.StripeCustomerID = customerID
	if err := s.userRepo.Save(ctx, user); err != nil {
		return "", err
	}

	return customerID, nil
}

func (s *StripeService) CreateSubscription(ctx context.Context, userID uuid.UUID, priceID string) (*domain.Subscription, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.StripeCustomerID == "" {
		// Create customer if not exists
		customerID, err := s.CreateCustomer(ctx, user)
		if err != nil {
			return nil, err
		}
		user.StripeCustomerID = customerID
	}

	subID, err := s.paymentGateway.CreateSubscription(user.StripeCustomerID, priceID)
	if err != nil {
		return nil, err
	}

	// Create subscription record
	sub := &domain.Subscription{
		ID:                   uuid.New(),
		UserID:               userID,
		StripeSubscriptionID: subID,
		Status:               domain.SubStatusActive, // Assuming active immediately for now
		PlanType:             string(domain.TierVault),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := s.subRepo.Save(ctx, sub); err != nil {
		return nil, err
	}

	// Update user status
	user.SubscriptionStatus = domain.SubStatusActive
	user.SubscriptionTier = domain.TierVault
	user.SubscriptionID = sub.ID.String()
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *StripeService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	eventInterface, err := s.paymentGateway.ConstructEvent(payload, signature)
	if err != nil {
		return err
	}

	event, ok := eventInterface.(stripe.Event)
	if !ok {
		return errors.New("invalid event type")
	}

	switch event.Type {
	case "customer.subscription.updated", "customer.subscription.deleted":
		var stripeSub stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &stripeSub)
		if err != nil {
			return err
		}

		sub, err := s.subRepo.FindByStripeSubscriptionID(ctx, stripeSub.ID)
		if err != nil {
			return err
		}
		if sub == nil {
			// Subscription not found, maybe log warning
			return nil
		}

		// Update status
		sub.Status = domain.SubscriptionStatus(stripeSub.Status)
		start := time.Unix(stripeSub.CurrentPeriodStart, 0)
		end := time.Unix(stripeSub.CurrentPeriodEnd, 0)
		sub.CurrentPeriodStart = &start
		sub.CurrentPeriodEnd = &end
		sub.CancelAtPeriodEnd = stripeSub.CancelAtPeriodEnd
		sub.UpdatedAt = time.Now()

		if err := s.subRepo.Save(ctx, sub); err != nil {
			return err
		}

		// Update user status
		user, err := s.userRepo.FindByID(ctx, sub.UserID)
		if err != nil {
			return err
		}

		user.SubscriptionStatus = sub.Status
		switch sub.Status {
		case domain.SubStatusActive:
			user.SubscriptionTier = domain.TierVault
		case domain.SubStatusCanceled, domain.SubStatusUnpaid:
			// Downgrade if needed, or keep as is until period end
			// For now, if not active, maybe revert to free?
			// user.SubscriptionTier = domain.TierFree
		}

		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}
	}

	return nil
}
