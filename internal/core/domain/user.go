package domain

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionTier string
type SubscriptionStatus string

const (
	TierFree  SubscriptionTier = "free"
	TierVault SubscriptionTier = "vault"
)

const (
	SubStatusActive   SubscriptionStatus = "active"
	SubStatusCanceled SubscriptionStatus = "canceled"
	SubStatusPastDue  SubscriptionStatus = "past_due"
	SubStatusNone     SubscriptionStatus = ""
)

type User struct {
	ID                 uuid.UUID          `json:"id"`
	Email              string             `json:"email"`
	PasswordHash       string             `json:"-"`
	SubscriptionTier   SubscriptionTier   `json:"subscription_tier"`
	SubscriptionStatus SubscriptionStatus `json:"subscription_status"`
	SubscriptionID     string             `json:"subscription_id,omitempty"`
	DisciplineIndex    float64            `json:"discipline_index"` // 0-100
	CurrentPrice       float64            `json:"current_price"`    // The Lazy Tax (Deprecated)
	VaultDeposit       float64            `json:"vault_deposit"`    // Monthly deposit (e.g., $20)
	EarnedRefund       float64            `json:"earned_refund"`    // Amount earned back so far
	TribeID            *uuid.UUID         `json:"tribe_id,omitempty"`
	ReferralCode       string             `json:"referral_code,omitempty"`
	SignedContract     bool               `json:"signed_contract"`
	CreatedAt          time.Time          `json:"created_at"`
}

func (u *User) IsVaultMember() bool {
	return u.SubscriptionTier == TierVault && u.SubscriptionStatus == SubStatusActive
}
