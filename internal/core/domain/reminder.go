package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReminderType represents the type of smart reminder
type ReminderType string

const (
	ReminderTypeFastStart ReminderType = "fast_start"
	ReminderTypeFastEnd   ReminderType = "fast_end"
	ReminderTypeHydration ReminderType = "hydration"
	ReminderTypeWeekly    ReminderType = "weekly_checkin"
)

// ScheduledReminder represents a scheduled reminder in the system
type ScheduledReminder struct {
	ID           uuid.UUID              `json:"id"`
	UserID       uuid.UUID              `json:"user_id"`
	ReminderType ReminderType           `json:"reminder_type"`
	ScheduledAt  time.Time              `json:"scheduled_at"`
	Sent         bool                   `json:"sent"`
	Message      string                 `json:"message,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// ReminderSettings represents user preferences for smart reminders
type ReminderSettings struct {
	UserID                   uuid.UUID `json:"user_id"`
	ReminderFastStart        bool      `json:"reminder_fast_start"`
	ReminderFastEnd          bool      `json:"reminder_fast_end"`
	ReminderHydration        bool      `json:"reminder_hydration"`
	PreferredFastStartHour   int       `json:"preferred_fast_start_hour"` // 0-23
	HydrationIntervalMinutes int       `json:"hydration_interval_minutes"`
}

// OptimalFastingWindow represents AI-suggested fasting times
type OptimalFastingWindow struct {
	SuggestedStartTime time.Time `json:"suggested_start_time"`
	SuggestedEndTime   time.Time `json:"suggested_end_time"`
	SuggestedDuration  int       `json:"suggested_duration_hours"`
	Reasoning          string    `json:"reasoning"`
	ConfidenceScore    float64   `json:"confidence_score"` // 0.0-1.0
}
