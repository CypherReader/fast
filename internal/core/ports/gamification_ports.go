package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type GamificationRepository interface {
	SaveUserBadge(ctx context.Context, badge *domain.UserBadge) error
	GetUserBadges(ctx context.Context, userID uuid.UUID) ([]domain.UserBadge, error)
	UpdateUserStreak(ctx context.Context, streak *domain.UserStreak) error
	GetUserStreak(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, error)
}

type GamificationService interface {
	CheckAndAwardBadges(ctx context.Context, userID uuid.UUID, eventType string, data interface{}) error
	UpdateStreak(ctx context.Context, userID uuid.UUID) error
	GetUserGamificationProfile(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, []domain.UserBadge, error)
}
