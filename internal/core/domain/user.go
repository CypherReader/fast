package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionTier string

const (
	TierFree    SubscriptionTier = "free"
	TierPremium SubscriptionTier = "premium"
	TierElite   SubscriptionTier = "elite"
)

type User struct {
	ID               uuid.UUID        `json:"id"`
	Email            string           `json:"email"`
	PasswordHash     string           `json:"-"`
	SubscriptionTier SubscriptionTier `json:"subscription_tier"`
	DisciplineIndex  float64          `json:"discipline_index"` // 0-100
	CurrentPrice     float64          `json:"current_price"`    // The Lazy Tax (Deprecated)
	VaultDeposit     float64          `json:"vault_deposit"`    // Monthly deposit (e.g., $20)
	EarnedRefund     float64          `json:"earned_refund"`    // Amount earned back so far
	TribeID          *uuid.UUID       `json:"tribe_id,omitempty"`
	SignedContract   bool             `json:"signed_contract"`
	CreatedAt        time.Time        `json:"created_at"`
}

func (u *User) IsPremium() bool {
	return u.SubscriptionTier == TierPremium || u.SubscriptionTier == TierElite
}

func (u *User) IsElite() bool {
	return u.SubscriptionTier == TierElite
}
