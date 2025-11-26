package domain

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the category of notification
type NotificationType string

const (
	NotificationTypeStreak       NotificationType = "streak"
	NotificationTypeVault        NotificationType = "vault"
	NotificationTypeSocial       NotificationType = "social"
	NotificationTypeGamification NotificationType = "gamification"
	NotificationTypeReferral     NotificationType = "referral"
	NotificationTypeFasting      NotificationType = "fasting"
)

// Notification represents a push notification sent to a user
type Notification struct {
	ID     uuid.UUID         `json:"id"`
	UserID uuid.UUID         `json:"user_id"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Type   NotificationType  `json:"type"`
	Data   map[string]string `json:"data,omitempty"` // Additional payload data
	SentAt time.Time         `json:"sent_at"`
	ReadAt *time.Time        `json:"read_at,omitempty"`
}

// FCMToken represents a Firebase Cloud Messaging token for a user's device
type FCMToken struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Token      string    `json:"token"`
	DeviceType string    `json:"device_type"` // "web", "android", "ios"
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}
