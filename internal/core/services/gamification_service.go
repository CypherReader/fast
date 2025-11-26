package services

import (
	"context"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type GamificationService struct {
	repo        ports.GamificationRepository
	fastingRepo ports.FastingRepository // Needed to check total hours etc.
}

func NewGamificationService(repo ports.GamificationRepository, fastingRepo ports.FastingRepository) *GamificationService {
	return &GamificationService{
		repo:        repo,
		fastingRepo: fastingRepo,
	}
}

func (s *GamificationService) GetUserGamificationProfile(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, []domain.UserBadge, error) {
	streak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if streak == nil {
		streak = &domain.UserStreak{UserID: userID}
	}

	badges, err := s.repo.GetUserBadges(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Populate static badge info
	for i := range badges {
		if info, ok := domain.Badges[badges[i].BadgeID]; ok {
			badges[i].BadgeInfo = &info
		}
	}

	return streak, badges, nil
}

func (s *GamificationService) UpdateStreak(ctx context.Context, userID uuid.UUID) error {
	streak, err := s.repo.GetUserStreak(ctx, userID)
	if err != nil {
		return err
	}
	if streak == nil {
		streak = &domain.UserStreak{UserID: userID}
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var lastActivityDate time.Time
	if !streak.LastActivityDate.IsZero() {
		lastActivityDate = time.Date(streak.LastActivityDate.Year(), streak.LastActivityDate.Month(), streak.LastActivityDate.Day(), 0, 0, 0, 0, streak.LastActivityDate.Location())
	}

	if lastActivityDate.Equal(today) {
		// Already updated for today
		return nil
	}

	if lastActivityDate.AddDate(0, 0, 1).Equal(today) {
		// Consecutive day
		streak.CurrentStreak++
	} else {
		// Streak broken (or first time)
		streak.CurrentStreak = 1
	}

	if streak.CurrentStreak > streak.LongestStreak {
		streak.LongestStreak = streak.CurrentStreak
	}
	streak.LastActivityDate = now

	return s.repo.UpdateUserStreak(ctx, streak)
}

func (s *GamificationService) CheckAndAwardBadges(ctx context.Context, userID uuid.UUID, eventType string, data interface{}) error {
	// 1. Get existing badges to avoid re-awarding
	existingBadges, err := s.repo.GetUserBadges(ctx, userID)
	if err != nil {
		return err
	}
	hasBadge := func(id domain.BadgeID) bool {
		for _, b := range existingBadges {
			if b.BadgeID == id {
				return true
			}
		}
		return false
	}

	// 2. Check criteria
	// Example: First Fast
	if eventType == "fast_completed" && !hasBadge(domain.BadgeFirstFast) {
		if err := s.awardBadge(ctx, userID, domain.BadgeFirstFast); err != nil {
			return err
		}
	}

	// Example: 100 Hours (Need to query total hours from fasting repo, or pass it in data)
	// For simplicity, let's assume we query it if not passed
	// ...

	// Example: Streak Badges
	streak, err := s.repo.GetUserStreak(ctx, userID)
	if err == nil && streak != nil {
		if streak.CurrentStreak >= 3 && !hasBadge(domain.BadgeStreak3) {
			if err := s.awardBadge(ctx, userID, domain.BadgeStreak3); err != nil {
				return err
			}
		}
		if streak.CurrentStreak >= 7 && !hasBadge(domain.BadgeStreak7) {
			if err := s.awardBadge(ctx, userID, domain.BadgeStreak7); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *GamificationService) awardBadge(ctx context.Context, userID uuid.UUID, badgeID domain.BadgeID) error {
	badge := &domain.UserBadge{
		UserID:   userID,
		BadgeID:  badgeID,
		EarnedAt: time.Now(),
	}
	return s.repo.SaveUserBadge(ctx, badge)
}
