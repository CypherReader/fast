package memory

import (
	"context"
	"fastinghero/internal/core/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ReminderRepository is an in-memory implementation of the reminder repository
type ReminderRepository struct {
	reminders map[uuid.UUID]*domain.ScheduledReminder
	settings  map[uuid.UUID]*domain.ReminderSettings
	mu        sync.RWMutex
}

// NewReminderRepository creates a new in-memory reminder repository
func NewReminderRepository() *ReminderRepository {
	return &ReminderRepository{
		reminders: make(map[uuid.UUID]*domain.ScheduledReminder),
		settings:  make(map[uuid.UUID]*domain.ReminderSettings),
	}
}

// Save saves a scheduled reminder
func (r *ReminderRepository) Save(ctx context.Context, reminder *domain.ScheduledReminder) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reminders[reminder.ID] = reminder
	return nil
}

// FindPending returns all pending reminders scheduled before the given time
func (r *ReminderRepository) FindPending(ctx context.Context, before time.Time) ([]domain.ScheduledReminder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.ScheduledReminder
	for _, reminder := range r.reminders {
		if !reminder.Sent && reminder.ScheduledAt.Before(before) {
			result = append(result, *reminder)
		}
	}
	return result, nil
}

// MarkSent marks a reminder as sent
func (r *ReminderRepository) MarkSent(ctx context.Context, reminderID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reminder, ok := r.reminders[reminderID]; ok {
		reminder.Sent = true
	}
	return nil
}

// DeleteByUserAndType deletes all reminders of a specific type for a user
func (r *ReminderRepository) DeleteByUserAndType(ctx context.Context, userID uuid.UUID, reminderType domain.ReminderType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	toDelete := []uuid.UUID{}
	for id, reminder := range r.reminders {
		if reminder.UserID == userID && reminder.ReminderType == reminderType && !reminder.Sent {
			toDelete = append(toDelete, id)
		}
	}
	for _, id := range toDelete {
		delete(r.reminders, id)
	}
	return nil
}

// GetUserSettings retrieves reminder settings for a user
func (r *ReminderRepository) GetUserSettings(ctx context.Context, userID uuid.UUID) (*domain.ReminderSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if settings, ok := r.settings[userID]; ok {
		return settings, nil
	}
	// Return default settings
	return &domain.ReminderSettings{
		UserID:                   userID,
		ReminderFastStart:        true,
		ReminderFastEnd:          true,
		ReminderHydration:        false,
		PreferredFastStartHour:   20,
		HydrationIntervalMinutes: 60,
	}, nil
}

// SaveUserSettings saves reminder settings for a user
func (r *ReminderRepository) SaveUserSettings(ctx context.Context, settings *domain.ReminderSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings[settings.UserID] = settings
	return nil
}
