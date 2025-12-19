package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMealRepository is a mock of ports.MealRepository
type MockMealRepository struct {
	mock.Mock
}

func (m *MockMealRepository) Save(ctx context.Context, meal *domain.Meal) error {
	args := m.Called(ctx, meal)
	return args.Error(0)
}

func (m *MockMealRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Meal, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Meal), args.Error(1)
}

func (m *MockMealRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Meal, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Meal), args.Error(1)
}

// MockCortexServiceForMeal is a mock of CortexService for meal tests
type MockCortexServiceForMeal struct {
	mock.Mock
}

func (m *MockCortexServiceForMeal) Chat(ctx context.Context, userID uuid.UUID, message string) (string, error) {
	args := m.Called(ctx, userID, message)
	return args.String(0), args.Error(1)
}

func (m *MockCortexServiceForMeal) GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error) {
	args := m.Called(ctx, userID, fastingHours)
	return args.String(0), args.Error(1)
}

func (m *MockCortexServiceForMeal) AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error) {
	args := m.Called(ctx, imageBase64, description)
	return args.String(0), args.Bool(1), args.Bool(2), args.Error(3)
}

func (m *MockCortexServiceForMeal) GetCravingHelp(ctx context.Context, userID uuid.UUID, cravingDescription string) (interface{}, error) {
	args := m.Called(ctx, userID, cravingDescription)
	return args.Get(0), args.Error(1)
}

// ============== LOG MEAL TESTS ==============

func TestMealService_LogMeal_Success(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(nil)

	meal, err := service.LogMeal(ctx, userID, "Grilled Chicken", 450, "lunch", "", "Healthy chicken with veggies")

	assert.NoError(t, err)
	assert.NotNil(t, meal)
	assert.Equal(t, "Grilled Chicken", meal.Name)
	assert.Equal(t, 450, meal.Calories)
	assert.Equal(t, "lunch", meal.MealType)
}

func TestMealService_LogMeal_WithImage(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockCortex.On("AnalyzeMeal", ctx, "base64image", "").Return("Healthy keto meal", true, true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(nil)

	meal, err := service.LogMeal(ctx, userID, "My Meal", 0, "dinner", "base64image", "")

	assert.NoError(t, err)
	assert.NotNil(t, meal)
	assert.True(t, meal.IsKeto)
	assert.True(t, meal.IsAuthentic)
}

func TestMealService_LogMeal_AnalysisFails(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockCortex.On("AnalyzeMeal", ctx, "badimage", "").Return("", false, false, errors.New("analysis failed"))
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(nil)

	meal, err := service.LogMeal(ctx, userID, "My Meal", 0, "dinner", "badimage", "")

	assert.NoError(t, err)
	assert.NotNil(t, meal)
	// On failure, defaults to optimistic
	assert.True(t, meal.IsAuthentic)
	assert.True(t, meal.IsKeto)
}

func TestMealService_LogMeal_NoName(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockCortex.On("AnalyzeMeal", ctx, "", "Salmon with asparagus").Return("Great meal!", true, true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(nil)

	meal, err := service.LogMeal(ctx, userID, "", 0, "dinner", "", "Salmon with asparagus")

	assert.NoError(t, err)
	assert.NotNil(t, meal)
	assert.Equal(t, "Salmon with asparagus", meal.Name) // Uses description as name
}

func TestMealService_LogMeal_NoNameNoDescription(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(nil)

	meal, err := service.LogMeal(ctx, userID, "", 0, "snack", "", "")

	assert.NoError(t, err)
	assert.NotNil(t, meal)
	assert.Equal(t, "Unnamed Meal", meal.Name)
}

func TestMealService_LogMeal_SaveError(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Meal")).Return(errors.New("db error"))

	meal, err := service.LogMeal(ctx, userID, "Test", 100, "lunch", "", "")

	assert.Error(t, err)
	assert.Nil(t, meal)
}

// ============== GET MEALS TESTS ==============

func TestMealService_GetMeals_Success(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	meals := []domain.Meal{
		{ID: uuid.New(), Name: "Breakfast"},
		{ID: uuid.New(), Name: "Lunch"},
	}

	mockRepo.On("FindByUserID", ctx, userID).Return(meals, nil)

	result, err := service.GetMeals(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestMealService_GetMeals_Empty(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindByUserID", ctx, userID).Return([]domain.Meal{}, nil)

	result, err := service.GetMeals(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestMealService_GetMeals_Error(t *testing.T) {
	mockRepo := new(MockMealRepository)
	mockCortex := new(MockCortexServiceForMeal)
	service := NewMealService(mockRepo, mockCortex)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("db error"))

	result, err := service.GetMeals(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
