package domain

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	MealType    string    `json:"meal_type"` // breakfast, lunch, dinner, snack
	Image       string    `json:"image"`     // Base64 string
	Description string    `json:"description"`
	LoggedAt    time.Time `json:"logged_at"`
	Calories    int       `json:"calories,omitempty"`
	Analysis    string    `json:"analysis"`     // DeepSeek analysis text
	IsKeto      bool      `json:"is_keto"`      // Parsed from analysis
	IsAuthentic bool      `json:"is_authentic"` // Parsed from analysis
}
