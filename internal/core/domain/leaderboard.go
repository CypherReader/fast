package domain

import "github.com/google/uuid"

type LeaderboardEntry struct {
	UserID            uuid.UUID `json:"user_id"`
	UserName          string    `json:"user_name"`
	TotalFastingHours float64   `json:"total_fasting_hours"`
	DisciplineScore   float64   `json:"discipline_score"`
	Rank              int       `json:"rank"`
}
