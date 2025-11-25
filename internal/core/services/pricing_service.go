package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"math"
)

const (
	MonthlyCharge = 30.0
	BaseFee       = 10.0
	VaultDeposit  = 20.0
	DailyMax      = 2.0
)

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

// CalculateVaultStatus returns the current vault status for the user
func (s *PricingService) CalculateVaultStatus(user *domain.User) (deposit float64, earned float64, potentialRefund float64) {
	deposit = user.VaultDeposit
	earned = user.EarnedRefund

	// Refund cannot exceed deposit
	potentialRefund = math.Min(earned, deposit)
	return
}

func (s *PricingService) AddDailyEarnings(ctx context.Context, user *domain.User, amount float64) {
	// Logic to cap daily earnings would go here (requires tracking daily earnings in DB)
	// For MVP, we just add to the total earned refund

	user.EarnedRefund += amount

	// Hard cap at vault deposit
	if user.EarnedRefund > user.VaultDeposit {
		user.EarnedRefund = user.VaultDeposit
	}
}

// Deprecated: Kept for backward compatibility until full migration
func (s *PricingService) CalculatePrice(ctx context.Context, user *domain.User) float64 {
	_, _, potentialRefund := s.CalculateVaultStatus(user)
	return MonthlyCharge - potentialRefund
}

func (s *PricingService) UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast bool, verifiedKetosis bool) {
	if completedFast {
		user.DisciplineIndex += 1
	}
	if verifiedKetosis {
		user.DisciplineIndex += 2
	}

	if user.DisciplineIndex > 100 {
		user.DisciplineIndex = 100
	}
	if user.DisciplineIndex < 0 {
		user.DisciplineIndex = 0
	}
}
