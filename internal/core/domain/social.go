package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SocialEventType string

const (
	SocialEventFastCompleted SocialEventType = "fast_completed"
	SocialEventKetoLogged    SocialEventType = "keto_logged"
	SocialEventTribeJoined   SocialEventType = "tribe_joined"
)

type SocialEvent struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	UserName  string          `json:"user_name"` // Denormalized for simpler feed queries
	Type      SocialEventType `json:"type"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}
