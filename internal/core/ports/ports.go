package ports

import (
	"context"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// Primary Ports (Services)

type AuthService interface {
	Register(ctx context.Context, email, password string, referralCode string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, error) // token, refresh, error
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

type FastingService interface {
	StartFast(ctx context.Context, userID uuid.UUID, plan domain.FastingPlanType, goalHours int, startTime *time.Time) (*domain.FastingSession, error)
	StopFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error)
	GetCurrentFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error)
	GetFastingHistory(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error)
}

type KetoService interface {
	LogEntry(ctx context.Context, userID uuid.UUID, entry domain.KetoEntry) error
}

type VaultService interface {
	CalculateVaultStatus(user *domain.User) (deposit float64, earned float64, potentialRefund float64)
	CalculateDailyEarning(disciplineIndex int) float64
	ProcessDailyEarnings(ctx context.Context) error
	AddDailyEarnings(ctx context.Context, user *domain.User, amount float64)
	CalculatePrice(ctx context.Context, user *domain.User) float64
	UpdateDisciplineIndex(ctx context.Context, user *domain.User, completedFast bool, verifiedKetosis bool)
}

// Secondary Ports (Repositories)

type CortexService interface {
	Chat(ctx context.Context, userID uuid.UUID, message string) (string, error)
	GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error)
	AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error)
}

// Secondary Ports (Repositories & Adapters)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByReferralCode(ctx context.Context, code string) (*domain.User, error)
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
	AnalyzeImage(ctx context.Context, imageBase64, prompt string) (string, error)
}

type ActivityService interface {
	SyncActivity(ctx context.Context, userID uuid.UUID, activity domain.Activity) error
	GetActivities(ctx context.Context, userID uuid.UUID) ([]domain.Activity, error)
	GetActivity(ctx context.Context, activityID string) (*domain.Activity, error)
}

type ActivityRepository interface {
	Save(ctx context.Context, activity *domain.Activity) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Activity, error)
	FindByID(ctx context.Context, id string) (*domain.Activity, error)
}

type MealRepository interface {
	Save(ctx context.Context, meal *domain.Meal) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Meal, error)
}

type MealService interface {
	LogMeal(ctx context.Context, userID uuid.UUID, name string, calories int, mealType string, image, description string) (*domain.Meal, error)
	GetMeals(ctx context.Context, userID uuid.UUID) ([]domain.Meal, error)
}

type RecipeRepository interface {
	FindAll(ctx context.Context) ([]domain.Recipe, error)
}

type RecipeService interface {
	GetRecipes(ctx context.Context, diet domain.DietType) ([]domain.Recipe, error)
}

type NotificationService interface {
	SendNotification(ctx context.Context, userID uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error
	SendBatchNotification(ctx context.Context, userIDs []uuid.UUID, title, body string, notifType domain.NotificationType, data map[string]string) error
	RegisterFCMToken(ctx context.Context, userID uuid.UUID, token, deviceType string) error
	UnregisterFCMToken(ctx context.Context, userID uuid.UUID, token string) error
}

type NotificationRepository interface {
	SaveToken(ctx context.Context, token *domain.FCMToken) error
	GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error)
	DeleteToken(ctx context.Context, tokenString string) error
	SaveNotification(ctx context.Context, notification *domain.Notification) error
	GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
}
