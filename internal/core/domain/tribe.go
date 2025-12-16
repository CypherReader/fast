package domain

import (
	"encoding/json"
	"time"
)

// Tribe represents a social group for fasting accountability
type Tribe struct {
	ID                string          `json:"id" db:"id"`
	Name              string          `json:"name" db:"name"`
	Slug              string          `json:"slug" db:"slug"`
	Description       string          `json:"description" db:"description"`
	AvatarURL         *string         `json:"avatar_url,omitempty" db:"avatar_url"`
	CoverPhotoURL     *string         `json:"cover_photo_url,omitempty" db:"cover_photo_url"`
	CreatorID         string          `json:"creator_id" db:"creator_id"`
	FastingSchedule   string          `json:"fasting_schedule" db:"fasting_schedule"` // "16:8", "18:6", "omad", "custom"
	PrimaryGoal       string          `json:"primary_goal" db:"primary_goal"`         // "weight_loss", "metabolic_health", etc.
	Category          json.RawMessage `json:"category" db:"category"`                 // JSON array of category tags
	Privacy           string          `json:"privacy" db:"privacy"`                   // "public", "private", "invite_only"
	Rules             *string         `json:"rules,omitempty" db:"rules"`
	MemberCount       int             `json:"member_count" db:"member_count"`
	ActiveMemberCount int             `json:"active_member_count" db:"active_member_count"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`

	// Computed fields (not in database)
	IsJoined bool   `json:"is_joined" db:"-"`           // Whether current user is a member
	UserRole string `json:"user_role,omitempty" db:"-"` // Current user's role if member
}

// TribeMembership represents a user's membership in a tribe
type TribeMembership struct {
	ID                   string     `json:"id" db:"id"`
	TribeID              string     `json:"tribe_id" db:"tribe_id"`
	UserID               string     `json:"user_id" db:"user_id"`
	Role                 string     `json:"role" db:"role"`     // "creator", "moderator", "member"
	Status               string     `json:"status" db:"status"` // "active", "pending", "left"
	JoinedAt             time.Time  `json:"joined_at" db:"joined_at"`
	LeftAt               *time.Time `json:"left_at,omitempty" db:"left_at"`
	NotificationsEnabled bool       `json:"notifications_enabled" db:"notifications_enabled"`
}

// TribeMember extends membership with user information
type TribeMember struct {
	TribeMembership
	UserName   string `json:"user_name" db:"user_name"`
	UserAvatar string `json:"user_avatar" db:"user_avatar"`
	UserStreak int    `json:"user_streak" db:"user_streak"` // Current fasting streak
}

// CreateTribeRequest represents the request to create a new tribe
type CreateTribeRequest struct {
	Name            string   `json:"name" validate:"required,min=3,max=50"`
	Description     string   `json:"description" validate:"required,min=10,max=500"`
	FastingSchedule string   `json:"fasting_schedule" validate:"required,oneof='16:8' '18:6' 'omad' 'custom'"`
	PrimaryGoal     string   `json:"primary_goal" validate:"required"`
	Category        []string `json:"category,omitempty"`
	Privacy         string   `json:"privacy" validate:"required,oneof=public private invite_only"`
	Rules           string   `json:"rules,omitempty" validate:"max=1000"`
	AvatarURL       string   `json:"avatar_url,omitempty"`
	CoverPhotoURL   string   `json:"cover_photo_url,omitempty"`
}

// UpdateTribeRequest represents the request to update tribe information
type UpdateTribeRequest struct {
	Description   *string   `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Category      *[]string `json:"category,omitempty"`
	Privacy       *string   `json:"privacy,omitempty" validate:"omitempty,oneof=public private invite_only"`
	Rules         *string   `json:"rules,omitempty" validate:"omitempty,max=1000"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	CoverPhotoURL *string   `json:"cover_photo_url,omitempty"`
}

// ListTribesQuery represents query parameters for listing tribes
type ListTribesQuery struct {
	Search          string `form:"search"`
	FastingSchedule string `form:"fasting_schedule"`
	PrimaryGoal     string `form:"primary_goal"`
	Privacy         string `form:"privacy"`
	SortBy          string `form:"sort_by"` // "newest", "popular", "active", "members"
	Limit           int    `form:"limit" validate:"min=1,max=100"`
	Offset          int    `form:"offset" validate:"min=0"`
}

// TribeStats represents statistics for a tribe
type TribeStats struct {
	TribeID              string  `json:"tribe_id"`
	TotalFasts           int     `json:"total_fasts"`
	TotalFastingHours    float64 `json:"total_fasting_hours"`
	AverageMemberStreak  float64 `json:"average_member_streak"`
	WeeklyGrowthPercent  float64 `json:"weekly_growth_percent"`
	ActiveMembersPercent float64 `json:"active_members_percent"`
}
