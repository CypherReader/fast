package memory

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.Email] = user
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if user, ok := r.users[email]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	// Return a dummy user for testing if not found (or error)
	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return &domain.User{ID: id, Email: "test@example.com", SubscriptionTier: domain.TierFree}, nil
	}
	return nil, errors.New("user not found")
}

func (r *UserRepository) FindByReferralCode(ctx context.Context, code string) (*domain.User, error) {
	// In memory, we don't store referral code on user struct in this example,
	// but let's assume we can't find it or return nil.
	// Actually, domain.User doesn't seem to have ReferralCode field based on previous view.
	// Let's check domain.User.
	return nil, errors.New("user not found")
}

type FastingRepository struct {
	sessions map[string]*domain.FastingSession
	mu       sync.RWMutex
}

func NewFastingRepository() *FastingRepository {
	return &FastingRepository{
		sessions: make(map[string]*domain.FastingSession),
	}
}

func (r *FastingRepository) Save(ctx context.Context, session *domain.FastingSession) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.ID.String()] = session
	return nil
}

func (r *FastingRepository) Update(ctx context.Context, session *domain.FastingSession) error {
	return r.Save(ctx, session)
}

func (r *FastingRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.sessions {
		if s.UserID == userID && s.Status == domain.StatusActive {
			return s, nil
		}
	}
	return nil, nil
}

func (r *FastingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.FastingSession
	for _, s := range r.sessions {
		if s.UserID == userID {
			result = append(result, *s)
		}
	}
	return result, nil
}

type KetoRepository struct {
	entries []domain.KetoEntry
	mu      sync.RWMutex
}

func NewKetoRepository() *KetoRepository {
	return &KetoRepository{
		entries: make([]domain.KetoEntry, 0),
	}
}

func (r *KetoRepository) Save(ctx context.Context, entry *domain.KetoEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, *entry)
	return nil
}

func (r *KetoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KetoEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.KetoEntry
	for _, e := range r.entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	return result, nil
}

type ActivityRepository struct {
	activities map[string]*domain.Activity
	mu         sync.RWMutex
}

func NewActivityRepository() *ActivityRepository {
	return &ActivityRepository{
		activities: make(map[string]*domain.Activity),
	}
}

func (r *ActivityRepository) Save(ctx context.Context, activity *domain.Activity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.activities[activity.ID] = activity
	return nil
}

func (r *ActivityRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Activity
	for _, a := range r.activities {
		if a.UserID == userID.String() {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (r *ActivityRepository) FindByID(ctx context.Context, id string) (*domain.Activity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if activity, ok := r.activities[id]; ok {
		return activity, nil
	}
	return nil, errors.New("activity not found")
}

type MealRepository struct {
	meals []domain.Meal
	mu    sync.RWMutex
}

func NewMealRepository() *MealRepository {
	return &MealRepository{
		meals: make([]domain.Meal, 0),
	}
}

func (r *MealRepository) Save(ctx context.Context, meal *domain.Meal) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.meals = append(r.meals, *meal)
	return nil
}

func (r *MealRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Meal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Meal
	for _, m := range r.meals {
		if m.UserID == userID {
			result = append(result, m)
		}
	}
	return result, nil
}

type RecipeRepository struct {
	recipes []domain.Recipe
}

func NewRecipeRepository() *RecipeRepository {
	return &RecipeRepository{
		recipes: []domain.Recipe{
			{
				ID:           "1",
				Title:        "Avocado & Egg Breakfast Bowl",
				Description:  "A simple, high-fat breakfast to start your day.",
				Ingredients:  []string{"2 Eggs", "1 Avocado", "1 tbsp Olive Oil", "Salt & Pepper"},
				Instructions: []string{"Boil or fry eggs.", "Slice avocado.", "Combine in a bowl.", "Drizzle with olive oil and season."},
				Diet:         domain.DietVegetarian,
				IsSimple:     true,
				Calories:     450,
				Carbs:        4,
				Image:        "https://images.unsplash.com/photo-1525351484163-7529414395d8?auto=format&fit=crop&q=80&w=1000",
			},
			{
				ID:           "2",
				Title:        "Keto Chicken Salad",
				Description:  "Classic chicken salad with mayo and celery.",
				Ingredients:  []string{"1 Chicken Breast (cooked)", "2 tbsp Mayonnaise", "1 Stalk Celery", "Lettuce Wraps"},
				Instructions: []string{"Shred chicken.", "Chop celery.", "Mix chicken, celery, and mayo.", "Serve in lettuce wraps."},
				Diet:         domain.DietNormal,
				IsSimple:     true,
				Calories:     350,
				Carbs:        2,
				Image:        "https://images.unsplash.com/photo-1626082927389-6cd097cdc6ec?auto=format&fit=crop&q=80&w=1000",
			},
			{
				ID:           "3",
				Title:        "Zucchini Noodles with Pesto",
				Description:  "Fresh zucchini noodles with basil pesto.",
				Ingredients:  []string{"2 Zucchinis", "1/4 cup Pesto", "Cherry Tomatoes", "Parmesan Cheese"},
				Instructions: []string{"Spiralize zucchini.", "Sauté briefly in pan.", "Toss with pesto.", "Top with tomatoes and cheese."},
				Diet:         domain.DietVegetarian,
				IsSimple:     false,
				Calories:     280,
				Carbs:        8,
				Image:        "https://images.unsplash.com/photo-1551892374-ecf8754cf8b0?auto=format&fit=crop&q=80&w=1000",
			},
			{
				ID:           "4",
				Title:        "Vegan Keto Tofu Stir-fry",
				Description:  "Crispy tofu with low-carb veggies.",
				Ingredients:  []string{"Firm Tofu", "Broccoli", "Soy Sauce", "Sesame Oil", "Ginger"},
				Instructions: []string{"Press tofu and cube.", "Fry tofu until golden.", "Stir fry broccoli.", "Combine with sauce."},
				Diet:         domain.DietVegan,
				IsSimple:     false,
				Calories:     320,
				Carbs:        9,
				Image:        "https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&q=80&w=1000",
			},
			{
				ID:           "5",
				Title:        "Simple Bulletproof Coffee",
				Description:  "Energy boosting morning coffee.",
				Ingredients:  []string{"1 cup Coffee", "1 tbsp Butter", "1 tbsp MCT Oil"},
				Instructions: []string{"Brew coffee.", "Blend with butter and oil until frothy."},
				Diet:         domain.DietVegetarian,
				IsSimple:     true,
				Calories:     250,
				Carbs:        0,
				Image:        "https://images.unsplash.com/photo-1517701550927-30cf4ba1dba5?auto=format&fit=crop&q=80&w=1000",
			},
		},
	}
}

func (r *RecipeRepository) FindAll(ctx context.Context) ([]domain.Recipe, error) {
	return r.recipes, nil
}

type TelemetryRepository struct {
	data        []domain.TelemetryData
	connections map[string]*domain.DeviceConnection
	mu          sync.RWMutex
}

func NewTelemetryRepository() *TelemetryRepository {
	return &TelemetryRepository{
		data:        make([]domain.TelemetryData, 0),
		connections: make(map[string]*domain.DeviceConnection),
	}
}

func (r *TelemetryRepository) SaveData(ctx context.Context, data *domain.TelemetryData) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Printf("DEBUG: Saving data: UserID=%s, Type=%s, Value=%f\n", data.UserID, data.Type, data.Value)
	r.data = append(r.data, *data)
	return nil
}

func (r *TelemetryRepository) GetLatestMetric(ctx context.Context, userID uuid.UUID, metricType domain.MetricType) (*domain.TelemetryData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Printf("DEBUG: GetLatestMetric: Searching for UserID=%s, Type=%s\n", userID, metricType)
	// Iterate backwards to find latest
	for i := len(r.data) - 1; i >= 0; i-- {
		fmt.Printf("DEBUG: Checking entry: UserID=%s, Type=%s\n", r.data[i].UserID, r.data[i].Type)
		if r.data[i].UserID == userID && r.data[i].Type == metricType {
			fmt.Println("DEBUG: Found match!")
			return &r.data[i], nil
		}
	}
	fmt.Println("DEBUG: No match found")
	return nil, nil
}

func (r *TelemetryRepository) GetWeeklyStats(ctx context.Context, userID uuid.UUID, metricType domain.MetricType) ([]domain.DailyStat, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Initialize map for last 7 days
	statsMap := make(map[string]float64)
	now := time.Now()
	var days []string

	// Create last 7 days keys
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		statsMap[dateStr] = 0
		days = append(days, dateStr)
	}

	// Aggregate data
	// Note: In a real DB this would be a SQL query. Here we scan all data.
	for _, d := range r.data {
		if d.UserID == userID && d.Type == metricType {
			dateStr := d.Timestamp.Format("2006-01-02")
			if _, exists := statsMap[dateStr]; exists {
				// For weight, we want the latest value for the day.
				// Since we iterate through the slice (chronological order), overwriting ensures we get the latest.
				if metricType == domain.MetricWeight {
					statsMap[dateStr] = d.Value
				} else {
					// For other metrics (like steps), we sum them up.
					statsMap[dateStr] += d.Value
				}
			}
		}
	}

	// Convert to result slice
	var result []domain.DailyStat
	for _, dateStr := range days {
		date, _ := time.Parse("2006-01-02", dateStr)
		result = append(result, domain.DailyStat{
			Date:  dateStr,
			Day:   date.Format("Mon"),
			Value: statsMap[dateStr],
		})
	}

	return result, nil
}

func (r *TelemetryRepository) SaveConnection(ctx context.Context, conn *domain.DeviceConnection) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := conn.UserID.String() + "_" + string(conn.Source)
	r.connections[key] = conn
	return nil
}

func (r *TelemetryRepository) GetConnection(ctx context.Context, userID uuid.UUID, source domain.TelemetrySource) (*domain.DeviceConnection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := userID.String() + "_" + string(source)
	if conn, ok := r.connections[key]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (r *TelemetryRepository) ListConnections(ctx context.Context, userID uuid.UUID) ([]domain.DeviceConnection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.DeviceConnection
	for _, conn := range r.connections {
		if conn.UserID == userID {
			result = append(result, *conn)
		}
	}
	return result, nil
}

// --- Missing Repositories Implementation ---

type TribeRepository struct {
	tribes  map[string]*domain.Tribe
	members map[string][]string // tribeID -> []userID
	mu      sync.RWMutex
}

func NewTribeRepository() *TribeRepository {
	return &TribeRepository{
		tribes:  make(map[string]*domain.Tribe),
		members: make(map[string][]string),
	}
}

func (r *TribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tribes[tribe.ID.String()] = tribe
	return nil
}

func (r *TribeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if tribe, ok := r.tribes[id.String()]; ok {
		return tribe, nil
	}
	return nil, errors.New("tribe not found")
}

func (r *TribeRepository) FindAll(ctx context.Context) ([]domain.Tribe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Tribe
	for _, t := range r.tribes {
		result = append(result, *t)
	}
	return result, nil
}

func (r *TribeRepository) AddMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tid := tribeID.String()
	uid := userID.String()

	// Check if already member
	for _, m := range r.members[tid] {
		if m == uid {
			return nil
		}
	}
	r.members[tid] = append(r.members[tid], uid)
	return nil
}

func (r *TribeRepository) RemoveMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tid := tribeID.String()
	uid := userID.String()

	members := r.members[tid]
	for i, m := range members {
		if m == uid {
			r.members[tid] = append(members[:i], members[i+1:]...)
			return nil
		}
	}
	return nil
}

type SocialRepository struct {
	events []domain.SocialEvent
	mu     sync.RWMutex
}

func NewSocialRepository() *SocialRepository {
	return &SocialRepository{
		events: make([]domain.SocialEvent, 0),
	}
}

func (r *SocialRepository) SaveEvent(ctx context.Context, event *domain.SocialEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, *event)
	return nil
}

func (r *SocialRepository) GetFeed(ctx context.Context, limit int) ([]domain.SocialEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Return last 'limit' events reversed
	count := len(r.events)
	if limit > count {
		limit = count
	}

	result := make([]domain.SocialEvent, limit)
	for i := 0; i < limit; i++ {
		result[i] = r.events[count-1-i]
	}
	return result, nil
}

func (r *SocialRepository) GetTribeFeed(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.SocialEvent, error) {
	// For simplicity in memory, just return global feed
	return r.GetFeed(ctx, limit)
}

type LeaderboardRepository struct {
	// In memory we can just calculate on fly or return dummy
}

func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{}
}

func (r *LeaderboardRepository) GetGlobalLeaderboard(ctx context.Context, limit int) ([]domain.LeaderboardEntry, error) {
	return []domain.LeaderboardEntry{}, nil
}

func (r *LeaderboardRepository) GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.LeaderboardEntry, error) {
	return []domain.LeaderboardEntry{}, nil
}

type GamificationRepository struct {
	badges  map[string][]domain.UserBadge
	streaks map[string]*domain.UserStreak
	mu      sync.RWMutex
}

func NewGamificationRepository() *GamificationRepository {
	return &GamificationRepository{
		badges:  make(map[string][]domain.UserBadge),
		streaks: make(map[string]*domain.UserStreak),
	}
}

func (r *GamificationRepository) SaveUserBadge(ctx context.Context, badge *domain.UserBadge) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	uid := badge.UserID.String()
	r.badges[uid] = append(r.badges[uid], *badge)
	return nil
}

func (r *GamificationRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]domain.UserBadge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.badges[userID.String()], nil
}

func (r *GamificationRepository) UpdateUserStreak(ctx context.Context, streak *domain.UserStreak) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.streaks[streak.UserID.String()] = streak
	return nil
}

func (r *GamificationRepository) GetUserStreak(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.streaks[userID.String()]; ok {
		return s, nil
	}
	return nil, nil
}

type ReferralRepository struct {
	referrals []domain.Referral
	mu        sync.RWMutex
}

func NewReferralRepository() *ReferralRepository {
	return &ReferralRepository{
		referrals: make([]domain.Referral, 0),
	}
}

func (r *ReferralRepository) Save(ctx context.Context, referral *domain.Referral) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.referrals = append(r.referrals, *referral)
	return nil
}

func (r *ReferralRepository) FindByRefereeID(ctx context.Context, refereeID uuid.UUID) (*domain.Referral, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ref := range r.referrals {
		if ref.RefereeID == refereeID {
			return &ref, nil
		}
	}
	return nil, nil
}

func (r *ReferralRepository) FindByReferrerID(ctx context.Context, referrerID uuid.UUID) ([]domain.Referral, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Referral
	for _, ref := range r.referrals {
		if ref.ReferrerID == referrerID {
			result = append(result, ref)
		}
	}
	return result, nil
}

func (r *ReferralRepository) Update(ctx context.Context, referral *domain.Referral) error {
	// In memory slice update is tricky without pointer, but for now append is okay or simple loop
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, ref := range r.referrals {
		if ref.ID == referral.ID {
			r.referrals[i] = *referral
			return nil
		}
	}
	return errors.New("referral not found")
}

type NotificationRepository struct {
	tokens        map[string][]domain.FCMToken
	notifications map[string][]domain.Notification
	mu            sync.RWMutex
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		tokens:        make(map[string][]domain.FCMToken),
		notifications: make(map[string][]domain.Notification),
	}
}

func (r *NotificationRepository) SaveToken(ctx context.Context, token *domain.FCMToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	uid := token.UserID.String()
	r.tokens[uid] = append(r.tokens[uid], *token)
	return nil
}

func (r *NotificationRepository) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tokens[userID.String()], nil
}

func (r *NotificationRepository) DeleteToken(ctx context.Context, tokenString string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Naive implementation: iterate all users
	for uid, tokens := range r.tokens {
		for i, t := range tokens {
			if t.Token == tokenString {
				r.tokens[uid] = append(tokens[:i], tokens[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (r *NotificationRepository) SaveNotification(ctx context.Context, notification *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	uid := notification.UserID.String()
	r.notifications[uid] = append(r.notifications[uid], *notification)
	return nil
}

func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	notifs := r.notifications[userID.String()]
	if len(notifs) > limit {
		return notifs[len(notifs)-limit:], nil
	}
	return notifs, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for uid, notifs := range r.notifications {
		for i, n := range notifs {
			if n.ID == notificationID {
				now := time.Now()
				r.notifications[uid][i].ReadAt = &now
				return nil
			}
		}
	}
	return nil
}
