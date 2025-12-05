package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
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
	var weightLbs, weightKg float64
	if unit == "kg" {
		weightKg = weight
		weightLbs = weight * 2.20462
	} else {
		weightLbs = weight
		weightKg = weight / 2.20462
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
	glasses := int(amount) // Simplified

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
