package domain

import (
	"time"

	"github.com/google/uuid"
)

type BadgeID string

const (
	BadgeFirstFast BadgeID = "first_fast"
	BadgeStreak3   BadgeID = "streak_3"
	BadgeStreak7   BadgeID = "streak_7"
	Badge100Hours  BadgeID = "100_hours"
	BadgeEarlyBird BadgeID = "early_bird" // Example: Finished fast before 10am
	BadgeNightOwl  BadgeID = "night_owl"  // Example: Started fast after 8pm
)

type Badge struct {
	ID          BadgeID `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"` // Emoji or URL
}

type UserBadge struct {
	UserID    uuid.UUID `json:"user_id"`
	BadgeID   BadgeID   `json:"badge_id"`
	EarnedAt  time.Time `json:"earned_at"`
	BadgeInfo *Badge    `json:"badge_info,omitempty"` // Populated on read
}

type UserStreak struct {
	UserID           uuid.UUID `json:"user_id"`
	CurrentStreak    int       `json:"current_streak"`
	LongestStreak    int       `json:"longest_streak"`
	LastActivityDate time.Time `json:"last_activity_date"`
}

var Badges = map[BadgeID]Badge{
	BadgeFirstFast: {ID: BadgeFirstFast, Name: "Beginner", Description: "Completed your first fast", Icon: "🌱"},
	BadgeStreak3:   {ID: BadgeStreak3, Name: "Consistency", Description: "3 day streak", Icon: "🔥"},
	BadgeStreak7:   {ID: BadgeStreak7, Name: "Dedicated", Description: "7 day streak", Icon: "🚀"},
	Badge100Hours:  {ID: Badge100Hours, Name: "Centurion", Description: "100 total fasting hours", Icon: "💯"},
}
