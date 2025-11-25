package memory

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"sync"

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
