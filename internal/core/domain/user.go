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
	SubStatusUnpaid   SubscriptionStatus = "unpaid"
	SubStatusNone     SubscriptionStatus = ""
)

type User struct {
	ID                       uuid.UUID          `json:"id"`
	Email                    string             `json:"email"`
	PasswordHash             string             `json:"-"`
	Name                     string             `json:"name,omitempty"`
	OnboardingCompleted      bool               `json:"onboarding_completed"`
	Goal                     string             `json:"goal,omitempty"`
	FastingPlan              string             `json:"fasting_plan,omitempty"`
	Sex                      string             `json:"sex,omitempty"`
	HeightCm                 float64            `json:"height_cm,omitempty"`
	CurrentWeightLbs         float64            `json:"current_weight_lbs,omitempty"`
	TargetWeightLbs          float64            `json:"target_weight_lbs,omitempty"`
	Timezone                 string             `json:"timezone"`
	Units                    string             `json:"units"`
	StripeCustomerID         string             `json:"stripe_customer_id,omitempty"`
	SubscriptionTier         SubscriptionTier   `json:"subscription_tier"`
	SubscriptionStatus       SubscriptionStatus `json:"subscription_status"`
	SubscriptionID           string             `json:"subscription_id,omitempty"`
	VaultEnabled             bool               `json:"vault_enabled"`
	TrialEndsAt              *time.Time         `json:"trial_ends_at,omitempty"`
	DisciplineIndex          float64            `json:"discipline_index"` // 0-100
	CurrentPrice             float64            `json:"current_price"`    // The Lazy Tax (Deprecated)
	VaultDeposit             float64            `json:"vault_deposit"`    // Monthly deposit (e.g., $20)
	EarnedRefund             float64            `json:"earned_refund"`    // Amount earned back so far
	TribeID                  *uuid.UUID         `json:"tribe_id,omitempty"`
	ReferralCode             string             `json:"referral_code,omitempty"`
	SignedContract           bool               `json:"signed_contract"`
	PushNotificationsEnabled bool               `json:"push_notifications_enabled"`
	NotificationToken        string             `json:"notification_token,omitempty"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
}

func (u *User) IsVaultMember() bool {
	return u.SubscriptionTier == TierVault && u.SubscriptionStatus == SubStatusActive
}

type UserProfileUpdate struct {
	Name             *string  `json:"name"`
	Goal             *string  `json:"goal"`
	FastingPlan      *string  `json:"fasting_plan"`
	Sex              *string  `json:"sex"`
	HeightCm         *float64 `json:"height_cm"`
	CurrentWeightLbs *float64 `json:"current_weight_lbs"`
	TargetWeightLbs  *float64 `json:"target_weight_lbs"`
	Timezone         *string  `json:"timezone"`
	Units            *string  `json:"units"`
}
