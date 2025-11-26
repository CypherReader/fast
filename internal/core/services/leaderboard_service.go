package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type LeaderboardService struct {
	repo ports.LeaderboardRepository
}

func NewLeaderboardService(repo ports.LeaderboardRepository) *LeaderboardService {
	return &LeaderboardService{
		repo: repo,
	}
}

func (s *LeaderboardService) GetGlobalLeaderboard(ctx context.Context) ([]domain.LeaderboardEntry, error) {
	return s.repo.GetGlobalLeaderboard(ctx, 50) // Top 50
}

func (s *LeaderboardService) GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID) ([]domain.LeaderboardEntry, error) {
	return s.repo.GetTribeLeaderboard(ctx, tribeID, 50) // Top 50
}
