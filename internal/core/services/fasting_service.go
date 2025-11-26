package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type FastingService struct {
	repo     ports.FastingRepository
	pricing  *PricingService
	userRepo ports.UserRepository
}

func NewFastingService(repo ports.FastingRepository, pricing *PricingService, userRepo ports.UserRepository) *FastingService {
	return &FastingService{
		repo:     repo,
		pricing:  pricing,
		userRepo: userRepo,
	}
}

func (s *FastingService) StartFast(ctx context.Context, userID uuid.UUID, plan domain.FastingPlanType, goalHours int, startTime *time.Time) (*domain.FastingSession, error) {
	// Check if active fast exists
	active, _ := s.repo.FindActiveByUserID(ctx, userID)
	if active != nil {
		return nil, errors.New("active fasting session already exists")
	}

	st := time.Now()
	if startTime != nil {
		st = *startTime
	}

	session := domain.NewFastingSession(userID, plan, goalHours, st)
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *FastingService) StopFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	session, err := s.repo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("no active fasting session")
	}

	// 1. Mark as completed
	session.Status = domain.StatusCompleted
	now := time.Now()
	session.EndTime = &now

	// 2. Calculate Duration
	duration := session.EndTime.Sub(session.StartTime).Hours()
	goalMet := duration >= float64(session.GoalHours)

	// 3. Update Discipline & Price
	user, err := s.userRepo.FindByID(ctx, userID)
	if err == nil {
		// Update discipline based on goal completion
		// If goal met: +1, if not: -1 (simple logic for now)
		// We reuse UpdateDisciplineIndex but might need to adjust it to handle "missed" logic better
		// For now, let's assume UpdateDisciplineIndex handles positive reinforcement
		// We'll add a manual check here for negative reinforcement if needed,
		// but PricingService.UpdateDisciplineIndex currently only adds.
		// Let's modify logic slightly:

		if goalMet {
			s.pricing.UpdateDisciplineIndex(ctx, user, true, false)
		} else {
			// Penalize for quitting early (Lazy Tax)
			user.DisciplineIndex -= 2
			if user.DisciplineIndex < 0 {
				user.DisciplineIndex = 0
			}
			// Recalculate price
			user.CurrentPrice = s.pricing.CalculatePrice(ctx, user)
		}

		_ = s.userRepo.Save(ctx, user)
	}

	// 4. Update repo
	if err := s.repo.Update(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *FastingService) GetCurrentFast(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	return s.repo.FindActiveByUserID(ctx, userID)
}

func (s *FastingService) GetFastingHistory(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error) {
	return s.repo.FindByUserID(ctx, userID)
}
