package domain

import (
	"time"

	"github.com/google/uuid"
)

type FastingPlanType string

const (
	PlanBeginner FastingPlanType = "beginner"
	Plan168      FastingPlanType = "16_8"
	Plan186      FastingPlanType = "18_6"
	PlanOMAD     FastingPlanType = "omad"
	Plan24h      FastingPlanType = "24h"
	Plan36h      FastingPlanType = "36h"
	PlanExtended FastingPlanType = "extended"
)

type FastingStatus string

const (
	StatusActive    FastingStatus = "active"
	StatusCompleted FastingStatus = "completed"
	StatusCancelled FastingStatus = "cancelled"
)

type FastingSession struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	StartTime time.Time       `json:"start_time"`
	EndTime   *time.Time      `json:"end_time,omitempty"`
	GoalHours int             `json:"goal_hours"`
	PlanType  FastingPlanType `json:"plan_type"`
	Status    FastingStatus   `json:"status"`
}

func NewFastingSession(userID uuid.UUID, plan FastingPlanType, goalHours int, startTime time.Time) *FastingSession {
	if startTime.IsZero() {
		startTime = time.Now()
	}
	return &FastingSession{
		ID:        uuid.New(),
		UserID:    userID,
		StartTime: startTime,
		GoalHours: goalHours,
		PlanType:  plan,
		Status:    StatusActive,
	}
}
