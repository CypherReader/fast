package domain

import (
	"time"

	"github.com/google/uuid"
)

type TelemetrySource string

const (
	SourceGarmin          TelemetrySource = "garmin"
	SourceAppleHealth     TelemetrySource = "apple_health"
	SourceOura            TelemetrySource = "oura"
	SourceWhoop           TelemetrySource = "whoop"
	SourceGoogleFit       TelemetrySource = "google_fit"
	TelemetrySourceManual TelemetrySource = "manual"
)

type MetricType string

const (
	MetricSteps          MetricType = "steps"
	MetricWeight         MetricType = "weight"
	MetricActiveCalories MetricType = "active_calories"
	MetricBodyFat        MetricType = "body_fat"
	MetricSleepScore     MetricType = "sleep_score"
	MetricHRV            MetricType = "hrv"
)

type TelemetryData struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	Source     TelemetrySource `json:"source"`
	Type       MetricType      `json:"type"`
	Value      float64         `json:"value"`
	Unit       string          `json:"unit"`
	Timestamp  time.Time       `json:"timestamp"`
	IsManual   bool            `json:"is_manual"`
	TrustScore float64         `json:"trust_score"` // 0.0 to 1.0
}

type DeviceConnection struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Source       TelemetrySource `json:"source"`
	ConnectedAt  time.Time       `json:"connected_at"`
	LastSyncAt   *time.Time      `json:"last_sync_at,omitempty"`
	Status       string          `json:"status"` // "connected", "disconnected", "error"
	AccessToken  string          `json:"-"`      // Encrypted in real app
	RefreshToken string          `json:"-"`      // Encrypted in real app
}
