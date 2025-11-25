package ports

import (
	"context"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

// Primary Ports (Services)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // token, refresh, error
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

type FastingService interface {
	StartFast(ctx context.Context, userID uuid.UUID, plan domain.FastingPlanType, goalHours int) (*domain.FastingSession, error)
	StopFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error)
	GetCurrentFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error)
	GetFastingHistory(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error)
}

type KetoService interface {
	LogEntry(ctx context.Context, userID uuid.UUID, entry domain.KetoEntry) error
}

// Secondary Ports (Repositories)

type CortexService interface {
	Chat(ctx context.Context, userID uuid.UUID, message string) (string, error)
	GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error)
}

// Secondary Ports (Repositories & Adapters)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type FastingRepository interface {
	Save(ctx context.Context, session *domain.FastingSession) error
	Update(ctx context.Context, session *domain.FastingSession) error
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error)
}

type KetoRepository interface {
	Save(ctx context.Context, entry *domain.KetoEntry) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KetoEntry, error)
}

type LLMProvider interface {
	GenerateResponse(ctx context.Context, prompt string, systemPrompt string) (string, error)
}
