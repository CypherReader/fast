package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type MealService struct {
	repo   ports.MealRepository
	cortex ports.CortexService
}

func NewMealService(repo ports.MealRepository, cortex ports.CortexService) *MealService {
	return &MealService{
		repo:   repo,
		cortex: cortex,
	}
}

func (s *MealService) LogMeal(ctx context.Context, userID uuid.UUID, image, description string) (*domain.Meal, error) {
	// 1. Analyze Meal
	analysis, isAuthentic, isKeto, err := s.cortex.AnalyzeMeal(ctx, image, description)
	if err != nil {
		// Fallback if analysis fails (or just log error and proceed with defaults)
		// For now, we'll proceed but note the error in analysis text
		analysis = "Analysis failed: " + err.Error()
		isAuthentic = true
		isKeto = true
	}

	meal := &domain.Meal{
		ID:          uuid.New(),
		UserID:      userID,
		Image:       image,
		Description: description,
		LoggedAt:    time.Now(),
		Calories:    0, // Placeholder
		Analysis:    analysis,
		IsAuthentic: isAuthentic,
		IsKeto:      isKeto,
	}
	if err := s.repo.Save(ctx, meal); err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *MealService) GetMeals(ctx context.Context, userID uuid.UUID) ([]domain.Meal, error) {
	return s.repo.FindByUserID(ctx, userID)
}
