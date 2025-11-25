package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type TribeRepository interface {
	Create(ctx context.Context, tribe *domain.Tribe) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error)
	Update(ctx context.Context, tribe *domain.Tribe) error
}

type TribeService struct {
	repo TribeRepository
}

func NewTribeService(repo TribeRepository) *TribeService {
	return &TribeService{repo: repo}
}

func (s *TribeService) CreateTribe(ctx context.Context, name string, creatorID uuid.UUID) (*domain.Tribe, error) {
	tribe := &domain.Tribe{
		ID:              uuid.New(),
		Name:            name,
		MemberIDs:       []uuid.UUID{creatorID},
		CollectiveScore: 100.0, // Start perfect
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, tribe); err != nil {
		return nil, err
	}

	return tribe, nil
}

func (s *TribeService) JoinTribe(ctx context.Context, tribeID uuid.UUID, userID uuid.UUID) error {
	tribe, err := s.repo.GetByID(ctx, tribeID)
	if err != nil {
		return err
	}
	if tribe == nil {
		return errors.New("tribe not found")
	}

	// Check if already a member
	for _, id := range tribe.MemberIDs {
		if id == userID {
			return errors.New("already a member")
		}
	}

	tribe.MemberIDs = append(tribe.MemberIDs, userID)
	tribe.UpdatedAt = time.Now()

	return s.repo.Update(ctx, tribe)
}
