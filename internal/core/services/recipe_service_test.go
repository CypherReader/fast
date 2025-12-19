package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRecipeRepository is a mock of ports.RecipeRepository
type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) FindAll(ctx context.Context) ([]domain.Recipe, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Recipe), args.Error(1)
}

// ============== GET RECIPES TESTS ==============

func TestRecipeService_GetRecipes_AllRecipes(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	recipes := []domain.Recipe{
		{Title: "Chicken Salad", Diet: domain.DietNormal},
		{Title: "Veggie Wrap", Diet: domain.DietVegetarian},
		{Title: "Tofu Bowl", Diet: domain.DietVegan},
	}

	mockRepo.On("FindAll", ctx).Return(recipes, nil)

	result, err := service.GetRecipes(ctx, "")

	assert.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestRecipeService_GetRecipes_FilterNormal(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	recipes := []domain.Recipe{
		{Title: "Chicken Salad", Diet: domain.DietNormal},
		{Title: "Veggie Wrap", Diet: domain.DietVegetarian},
		{Title: "Tofu Bowl", Diet: domain.DietVegan},
	}

	mockRepo.On("FindAll", ctx).Return(recipes, nil)

	result, err := service.GetRecipes(ctx, domain.DietNormal)

	assert.NoError(t, err)
	assert.Len(t, result, 3) // Normal shows all
}

func TestRecipeService_GetRecipes_FilterVegetarian(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	recipes := []domain.Recipe{
		{Title: "Chicken Salad", Diet: domain.DietNormal},
		{Title: "Veggie Wrap", Diet: domain.DietVegetarian},
		{Title: "Tofu Bowl", Diet: domain.DietVegan},
	}

	mockRepo.On("FindAll", ctx).Return(recipes, nil)

	result, err := service.GetRecipes(ctx, domain.DietVegetarian)

	assert.NoError(t, err)
	assert.Len(t, result, 2) // Vegetarian + Vegan
}

func TestRecipeService_GetRecipes_FilterVegan(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	recipes := []domain.Recipe{
		{Title: "Chicken Salad", Diet: domain.DietNormal},
		{Title: "Veggie Wrap", Diet: domain.DietVegetarian},
		{Title: "Tofu Bowl", Diet: domain.DietVegan},
	}

	mockRepo.On("FindAll", ctx).Return(recipes, nil)

	result, err := service.GetRecipes(ctx, domain.DietVegan)

	assert.NoError(t, err)
	assert.Len(t, result, 1) // Only Vegan
	assert.Equal(t, "Tofu Bowl", result[0].Title)
}

func TestRecipeService_GetRecipes_Empty(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return([]domain.Recipe{}, nil)

	result, err := service.GetRecipes(ctx, domain.DietVegan)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestRecipeService_GetRecipes_Error(t *testing.T) {
	mockRepo := new(MockRecipeRepository)
	service := NewRecipeService(mockRepo)
	ctx := context.Background()

	mockRepo.On("FindAll", ctx).Return(nil, errors.New("db error"))

	result, err := service.GetRecipes(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, result)
}
