package domain

import (
	"time"

	"github.com/google/uuid"
)

type WeightLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	WeightLbs float64   `json:"weight_lbs"`
	WeightKg  float64   `json:"weight_kg"`
	LoggedAt  time.Time `json:"logged_at"`
	CreatedAt time.Time `json:"created_at"`
}

type HydrationLog struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	GlassesCount int       `json:"glasses_count"`
	LoggedDate   time.Time `json:"logged_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Steps          int       `json:"steps"`
	DistanceKm     float64   `json:"distance_km"`
	CaloriesBurned int       `json:"calories_burned"`
	LoggedDate     time.Time `json:"logged_date"`
	CreatedAt      time.Time `json:"created_at"`
}

type CommitmentContract struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	CommitmentText string    `json:"commitment_text"`
	Goals          string    `json:"goals"` // JSON string
	SignedAt       time.Time `json:"signed_at"`
}
