package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type OnboardingService struct {
	userRepo ports.UserRepository
}

func NewOnboardingService(userRepo ports.UserRepository) *OnboardingService {
	return &OnboardingService{
		userRepo: userRepo,
	}
}

func (s *OnboardingService) UpdateProfile(ctx context.Context, userID uuid.UUID, profile domain.UserProfileUpdate) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if profile.Name != nil {
		user.Name = *profile.Name
	}
	if profile.Goal != nil {
		user.Goal = *profile.Goal
	}
	if profile.FastingPlan != nil {
		user.FastingPlan = *profile.FastingPlan
	}
	if profile.Sex != nil {
		user.Sex = *profile.Sex
	}
	if profile.HeightCm != nil {
		user.HeightCm = *profile.HeightCm
	}
	if profile.CurrentWeightLbs != nil {
		user.CurrentWeightLbs = *profile.CurrentWeightLbs
	}
	if profile.TargetWeightLbs != nil {
		user.TargetWeightLbs = *profile.TargetWeightLbs
	}
	if profile.Timezone != nil {
		user.Timezone = *profile.Timezone
	}
	if profile.Units != nil {
		user.Units = *profile.Units
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *OnboardingService) CompleteOnboarding(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.OnboardingCompleted = true
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
