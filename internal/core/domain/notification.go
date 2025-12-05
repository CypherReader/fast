package domain

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	Link      string    `json:"link,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FCMToken struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Token      string    `json:"token"`
	DeviceType string    `json:"device_type"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type NotificationType string

const (
	NotificationTypeFastComplete  NotificationType = "fast_complete"
	NotificationTypeVaultRefund   NotificationType = "vault_refund"
	NotificationTypeChallenge     NotificationType = "challenge_invite"
	NotificationTypePotWinner     NotificationType = "pot_winner"
	NotificationTypeStreakWarning NotificationType = "streak_warning"
	NotificationTypeFriendInvite  NotificationType = "friend_invite"
)
