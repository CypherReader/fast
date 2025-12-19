package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ============== UPDATE PROFILE TESTS ==============

func TestOnboardingService_UpdateProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:   userID,
		Name: "Old Name",
		Goal: "lose_weight",
	}

	newName := "New Name"
	newGoal := "maintain"
	profile := domain.UserProfileUpdate{
		Name: &newName,
		Goal: &newGoal,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	result, err := service.UpdateProfile(ctx, userID, profile)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "maintain", result.Goal)
}

func TestOnboardingService_UpdateProfile_PartialUpdate(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:          userID,
		Name:        "Old Name",
		Goal:        "lose_weight",
		FastingPlan: "16:8",
	}

	newPlan := "18:6"
	profile := domain.UserProfileUpdate{
		FastingPlan: &newPlan,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	result, err := service.UpdateProfile(ctx, userID, profile)

	assert.NoError(t, err)
	assert.Equal(t, "Old Name", result.Name)    // Unchanged
	assert.Equal(t, "18:6", result.FastingPlan) // Updated
}

func TestOnboardingService_UpdateProfile_AllFields(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID}

	name := "Test User"
	goal := "build_muscle"
	plan := "20:4"
	sex := "male"
	height := 180.0
	currentWeight := 185.0
	targetWeight := 175.0
	timezone := "America/New_York"
	units := "imperial"

	profile := domain.UserProfileUpdate{
		Name:             &name,
		Goal:             &goal,
		FastingPlan:      &plan,
		Sex:              &sex,
		HeightCm:         &height,
		CurrentWeightLbs: &currentWeight,
		TargetWeightLbs:  &targetWeight,
		Timezone:         &timezone,
		Units:            &units,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	result, err := service.UpdateProfile(ctx, userID, profile)

	assert.NoError(t, err)
	assert.Equal(t, name, result.Name)
	assert.Equal(t, goal, result.Goal)
	assert.Equal(t, plan, result.FastingPlan)
	assert.Equal(t, sex, result.Sex)
	assert.Equal(t, height, result.HeightCm)
}

func TestOnboardingService_UpdateProfile_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("not found"))

	result, err := service.UpdateProfile(ctx, userID, domain.UserProfileUpdate{})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestOnboardingService_UpdateProfile_SaveError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{ID: userID}
	name := "New Name"
	profile := domain.UserProfileUpdate{Name: &name}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(errors.New("save error"))

	result, err := service.UpdateProfile(ctx, userID, profile)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ============== COMPLETE ONBOARDING TESTS ==============

func TestOnboardingService_CompleteOnboarding_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:                  userID,
		OnboardingCompleted: false,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	result, err := service.CompleteOnboarding(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.OnboardingCompleted)
}

func TestOnboardingService_CompleteOnboarding_AlreadyComplete(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	user := &domain.User{
		ID:                  userID,
		OnboardingCompleted: true,
	}

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil)
	mockUserRepo.On("Save", ctx, user).Return(nil)

	result, err := service.CompleteOnboarding(ctx, userID)

	assert.NoError(t, err)
	assert.True(t, result.OnboardingCompleted)
}

func TestOnboardingService_CompleteOnboarding_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	service := NewOnboardingService(mockUserRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("not found"))

	result, err := service.CompleteOnboarding(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
