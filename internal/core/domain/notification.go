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
	NotificationTypeSOSFlare      NotificationType = "sos_flare"     // Tribe member needs help
	NotificationTypeHypeReceived  NotificationType = "hype_received" // User received hype support
	NotificationTypeSOSResolved   NotificationType = "sos_resolved"  // SOS was resolved (rescued)
	NotificationTypeCortexBackup  NotificationType = "cortex_backup" // Cortex auto-responded to SOS
	// Smart Reminders
	NotificationTypeFastStartReminder NotificationType = "fast_start_reminder" // Time to start fasting
	NotificationTypeFastEndReminder   NotificationType = "fast_end_reminder"   // Fast completing soon
	NotificationTypeHydrationReminder NotificationType = "hydration_reminder"  // Drink water reminder
	NotificationTypeWeeklyCheckIn     NotificationType = "weekly_checkin"      // Weekly AI summary
)
