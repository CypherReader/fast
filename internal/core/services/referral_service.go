package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ReferralService struct {
	referralRepo ports.ReferralRepository
	userRepo     ports.UserRepository
	vaultService ports.VaultService
}

func NewReferralService(referralRepo ports.ReferralRepository, userRepo ports.UserRepository, vaultService ports.VaultService) *ReferralService {
	return &ReferralService{
		referralRepo: referralRepo,
		userRepo:     userRepo,
		vaultService: vaultService,
	}
}

func (s *ReferralService) GenerateReferralCode(ctx context.Context, userID uuid.UUID) (string, error) {
	// Check if user already has a code
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user.ReferralCode != "" {
		return user.ReferralCode, nil
	}

	// Generate new code
	code, err := generateRandomCode(6)
	if err != nil {
		return "", err
	}

	user.ReferralCode = code
	if err := s.userRepo.Save(ctx, user); err != nil {
		return "", err
	}

	return code, nil
}

func (s *ReferralService) GetReferralCode(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user.ReferralCode == "" {
		return s.GenerateReferralCode(ctx, userID)
	}
	return user.ReferralCode, nil
}

func (s *ReferralService) TrackReferral(ctx context.Context, referrerCode string, refereeID uuid.UUID) error {
	if referrerCode == "" {
		return nil
	}

	// Find referrer by code
	// Note: We need a FindByReferralCode method in UserRepository.
	// For now, let's assume we might need to add it or iterate (inefficient).
	// Better: Add FindByReferralCode to UserRepository.
	referrer, err := s.userRepo.FindByReferralCode(ctx, referrerCode)
	if err != nil {
		// If not found, just ignore (invalid code)
		return nil
	}

	if referrer.ID == refereeID {
		return errors.New("cannot refer yourself")
	}

	// Check if referral already exists
	existing, _ := s.referralRepo.FindByRefereeID(ctx, refereeID)
	if existing != nil {
		return errors.New("user already referred")
	}

	referral := &domain.Referral{
		ID:          uuid.New(),
		ReferrerID:  referrer.ID,
		RefereeID:   refereeID,
		Status:      domain.ReferralStatusPending,
		RewardValue: 5.00, // $5.00 reward
		CreatedAt:   time.Now(),
	}

	return s.referralRepo.Save(ctx, referral)
}

func (s *ReferralService) CompleteReferral(ctx context.Context, refereeID uuid.UUID) error {
	referral, err := s.referralRepo.FindByRefereeID(ctx, refereeID)
	if err != nil {
		return nil // No referral found, nothing to do
	}

	if referral.Status == domain.ReferralStatusCompleted {
		return nil // Already completed
	}

	// Award referrer
	referrer, err := s.userRepo.FindByID(ctx, referral.ReferrerID)
	if err != nil {
		return err
	}
	s.vaultService.AddDailyEarnings(ctx, referrer, referral.RewardValue)

	// Award referee
	referee, err := s.userRepo.FindByID(ctx, refereeID)
	if err != nil {
		return err
	}
	s.vaultService.AddDailyEarnings(ctx, referee, referral.RewardValue)

	// Mark as completed
	now := time.Now()
	referral.Status = domain.ReferralStatusCompleted
	referral.CompletedAt = &now

	return s.referralRepo.Update(ctx, referral)
}

func (s *ReferralService) GetReferralStats(ctx context.Context, userID uuid.UUID) (float64, int, error) {
	referrals, err := s.referralRepo.FindByReferrerID(ctx, userID)
	if err != nil {
		return 0, 0, err
	}

	var totalEarned float64
	var count int

	for _, r := range referrals {
		if r.Status == domain.ReferralStatusCompleted {
			totalEarned += r.RewardValue
			count++
		}
	}

	return totalEarned, count, nil
}

func generateRandomCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(base64.URLEncoding.EncodeToString(bytes)[:length]), nil
}
