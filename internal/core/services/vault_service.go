package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"math"
)

const (
	MonthlyCharge = 30.0
	BaseFee       = 10.0
	VaultDeposit  = 20.0
	DailyMax      = 2.0
)

type VaultService struct {
	userRepo       ports.UserRepository
	paymentGateway ports.PaymentGateway
}

func NewVaultService(userRepo ports.UserRepository, paymentGateway ports.PaymentGateway) *VaultService {
	return &VaultService{
		userRepo:       userRepo,
		paymentGateway: paymentGateway,
	}
}

// CalculateDailyEarning calculates how much a user earns based on their discipline score
func (s *VaultService) CalculateDailyEarning(disciplineIndex int) float64 {
	// Formula: $2.00 * (DisciplineIndex / 100)
	earning := DailyMax * (float64(disciplineIndex) / 100.0)
	if earning > DailyMax {
		earning = DailyMax
	}
	return earning
}

// ProcessDailyEarnings iterates over all users and updates their earnings
func (s *VaultService) ProcessDailyEarnings(ctx context.Context) error {
	// In a real app, we would stream users or process in batches
	// For MVP, we'll just fetch all users (assuming low volume)
	// TODO: Add GetAllUsers to UserRepository
	return nil
}

// CalculateVaultStatus returns the current vault status for the user
func (s *VaultService) CalculateVaultStatus(user *domain.User) (deposit float64, earned float64, potentialRefund float64) {
	if !user.IsVaultMember() {
		return 0, 0, 0
	}

	deposit = user.VaultDeposit
	earned = user.EarnedRefund

	// Refund cannot exceed deposit
	potentialRefund = math.Min(earned, deposit)
	return
}

func (s *VaultService) AddDailyEarnings(ctx context.Context, user *domain.User, amount float64) {
	if !user.IsVaultMember() {
		return
	}

	user.EarnedRefund += amount

	// Hard cap at vault deposit
	if user.EarnedRefund > user.VaultDeposit {
		user.EarnedRefund = user.VaultDeposit
	}
}

// Deprecated: Kept for backward compatibility until full migration
func (s *VaultService) CalculatePrice(ctx context.Context, user *domain.User) float64 {
	_, _, potentialRefund := s.CalculateVaultStatus(user)
	return MonthlyCharge - potentialRefund
}

func (s *VaultService) UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast bool, verifiedKetosis bool) {
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

func (s *VaultService) ProcessMonthlyRefunds(ctx context.Context) error {
	// Logic to process refunds at the end of the billing cycle
	return nil
}
