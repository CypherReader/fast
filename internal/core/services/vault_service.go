package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"math"
	"time"

	"github.com/google/uuid"
)

// V2 Pricing
const (
	MonthlySubscription = 4.99 // Standard subscription price for V2
)

// Deprecated vault constants (kept for backward compatibility)
// HIDDEN FOR V2: These are used by existing vault logic but not exposed to users
const (
	MonthlyCharge = 30.0 // Deprecated: old vault pricing
	BaseFee       = 10.0 // Deprecated: old vault fee
	VaultDeposit  = 20.0 // Deprecated: now using MonthlySubscription (4.99)
	DailyMax      = 2.0  // Deprecated: vault daily earning cap
)

type VaultService struct {
	userRepo       ports.UserRepository
	vaultRepo      ports.VaultRepository
	paymentGateway ports.PaymentGateway
}

func NewVaultService(userRepo ports.UserRepository, vaultRepo ports.VaultRepository, paymentGateway ports.PaymentGateway) *VaultService {
	return &VaultService{
		userRepo:       userRepo,
		vaultRepo:      vaultRepo,
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

	// 1. Update User Record
	user.EarnedRefund += amount
	if user.EarnedRefund > user.VaultDeposit {
		user.EarnedRefund = user.VaultDeposit
	}

	// 2. Update VaultParticipation Record
	// Find current month's participation
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	vault, err := s.vaultRepo.FindByUserIDAndMonth(ctx, user.ID, monthStart)
	if err == nil && vault != nil {
		vault.AmountRecovered += amount
		if vault.AmountRecovered > vault.DepositAmount {
			vault.AmountRecovered = vault.DepositAmount
		}
		vault.UpdatedAt = time.Now()
		_ = s.vaultRepo.Save(ctx, vault)
	} else {
		// Create if not exists (lazy creation)
		vault = &domain.VaultParticipation{
			ID:              uuid.New(),
			UserID:          user.ID,
			MonthStart:      monthStart,
			MonthEnd:        monthStart.AddDate(0, 1, 0).Add(-time.Second),
			DepositAmount:   user.VaultDeposit,
			AmountRecovered: amount, // Start with this amount
			OptedIn:         true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		_ = s.vaultRepo.Save(ctx, vault)
	}
}

// Deprecated: Kept for backward compatibility until full migration
func (s *VaultService) CalculatePrice(ctx context.Context, user *domain.User) float64 {
	_, _, potentialRefund := s.CalculateVaultStatus(user)
	return MonthlyCharge - potentialRefund
}

func (s *VaultService) UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast bool, verifiedKetosis bool) {
	if completedFast {
		user.DisciplineIndex += 1.0
	}
	if verifiedKetosis {
		user.DisciplineIndex += 2.0
	}

	if user.DisciplineIndex > 100 {
		user.DisciplineIndex = 100
	}
	if user.DisciplineIndex < 0 {
		user.DisciplineIndex = 0
	}

	// Calculate and add earnings if user is a vault member
	if user.IsVaultMember() {
		earning := s.CalculateDailyEarning(int(user.DisciplineIndex))
		s.AddDailyEarnings(ctx, user, earning)
	}
}

func (s *VaultService) ProcessMonthlyRefunds(ctx context.Context) error {
	// Logic to process refunds at the end of the billing cycle
	return nil
}

func (s *VaultService) GetCurrentParticipation(ctx context.Context, userID uuid.UUID) (*domain.VaultParticipation, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	return s.vaultRepo.FindByUserIDAndMonth(ctx, userID, monthStart)
}
