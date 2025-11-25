package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tribe struct {
	ID              uuid.UUID   `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	MemberIDs       []uuid.UUID `json:"member_ids"`
	CollectiveScore float64     `json:"collective_score"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
