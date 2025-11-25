package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type ActivityService struct {
	repo ports.ActivityRepository
}

func NewActivityService(repo ports.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) SyncActivity(ctx context.Context, userID uuid.UUID, activity domain.Activity) error {
	activity.UserID = userID.String()
	// In a real app, we might check for duplicates or merge data here
	return s.repo.Save(ctx, &activity)
}

func (s *ActivityService) GetActivities(ctx context.Context, userID uuid.UUID) ([]domain.Activity, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *ActivityService) GetActivity(ctx context.Context, activityID string) (*domain.Activity, error) {
	return s.repo.FindByID(ctx, activityID)
}
