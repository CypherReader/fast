package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type SocialRepository interface {
	SaveEvent(ctx context.Context, event *domain.SocialEvent) error
	GetFeed(ctx context.Context, limit int) ([]domain.SocialEvent, error)
	GetTribeFeed(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.SocialEvent, error)
}

type SocialService interface {
	LogEvent(ctx context.Context, userID uuid.UUID, eventType domain.SocialEventType, data interface{}) error
	GetFeed(ctx context.Context) ([]domain.SocialEvent, error)
}
