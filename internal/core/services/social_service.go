package services

import (
	"context"
	"errors"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"time"

	"github.com/google/uuid"
)

type SocialService struct {
	repo ports.SocialRepository
}

func NewSocialService(repo ports.SocialRepository) *SocialService {
	return &SocialService{
		repo: repo,
	}
}

func (s *SocialService) AddFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	if userID == friendID {
		return errors.New("cannot add self as friend")
	}

	// Check if already friends (bi-directional check needed in real app, simplified here)
	friends, err := s.repo.FindFriends(ctx, userID)
	if err != nil {
		return err
	}
	for _, f := range friends {
		if f.FriendID == friendID {
			return errors.New("already friends or request pending")
		}
	}

	fn := &domain.FriendNetwork{
		ID:        uuid.New(),
		UserID:    userID,
		FriendID:  friendID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	return s.repo.SaveFriendNetwork(ctx, fn)
}

func (s *SocialService) GetFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendNetwork, error) {
	return s.repo.FindFriends(ctx, userID)
}

func (s *SocialService) CreateTribe(ctx context.Context, userID uuid.UUID, name, description string, isPublic bool) (*domain.Tribe, error) {
	privacy := "public"
	if !isPublic {
		privacy = "private"
	}
	tribe := &domain.Tribe{
		ID:          uuid.New().String(),
		CreatorID:   userID.String(),
		Name:        name,
		Description: description,
		Privacy:     privacy,
		MemberCount: 1, // Owner is first member
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.SaveTribe(ctx, tribe); err != nil {
		return nil, err
	}

	// Log event
	event := &domain.SocialEvent{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: domain.EventTribeJoined, // Creator joins their own tribe
		Data:      `{"tribe_name": "` + name + `"}`,
		CreatedAt: time.Now(),
	}
	_ = s.repo.SaveEvent(ctx, event) // Ignore error for now, non-critical

	return tribe, nil
}

func (s *SocialService) GetTribe(ctx context.Context, tribeID uuid.UUID) (*domain.Tribe, error) {
	return s.repo.FindTribeByID(ctx, tribeID)
}

func (s *SocialService) CreateChallenge(ctx context.Context, userID uuid.UUID, name string, challengeType domain.ChallengeType, goal int, startDate, endDate time.Time) (*domain.FriendChallenge, error) {
	challenge := &domain.FriendChallenge{
		ID:            uuid.New(),
		CreatorID:     userID,
		Name:          name,
		ChallengeType: challengeType,
		Goal:          goal,
		StartDate:     startDate,
		EndDate:       endDate,
		Status:        "active",
		CreatedAt:     time.Now(),
	}

	if err := s.repo.SaveChallenge(ctx, challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}

func (s *SocialService) GetChallenges(ctx context.Context, userID uuid.UUID) ([]domain.FriendChallenge, error) {
	return s.repo.FindChallengesByUserID(ctx, userID)
}

func (s *SocialService) ListTribes(ctx context.Context, limit, offset int) ([]domain.Tribe, error) {
	return s.repo.FindAllTribes(ctx, limit, offset)
}

func (s *SocialService) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.SocialEvent, error) {
	return s.repo.GetFeed(ctx, userID, limit, offset)
}
