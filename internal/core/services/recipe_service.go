package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
)

type RecipeService struct {
	repo ports.RecipeRepository
}

func NewRecipeService(repo ports.RecipeRepository) *RecipeService {
	return &RecipeService{repo: repo}
}

func (s *RecipeService) GetRecipes(ctx context.Context, diet domain.DietType) ([]domain.Recipe, error) {
	allRecipes, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// If diet is "normal" (or empty/all), return everything?
	// Or if "normal" is a specific category?
	// The requirement says "filter by diet".
	// Let's assume if diet is empty, return all.
	if diet == "" {
		return allRecipes, nil
	}

	var filtered []domain.Recipe
	for _, r := range allRecipes {
		// If user selects "Vegan", show Vegan.
		// If user selects "Vegetarian", show Vegetarian AND Vegan (since vegan is veg).
		// If user selects "Normal", show everything? Or just "Normal" tagged ones?
		// Let's stick to strict matching for now, or simple logic.

		// Logic:
		// Vegan -> Vegan only
		// Vegetarian -> Vegetarian + Vegan
		// Normal -> All (including meat)

		if diet == domain.DietNormal {
			filtered = append(filtered, r)
		} else if diet == domain.DietVegetarian {
			if r.Diet == domain.DietVegetarian || r.Diet == domain.DietVegan {
				filtered = append(filtered, r)
			}
		} else if diet == domain.DietVegan {
			if r.Diet == domain.DietVegan {
				filtered = append(filtered, r)
			}
		}
	}
	return filtered, nil
}
