package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type TribeRepository interface {
	Save(ctx context.Context, tribe *domain.Tribe) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error)
	FindAll(ctx context.Context) ([]domain.Tribe, error)
	AddMember(ctx context.Context, tribeID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, tribeID, userID uuid.UUID) error
}

type TribeService interface {
	CreateTribe(ctx context.Context, name, description string, leaderID uuid.UUID) (*domain.Tribe, error)
	JoinTribe(ctx context.Context, tribeID, userID uuid.UUID) error
	LeaveTribe(ctx context.Context, tribeID, userID uuid.UUID) error
	GetTribeDetails(ctx context.Context, tribeID uuid.UUID) (*domain.Tribe, error)
	ListTribes(ctx context.Context) ([]domain.Tribe, error)
}
