package domain

import (
	"time"

	"github.com/google/uuid"
)

type KetoSource string

const (
	SourceManual KetoSource = "manual"
	SourceDevice KetoSource = "device"
)

type KetoEntry struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	LoggedAt     time.Time  `json:"logged_at"`
	KetoneLevel  *float64   `json:"ketone_level,omitempty"`  // Premium
	AcetoneLevel *float64   `json:"acetone_level,omitempty"` // Premium
	Source       KetoSource `json:"source"`
}
