package memory

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"sort"
	"strings"
	"sync"
	"time"
)

type MemoryTribeRepository struct {
	tribes      map[string]*domain.Tribe
	memberships map[string]*domain.TribeMembership // ID -> Membership
	mu          sync.RWMutex
}

func NewTribeRepository() *MemoryTribeRepository {
	return &MemoryTribeRepository{
		tribes:      make(map[string]*domain.Tribe),
		memberships: make(map[string]*domain.TribeMembership),
	}
}

// Tribe CRUD

func (r *MemoryTribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tribes[tribe.ID] = tribe
	return nil
}

func (r *MemoryTribeRepository) Update(ctx context.Context, tribe *domain.Tribe) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tribes[tribe.ID]; !exists {
		return errors.New("tribe not found")
	}
	r.tribes[tribe.ID] = tribe
	return nil
}

func (r *MemoryTribeRepository) FindByID(ctx context.Context, id string) (*domain.Tribe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tribe, exists := r.tribes[id]
	if !exists {
		return nil, errors.New("tribe not found")
	}
	return tribe, nil
}

func (r *MemoryTribeRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tribe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tribes {
		if t.Slug == slug {
			return t, nil
		}
	}
	return nil, errors.New("tribe not found")
}

func (r *MemoryTribeRepository) List(ctx context.Context, query domain.ListTribesQuery) ([]domain.Tribe, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.Tribe
	for _, t := range r.tribes {
		if t.DeletedAt != nil {
			continue
		}
		// Filter
		if query.Search != "" {
			if !strings.Contains(strings.ToLower(t.Name), strings.ToLower(query.Search)) {
				continue
			}
		}
		if query.FastingSchedule != "" && t.FastingSchedule != query.FastingSchedule {
			continue
		}
		if query.PrimaryGoal != "" && t.PrimaryGoal != query.PrimaryGoal {
			continue
		}
		if query.Privacy != "" && t.Privacy != query.Privacy {
			continue
		}
		result = append(result, *t)
	}

	total := len(result)

	// Sort
	sort.Slice(result, func(i, j int) bool {
		switch query.SortBy {
		case "members":
			return result[i].MemberCount > result[j].MemberCount
		case "newest":
			return result[i].CreatedAt.After(result[j].CreatedAt)
		default:
			return result[i].CreatedAt.After(result[j].CreatedAt)
		}
	})

	// Pagination
	start := query.Offset
	if start > total {
		start = total
	}
	end := start + query.Limit
	if end > total {
		end = total
	}

	return result[start:end], total, nil
}

func (r *MemoryTribeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t, exists := r.tribes[id]; exists {
		now := time.Now()
		t.DeletedAt = &now
	}
	return nil
}

// Memberships

func (r *MemoryTribeRepository) SaveMembership(ctx context.Context, membership *domain.TribeMembership) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.memberships[membership.ID] = membership
	return nil
}

func (r *MemoryTribeRepository) UpdateMembership(ctx context.Context, membership *domain.TribeMembership) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.memberships[membership.ID]; !exists {
		return errors.New("membership not found")
	}
	r.memberships[membership.ID] = membership
	return nil
}

func (r *MemoryTribeRepository) FindMembership(ctx context.Context, tribeID, userID string) (*domain.TribeMembership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, m := range r.memberships {
		if m.TribeID == tribeID && m.UserID == userID {
			return m, nil
		}
	}
	return nil, errors.New("membership not found")
}

func (r *MemoryTribeRepository) GetMembersByTribeID(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var members []domain.TribeMember
	for _, m := range r.memberships {
		if m.TribeID == tribeID && m.Status == "active" {
			members = append(members, domain.TribeMember{
				TribeMembership: *m,
				UserName:        "Test User", // Mock data
				UserAvatar:      "",
				UserStreak:      5,
			})
		}
	}

	// Pagination
	start := offset
	if start > len(members) {
		start = len(members)
	}
	end := start + limit
	if end > len(members) {
		end = len(members)
	}

	return members[start:end], nil
}

func (r *MemoryTribeRepository) GetUserTribes(ctx context.Context, userID string, status string) ([]domain.Tribe, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tribes []domain.Tribe
	for _, m := range r.memberships {
		if m.UserID == userID {
			if status != "" && m.Status != status {
				continue
			}
			if t, exists := r.tribes[m.TribeID]; exists {
				tribes = append(tribes, *t)
			}
		}
	}
	return tribes, nil
}

func (r *MemoryTribeRepository) GetMembershipCount(ctx context.Context, tribeID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, m := range r.memberships {
		if m.TribeID == tribeID && m.Status == "active" {
			count++
		}
	}
	return count, nil
}

func (r *MemoryTribeRepository) DeleteMembership(ctx context.Context, tribeID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, m := range r.memberships {
		if m.TribeID == tribeID && m.UserID == userID {
			delete(r.memberships, id)
			return nil
		}
	}
	return errors.New("membership not found")
}

// Stats

func (r *MemoryTribeRepository) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	return &domain.TribeStats{
		TribeID:              tribeID,
		TotalFasts:           120,
		TotalFastingHours:    2400.5,
		AverageMemberStreak:  7.5,
		WeeklyGrowthPercent:  15.0,
		ActiveMembersPercent: 85.0,
	}, nil
}

func (r *MemoryTribeRepository) UpdateMemberCounts(ctx context.Context, tribeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	count := 0
	for _, m := range r.memberships {
		if m.TribeID == tribeID && m.Status == "active" {
			count++
		}
	}

	if t, exists := r.tribes[tribeID]; exists {
		t.MemberCount = count
	}
	return nil
}
