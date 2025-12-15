package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProgressService struct {
	repo ports.ProgressRepository
}

func NewProgressService(repo ports.ProgressRepository) *ProgressService {
	return &ProgressService{repo: repo}
}

func (s *ProgressService) LogWeight(ctx context.Context, userID uuid.UUID, weight float64, unit string) (*domain.WeightLog, error) {
	// Input validation
	if weight <= 0 {
		return nil, errors.New("weight must be a positive number")
	}

	// Reasonable bounds to prevent extreme values
	const maxWeightKg = 500.0 // ~1100 lbs
	const minWeightKg = 20.0  // ~44 lbs

	var weightLbs, weightKg float64

	switch unit {
	case "kg":
		if weight < minWeightKg || weight > maxWeightKg {
			return nil, fmt.Errorf("weight must be between %.1f and %.1f kg", minWeightKg, maxWeightKg)
		}
		weightKg = weight
		weightLbs = weight * 2.20462
	case "lbs", "lb":
		minWeightLbs := minWeightKg * 2.20462
		maxWeightLbs := maxWeightKg * 2.20462
		if weight < minWeightLbs || weight > maxWeightLbs {
			return nil, fmt.Errorf("weight must be between %.1f and %.1f lbs", minWeightLbs, maxWeightLbs)
		}
		weightLbs = weight
		weightKg = weight / 2.20462
	default:
		return nil, fmt.Errorf("invalid unit: must be 'kg' or 'lbs'")
	}

	log := &domain.WeightLog{
		ID:        uuid.New(),
		UserID:    userID,
		WeightLbs: weightLbs,
		WeightKg:  weightKg,
		LoggedAt:  time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.repo.SaveWeightLog(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func (s *ProgressService) GetWeightHistory(ctx context.Context, userID uuid.UUID, days int) ([]domain.WeightLog, error) {
	return s.repo.GetWeightHistory(ctx, userID, days)
}

func (s *ProgressService) LogHydration(ctx context.Context, userID uuid.UUID, amount float64, unit string) (*domain.HydrationLog, error) {
	// Validate amount is within reasonable bounds
	if amount < 0 || amount > 100 {
		return nil, errors.New("hydration amount must be between 0 and 100 glasses")
	}

	glasses := int(amount)

	log := &domain.HydrationLog{
		ID:           uuid.New(),
		UserID:       userID,
		GlassesCount: glasses,
		LoggedDate:   time.Now(),
		CreatedAt:    time.Now(),
	}

	if err := s.repo.SaveHydrationLog(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func (s *ProgressService) GetDailyHydration(ctx context.Context, userID uuid.UUID) (*domain.HydrationLog, error) {
	return s.repo.GetHydrationLog(ctx, userID, time.Now())
}
