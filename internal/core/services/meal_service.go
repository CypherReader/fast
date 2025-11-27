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

func (s *MealService) LogMeal(ctx context.Context, userID uuid.UUID, name string, calories int, mealType string, image, description string) (*domain.Meal, error) {
	var analysis string
	var isAuthentic, isKeto bool
	var err error

	// Only analyze if image or description is provided AND we don't have manual data (or we want to augment it)
	// For now, if manual data is provided, we skip analysis to save tokens/time, unless explicitly requested?
	// Let's say if image is provided, we always analyze.
	if image != "" || (description != "" && name == "") {
		// 1. Analyze Meal
		analysis, isAuthentic, isKeto, err = s.cortex.AnalyzeMeal(ctx, image, description)
		if err != nil {
			// Fallback
			analysis = "Analysis failed: " + err.Error()
			isAuthentic = true // Default to optimistic
			isKeto = true
		}
	}

	// If name is empty but we have description, use description as name
	if name == "" {
		if len(description) > 50 {
			name = description[:47] + "..."
		} else {
			name = description
		}
	}
	if name == "" {
		name = "Unnamed Meal"
	}

	meal := &domain.Meal{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		MealType:    mealType,
		Image:       image,
		Description: description,
		LoggedAt:    time.Now(),
		Calories:    calories,
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
