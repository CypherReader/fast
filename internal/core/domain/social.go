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

type Tribe struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
	MemberCount int       `json:"member_count"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TribeMember struct {
	ID       uuid.UUID `json:"id"`
	TribeID  uuid.UUID `json:"tribe_id"`
	UserID   uuid.UUID `json:"user_id"`
	IsAdmin  bool      `json:"is_admin"`
	JoinedAt time.Time `json:"joined_at"`
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
