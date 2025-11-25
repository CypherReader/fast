package domain

import (
	"time"

	"github.com/google/uuid"
)

type SocialPost struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	Type      string    `json:"type"` // streak, meal
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}
