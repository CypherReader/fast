package domain

import (
	"time"

	"github.com/google/uuid"
)

type VaultParticipation struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	MonthStart      time.Time  `json:"month_start"`
	MonthEnd        time.Time  `json:"month_end"`
	DepositAmount   float64    `json:"deposit_amount"`
	FastsCompleted  int        `json:"fasts_completed"`
	AmountRecovered float64    `json:"amount_recovered"`
	RefundProcessed bool       `json:"refund_processed"`
	RefundDate      *time.Time `json:"refund_date,omitempty"`
	OptedIn         bool       `json:"opted_in"`
	ForfeitedAmount float64    `json:"forfeited_amount"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
