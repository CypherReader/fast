package domain

import (
	"time"

	"github.com/google/uuid"
)

// SOSFlare represents a distress signal sent by a struggling user
type SOSFlare struct {
	ID               uuid.UUID   `json:"id"`
	UserID           uuid.UUID   `json:"user_id"`
	FastingID        uuid.UUID   `json:"fasting_id"`
	TribeID          *uuid.UUID  `json:"tribe_id,omitempty"`
	Description      string      `json:"description"` // User's craving description
	HoursFasted      float64     `json:"hours_fasted"`
	Status           SOSStatus   `json:"status"`
	HypeCount        int         `json:"hype_count"`
	IsAnonymous      bool        `json:"is_anonymous"` // Hide identity from tribe
	RespondedUserIDs []uuid.UUID `json:"responded_user_ids,omitempty"`
	CortexResponded  bool        `json:"cortex_responded"` // Did Cortex auto-respond?
	CreatedAt        time.Time   `json:"created_at"`
	ResolvedAt       *time.Time  `json:"resolved_at,omitempty"`
}

type SOSStatus string

const (
	SOSStatusActive  SOSStatus = "active"
	SOSStatusRescued SOSStatus = "rescued" // User survived the urge
	SOSStatusFailed  SOSStatus = "failed"  // User broke the fast
)

// HypeResponse represents encouragement from a tribe member
type HypeResponse struct {
	ID         uuid.UUID `json:"id"`
	SOSID      uuid.UUID `json:"sos_id"`
	FromUserID uuid.UUID `json:"from_user_id"`
	FromName   string    `json:"from_name"` // For display in UI
	Message    string    `json:"message,omitempty"`
	Emoji      string    `json:"emoji"` // 🔥 or ⚡
	CreatedAt  time.Time `json:"created_at"`
}

// SOSSettings contains user preferences for SOS behavior
type SOSSettings struct {
	NotifyTribeOnSOS bool       `json:"notify_tribe_on_sos"`
	AnonymousMode    bool       `json:"anonymous_mode"` // Hide identity in tribe notifications
	LastSOSAt        *time.Time `json:"last_sos_at,omitempty"`
}

// GetHypeLimit returns the max hypes a user can send per day based on tribe size
// Scaling system: smaller tribes = unlimited, larger tribes = rate limited
func GetHypeLimit(tribeMemberCount int) int {
	switch {
	case tribeMemberCount <= 10:
		return 100 // Essentially unlimited for small tribes
	case tribeMemberCount <= 50:
		return 20 // Medium tribes
	case tribeMemberCount <= 200:
		return 10 // Large tribes
	default:
		return 5 // Mega tribes (prevent spam)
	}
}
