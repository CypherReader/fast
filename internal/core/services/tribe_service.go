package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type TribeService struct {
	repo     ports.TribeRepository
	userRepo ports.UserRepository
}

func NewTribeService(repo ports.TribeRepository, userRepo ports.UserRepository) *TribeService {
	return &TribeService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *TribeService) CreateTribe(ctx context.Context, name, description string, leaderID uuid.UUID) (*domain.Tribe, error) {
	tribe := &domain.Tribe{
		ID:              uuid.New(),
		Name:            name,
		Description:     description,
		LeaderID:        leaderID,
		MemberCount:     1, // Leader is the first member
		TotalDiscipline: 0,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.Save(ctx, tribe); err != nil {
		return nil, err
	}

	// Add leader as member
	if err := s.repo.AddMember(ctx, tribe.ID, leaderID); err != nil {
		return nil, err
	}

	return tribe, nil
}

func (s *TribeService) JoinTribe(ctx context.Context, tribeID, userID uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.TribeID != nil {
		return errors.New("user is already in a tribe")
	}

	if err := s.repo.AddMember(ctx, tribeID, userID); err != nil {
		return err
	}

	return nil
}

func (s *TribeService) LeaveTribe(ctx context.Context, tribeID, userID uuid.UUID) error {
	return s.repo.RemoveMember(ctx, tribeID, userID)
}

func (s *TribeService) GetTribeDetails(ctx context.Context, tribeID uuid.UUID) (*domain.Tribe, error) {
	return s.repo.FindByID(ctx, tribeID)
}

func (s *TribeService) ListTribes(ctx context.Context) ([]domain.Tribe, error) {
	return s.repo.FindAll(ctx)
}
