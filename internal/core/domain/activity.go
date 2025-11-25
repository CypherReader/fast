package domain

import "time"

type ActivityType string

const (
	ActivityTypeWalk    ActivityType = "WALK"
	ActivityTypeRun     ActivityType = "RUN"
	ActivityTypeCycle   ActivityType = "CYCLE"
	ActivityTypeWorkout ActivityType = "WORKOUT"
)

type RoutePoint struct {
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Timestamp time.Time `json:"timestamp"`
}

type Activity struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	Type      ActivityType  `json:"type"`
	StartTime time.Time     `json:"start_time"`
	Duration  time.Duration `json:"duration"` // stored in nanoseconds, serialized as string usually, but we'll handle it
	Steps     int           `json:"steps"`
	Distance  float64       `json:"distance"` // in meters
	Calories  int           `json:"calories"`
	Route     []RoutePoint  `json:"route,omitempty"`
	Source    string        `json:"source"` // e.g., "GARMIN", "APPLE", "MANUAL"
}
