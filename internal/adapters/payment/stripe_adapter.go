package payment

import (
	"fastinghero/internal/core/ports"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/payout"
	"github.com/stripe/stripe-go/v74/subscription"
	"github.com/stripe/stripe-go/v74/webhook"
)

type StripeAdapter struct {
	apiKey        string
	webhookSecret string
}

// Ensure StripeAdapter implements PaymentGateway
var _ ports.PaymentGateway = (*StripeAdapter)(nil)

func NewStripeAdapter(apiKey, webhookSecret string) *StripeAdapter {
	stripe.Key = apiKey
	return &StripeAdapter{
		apiKey:        apiKey,
		webhookSecret: webhookSecret,
	}
}

func (s *StripeAdapter) CreateCustomer(email, name string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return c.ID, nil
}

func (s *StripeAdapter) CreateSubscription(customerID, priceID string) (string, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
	}
	sub, err := subscription.New(params)
	if err != nil {
		return "", err
	}
	return sub.ID, nil
}

func (s *StripeAdapter) CreatePayout(amount float64, currency, destination string) (string, error) {
	// Stripe expects amount in cents for USD
	amountInt := int64(amount * 100)

	params := &stripe.PayoutParams{
		Amount:      stripe.Int64(amountInt),
		Currency:    stripe.String(currency),
		Destination: stripe.String(destination),
	}
	p, err := payout.New(params)
	if err != nil {
		return "", err
	}
	return p.ID, nil
}

func (s *StripeAdapter) ConstructEvent(payload []byte, header string) (interface{}, error) {
	event, err := webhook.ConstructEvent(payload, header, s.webhookSecret)
	if err != nil {
		return nil, err
	}
	return event, nil
}
