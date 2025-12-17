package memory

import (
	"context"
	"fastinghero/internal/core/domain"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemorySOSRepository struct {
	mu            sync.RWMutex
	sosFlares     map[uuid.UUID]*domain.SOSFlare
	hypeResponses map[uuid.UUID][]domain.HypeResponse // sosID -> []hypes
}

func NewMemorySOSRepository() *MemorySOSRepository {
	return &MemorySOSRepository{
		sosFlares:     make(map[uuid.UUID]*domain.SOSFlare),
		hypeResponses: make(map[uuid.UUID][]domain.HypeResponse),
	}
}

func (r *MemorySOSRepository) Save(ctx context.Context, sos *domain.SOSFlare) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sosFlares[sos.ID] = sos
	return nil
}

func (r *MemorySOSRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SOSFlare, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sos, exists := r.sosFlares[id]
	if !exists {
		return nil, nil
	}

	return sos, nil
}

func (r *MemorySOSRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.SOSFlare, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, sos := range r.sosFlares {
		if sos.UserID == userID && sos.Status == domain.SOSStatusActive {
			return sos, nil
		}
	}

	return nil, nil
}

func (r *MemorySOSRepository) UpdateStatus(ctx context.Context, sosID uuid.UUID, status domain.SOSStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sos, exists := r.sosFlares[sosID]
	if !exists {
		return nil
	}

	sos.Status = status
	now := time.Now()
	sos.ResolvedAt = &now
	return nil
}

func (r *MemorySOSRepository) UpdateCortexResponse(ctx context.Context, sosID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sos, exists := r.sosFlares[sosID]
	if !exists {
		return nil
	}

	sos.CortexResponded = true
	return nil
}

func (r *MemorySOSRepository) SaveHypeResponse(ctx context.Context, hype *domain.HypeResponse) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hypeResponses[hype.SOSID] = append(r.hypeResponses[hype.SOSID], *hype)
	return nil
}

func (r *MemorySOSRepository) GetHypeResponses(ctx context.Context, sosID uuid.UUID) ([]domain.HypeResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hypes, exists := r.hypeResponses[sosID]
	if !exists {
		return []domain.HypeResponse{}, nil
	}

	return hypes, nil
}

func (r *MemorySOSRepository) GetUserHypeCount(ctx context.Context, userID uuid.UUID, since time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, hypes := range r.hypeResponses {
		for _, hype := range hypes {
			if hype.FromUserID == userID && hype.CreatedAt.After(since) {
				count++
			}
		}
	}

	return count, nil
}

func (r *MemorySOSRepository) FindAllActive(ctx context.Context) ([]*domain.SOSFlare, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeFlares []*domain.SOSFlare
	for _, sos := range r.sosFlares {
		if sos.Status == domain.SOSStatusActive {
			activeFlares = append(activeFlares, sos)
		}
	}

	return activeFlares, nil
}

func (r *MemorySOSRepository) IncrementHypeCount(ctx context.Context, sosID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sos, exists := r.sosFlares[sosID]
	if !exists {
		return nil
	}

	sos.HypeCount++
	return nil
}
