package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type KetoService struct {
	repo     ports.KetoRepository
	userRepo ports.UserRepository
}

func NewKetoService(repo ports.KetoRepository, userRepo ports.UserRepository) *KetoService {
	return &KetoService{repo: repo, userRepo: userRepo}
}

func (s *KetoService) LogEntry(ctx context.Context, userID uuid.UUID, entry domain.KetoEntry) error {
	// Check if user is premium if hard data is present
	if entry.KetoneLevel != nil || entry.AcetoneLevel != nil {
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			return err
		}
		if !user.IsPremium() {
			return errors.New("premium subscription required for hard data inputs")
		}
	}

	entry.UserID = userID
	return s.repo.Save(ctx, &entry)
}
