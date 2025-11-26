package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tribe struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	LeaderID        uuid.UUID `json:"leader_id"`
	MemberCount     int       `json:"member_count"`
	TotalDiscipline float64   `json:"total_discipline"`
	CreatedAt       time.Time `json:"created_at"`
}
