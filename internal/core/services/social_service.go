package services

import (
	"context"
	"encoding/json"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type SocialService struct {
	repo     ports.SocialRepository
	userRepo ports.UserRepository
}

func NewSocialService(repo ports.SocialRepository, userRepo ports.UserRepository) *SocialService {
	return &SocialService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *SocialService) LogEvent(ctx context.Context, userID uuid.UUID, eventType domain.SocialEventType, data interface{}) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event := &domain.SocialEvent{
		ID:        uuid.New(),
		UserID:    userID,
		UserName:  user.Email, // Ideally use a display name if available
		Type:      eventType,
		Data:      jsonData,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveEvent(ctx, event)
}

func (s *SocialService) GetFeed(ctx context.Context) ([]domain.SocialEvent, error) {
	// For MVP, just return global feed
	// In future, we can check if user is in a tribe and call GetTribeFeed
	return s.repo.GetFeed(ctx, 20)
}
