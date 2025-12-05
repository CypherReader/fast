package domain

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                   uuid.UUID          `json:"id"`
	UserID               uuid.UUID          `json:"user_id"`
	PlanType             string             `json:"plan_type"`
	SubscriptionPrice    float64            `json:"subscription_price"`
	StripeSubscriptionID string             `json:"stripe_subscription_id"`
	Status               SubscriptionStatus `json:"status"`
	CurrentPeriodStart   *time.Time         `json:"current_period_start"`
	CurrentPeriodEnd     *time.Time         `json:"current_period_end"`
	CancelAtPeriodEnd    bool               `json:"cancel_at_period_end"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}
