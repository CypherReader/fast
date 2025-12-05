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
	repo         ports.FastingRepository
	vaultService ports.VaultService
	userRepo     ports.UserRepository
}

func NewFastingService(repo ports.FastingRepository, vaultService ports.VaultService, userRepo ports.UserRepository) *FastingService {
	return &FastingService{
		repo:         repo,
		vaultService: vaultService,
		userRepo:     userRepo,
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

	// Link to Vault Participation if exists
	vault, err := s.vaultService.GetCurrentParticipation(ctx, userID)
	if err == nil && vault != nil {
		session.VaultParticipationID = &vault.ID
	}

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
	session.ActualDurationHours = duration
	goalMet := duration >= float64(session.GoalHours)

	// Calculate Phase Reached
	if duration >= 72 {
		session.PhaseReached = "Immune Regeneration"
	} else if duration >= 48 {
		session.PhaseReached = "Deep Autophagy"
	} else if duration >= 24 {
		session.PhaseReached = "Autophagy"
	} else if duration >= 18 {
		session.PhaseReached = "Ketosis"
	} else if duration >= 12 {
		session.PhaseReached = "Catabolic"
	} else {
		session.PhaseReached = "Anabolic"
	}

	// 3. Update Discipline & Price
	user, err := s.userRepo.FindByID(ctx, userID)
	if err == nil {
		// Update discipline based on goal completion
		if goalMet {
			s.vaultService.UpdateDisciplineIndex(ctx, user, true, false)
		} else {
			// Penalize for quitting early (Lazy Tax)
			user.DisciplineIndex -= 2
			if user.DisciplineIndex < 0 {
				user.DisciplineIndex = 0
			}
			// Recalculate price
			user.CurrentPrice = s.vaultService.CalculatePrice(ctx, user)
		}
	}

	// 4. Save updated session
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, err
	}

	// 5. Save updated user (if discipline/price changed)
	if err := s.userRepo.Save(ctx, user); err != nil {
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
