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
				// For steps, we might want max or sum depending on how we log.
				// Assuming we log cumulative or discrete chunks, let's sum for now.
				// If we logged "total steps at time X", we'd want MAX.
				// Let's assume discrete chunks for manual entry (e.g. +5000 steps).
				statsMap[dateStr] += d.Value
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
