package ports

import (
	"context"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

// Primary Ports (Services)

type AuthService interface {
	Register(ctx context.Context, email, password, name, referralCode string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, string, *domain.User, error) // token, refresh, user, error
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

type OnboardingService interface {
	UpdateProfile(ctx context.Context, userID uuid.UUID, profile domain.UserProfileUpdate) (*domain.User, error)
	CompleteOnboarding(ctx context.Context, userID uuid.UUID) (*domain.User, error)
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
	GetCurrentParticipation(ctx context.Context, userID uuid.UUID) (*domain.VaultParticipation, error)
}

type PaymentService interface {
	CreateCustomer(ctx context.Context, user *domain.User) (string, error)
	CreateSubscription(ctx context.Context, userID uuid.UUID, priceID string) (*domain.Subscription, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
}

// Secondary Ports (Repositories)

type CortexService interface {
	Chat(ctx context.Context, userID uuid.UUID, message string) (string, error)
	GenerateInsight(ctx context.Context, userID uuid.UUID, fastingHours float64) (string, error)
	AnalyzeMeal(ctx context.Context, imageBase64, description string) (string, bool, bool, error)
	GetCravingHelp(ctx context.Context, userID uuid.UUID, cravingDescription string) (interface{}, error)
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
	GetHistory(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error)
}

type NotificationRepository interface {
	SaveToken(ctx context.Context, token *domain.FCMToken) error
	GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error)
	DeleteToken(ctx context.Context, tokenString string) error
	Save(ctx context.Context, notification *domain.Notification) error
	FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

type SubscriptionRepository interface {
	Save(ctx context.Context, sub *domain.Subscription) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error)
	FindByStripeSubscriptionID(ctx context.Context, stripeSubID string) (*domain.Subscription, error)
}

type VaultRepository interface {
	Save(ctx context.Context, vault *domain.VaultParticipation) error
	FindByUserIDAndMonth(ctx context.Context, userID uuid.UUID, monthStart time.Time) (*domain.VaultParticipation, error)
}
type SocialService interface {
	AddFriend(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendNetwork, error)
	CreateTribe(ctx context.Context, userID uuid.UUID, name, description string, isPublic bool) (*domain.Tribe, error)
	GetTribe(ctx context.Context, tribeID uuid.UUID) (*domain.Tribe, error)
	CreateChallenge(ctx context.Context, userID uuid.UUID, name string, challengeType domain.ChallengeType, goal int, startDate, endDate time.Time) (*domain.FriendChallenge, error)
	GetChallenges(ctx context.Context, userID uuid.UUID) ([]domain.FriendChallenge, error)
	ListTribes(ctx context.Context, limit, offset int) ([]domain.Tribe, error)
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.SocialEvent, error)
}

type SocialRepository interface {
	SaveFriendNetwork(ctx context.Context, fn *domain.FriendNetwork) error
	FindFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendNetwork, error)
	SaveTribe(ctx context.Context, tribe *domain.Tribe) error
	FindTribeByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error)
	SaveChallenge(ctx context.Context, c *domain.FriendChallenge) error
	FindChallengesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FriendChallenge, error)
	FindAllTribes(ctx context.Context, limit, offset int) ([]domain.Tribe, error)
	GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.SocialEvent, error)
	SaveEvent(ctx context.Context, event *domain.SocialEvent) error
}

type ProgressService interface {
	LogWeight(ctx context.Context, userID uuid.UUID, weight float64, unit string) (*domain.WeightLog, error)
	GetWeightHistory(ctx context.Context, userID uuid.UUID, days int) ([]domain.WeightLog, error)
	LogHydration(ctx context.Context, userID uuid.UUID, amount float64, unit string) (*domain.HydrationLog, error)
	GetDailyHydration(ctx context.Context, userID uuid.UUID) (*domain.HydrationLog, error)
}

type ProgressRepository interface {
	SaveWeightLog(ctx context.Context, log *domain.WeightLog) error
	GetWeightHistory(ctx context.Context, userID uuid.UUID, days int) ([]domain.WeightLog, error)
	SaveHydrationLog(ctx context.Context, log *domain.HydrationLog) error
	GetHydrationLog(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.HydrationLog, error)
	SaveActivityLog(ctx context.Context, log *domain.ActivityLog) error
	GetActivityStats(ctx context.Context, userID uuid.UUID, days int) ([]domain.ActivityLog, error)
}

// TribeRepository defines the interface for tribe data persistence
type TribeRepository interface {
	// Tribe CRUD
	Save(ctx context.Context, tribe *domain.Tribe) error
	Update(ctx context.Context, tribe *domain.Tribe) error
	FindByID(ctx context.Context, id string) (*domain.Tribe, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Tribe, error)
	List(ctx context.Context, query domain.ListTribesQuery) ([]domain.Tribe, int, error) // tribes, total count, error
	Delete(ctx context.Context, id string) error

	// Memberships
	SaveMembership(ctx context.Context, membership *domain.TribeMembership) error
	UpdateMembership(ctx context.Context, membership *domain.TribeMembership) error
	FindMembership(ctx context.Context, tribeID, userID string) (*domain.TribeMembership, error)
	GetMembersByTribeID(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error)
	GetUserTribes(ctx context.Context, userID string, status string) ([]domain.Tribe, error)
	GetMembershipCount(ctx context.Context, tribeID string) (int, error)
	DeleteMembership(ctx context.Context, tribeID, userID string) error

	// Stats
	GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error)
	UpdateMemberCounts(ctx context.Context, tribeID string) error
}

// TribeService defines the business logic for tribe operations
type TribeService interface {
	// Tribe management
	CreateTribe(ctx context.Context, userID string, req domain.CreateTribeRequest) (*domain.Tribe, error)
	GetTribe(ctx context.Context, tribeID string, currentUserID *string) (*domain.Tribe, error)
	UpdateTribe(ctx context.Context, tribeID, userID string, req domain.UpdateTribeRequest) (*domain.Tribe, error)
	DeleteTribe(ctx context.Context, tribeID, userID string) error
	ListTribes(ctx context.Context, query domain.ListTribesQuery, currentUserID *string) ([]domain.Tribe, int, error)

	// Membership operations
	JoinTribe(ctx context.Context, tribeID, userID string) error
	LeaveTribe(ctx context.Context, tribeID, userID string) error
	GetTribeMembers(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error)
	GetMyTribes(ctx context.Context, userID string) ([]domain.Tribe, error)

	// Stats and analytics
	GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error)
}
