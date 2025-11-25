package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type SocialService struct {
	// In a real app, we'd have a SocialRepository
}

func NewSocialService() *SocialService {
	return &SocialService{}
}

func (s *SocialService) GetFeed(ctx context.Context) ([]domain.SocialPost, error) {
	// Mock data for now
	return []domain.SocialPost{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Username:  "FastingKing",
			Content:   "Just hit 36 hours! Feeling amazing.",
			Type:      "streak",
			Likes:     12,
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Username:  "KetoQueen",
			Content:   "Broke my fast with avocado and eggs.",
			Type:      "meal",
			Likes:     45,
			CreatedAt: time.Now().Add(-5 * time.Hour),
		},
	}, nil
}
