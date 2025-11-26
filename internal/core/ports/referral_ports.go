package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type ReferralRepository interface {
	Save(ctx context.Context, referral *domain.Referral) error
	FindByRefereeID(ctx context.Context, refereeID uuid.UUID) (*domain.Referral, error)
	FindByReferrerID(ctx context.Context, referrerID uuid.UUID) ([]domain.Referral, error)
	Update(ctx context.Context, referral *domain.Referral) error
}

type ReferralService interface {
	GenerateReferralCode(ctx context.Context, userID uuid.UUID) (string, error)
	GetReferralCode(ctx context.Context, userID uuid.UUID) (string, error)
	TrackReferral(ctx context.Context, referrerCode string, refereeID uuid.UUID) error
	CompleteReferral(ctx context.Context, refereeID uuid.UUID) error
	GetReferralStats(ctx context.Context, userID uuid.UUID) (totalEarned float64, count int, err error)
}
