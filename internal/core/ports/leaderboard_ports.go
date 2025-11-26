package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type LeaderboardRepository interface {
	GetGlobalLeaderboard(ctx context.Context, limit int) ([]domain.LeaderboardEntry, error)
	GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.LeaderboardEntry, error)
}

type LeaderboardService interface {
	GetGlobalLeaderboard(ctx context.Context) ([]domain.LeaderboardEntry, error)
	GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID) ([]domain.LeaderboardEntry, error)
}
