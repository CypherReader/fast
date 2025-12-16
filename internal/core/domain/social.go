package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeType string

const (
	ChallengeTypeFasting   ChallengeType = "fasting"
	ChallengeTypeSteps     ChallengeType = "steps"
	ChallengeTypeHydration ChallengeType = "hydration"
)

type EventType string

const (
	EventFastCompleted EventType = "fast_completed"
	EventTribeJoined   EventType = "tribe_joined"
	EventChallengeWon  EventType = "challenge_won"
)

type SocialEvent struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	EventType EventType `json:"event_type"`
	Data      string    `json:"data"` // JSON or text description
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
}

type FriendNetwork struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	FriendID    uuid.UUID `json:"friend_id"`
	Status      string    `json:"status"` // pending, accepted, blocked
	ConnectedAt time.Time `json:"connected_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Tribe represents a social group for fasting accountability
type Tribe struct {
	ID                string     `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	Slug              string     `json:"slug" db:"slug"`
	Description       string     `json:"description" db:"description"`
	AvatarURL         *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	CoverPhotoURL     *string    `json:"cover_photo_url,omitempty" db:"cover_photo_url"`
	CreatorID         string     `json:"creator_id" db:"creator_id"`
	FastingSchedule   string     `json:"fasting_schedule" db:"fasting_schedule"` // "16:8", "18:6", "omad", "custom"
	PrimaryGoal       string     `json:"primary_goal" db:"primary_goal"`         // "weight_loss", "metabolic_health", etc.
	Category          []byte     `json:"category" db:"category"`                 // JSON array of category tags
	Privacy           string     `json:"privacy" db:"privacy"`                   // "public", "private", "invite_only"
	Rules             *string    `json:"rules,omitempty" db:"rules"`
	MemberCount       int        `json:"member_count" db:"member_count"`
	ActiveMemberCount int        `json:"active_member_count" db:"active_member_count"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

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

type TribePool struct {
	ID                uuid.UUID `json:"id"`
	TribeID           uuid.UUID `json:"tribe_id"`
	MonthStart        time.Time `json:"month_start"`
	MonthEnd          time.Time `json:"month_end"`
	TotalPot          float64   `json:"total_pot"`
	ParticipantCount  int       `json:"participant_count"`
	Status            string    `json:"status"` // 'open', 'active', 'completed'
	FirstPlacePayout  float64   `json:"first_place_payout"`
	SecondPlacePayout float64   `json:"second_place_payout"`
	ThirdPlacePayout  float64   `json:"third_place_payout"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TribePoolParticipant struct {
	ID              uuid.UUID `json:"id"`
	PoolID          uuid.UUID `json:"pool_id"`
	UserID          uuid.UUID `json:"user_id"`
	DepositAmount   float64   `json:"deposit_amount"`
	FastsCompleted  int       `json:"fasts_completed"`
	FinalRank       int       `json:"final_rank"`
	PayoutAmount    float64   `json:"payout_amount"`
	PayoutProcessed bool      `json:"payout_processed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type FriendChallenge struct {
	ID              uuid.UUID     `json:"id"`
	CreatorID       uuid.UUID     `json:"creator_id"`
	Name            string        `json:"name"`
	ChallengeType   ChallengeType `json:"challenge_type"`
	Goal            int           `json:"goal"`
	StartDate       time.Time     `json:"start_date"`
	EndDate         time.Time     `json:"end_date"`
	Status          string        `json:"status"` // active, completed
	WinnerID        *uuid.UUID    `json:"winner_id,omitempty"`
	PayoutAmount    float64       `json:"payout_amount"`
	PayoutProcessed bool          `json:"payout_processed"`
	PayoutDate      *time.Time    `json:"payout_date,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	// Legacy fields for compatibility if needed, or remove if fully migrating
	ChallengerID      uuid.UUID `json:"challenger_id,omitempty"`
	ChallengedID      uuid.UUID `json:"challenged_id,omitempty"`
	MonthStart        time.Time `json:"month_start,omitempty"`
	MonthEnd          time.Time `json:"month_end,omitempty"`
	PotAmount         float64   `json:"pot_amount,omitempty"`
	ChallengerDeposit float64   `json:"challenger_deposit,omitempty"`
	ChallengedDeposit float64   `json:"challenged_deposit,omitempty"`
	ChallengerFasts   int       `json:"challenger_fasts,omitempty"`
	ChallengedFasts   int       `json:"challenged_fasts,omitempty"`
}
