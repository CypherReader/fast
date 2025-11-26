package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReferralStatus string

const (
	ReferralStatusPending   ReferralStatus = "pending"
	ReferralStatusCompleted ReferralStatus = "completed"
)

type Referral struct {
	ID          uuid.UUID      `json:"id"`
	ReferrerID  uuid.UUID      `json:"referrer_id"`
	RefereeID   uuid.UUID      `json:"referee_id"`
	Status      ReferralStatus `json:"status"`
	RewardValue float64        `json:"reward_value"` // Amount awarded to each user
	CreatedAt   time.Time      `json:"created_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
}
