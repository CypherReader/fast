package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type TribeService struct {
	repo ports.TribeRepository
}

func NewTribeService(repo ports.TribeRepository) *TribeService {
	return &TribeService{
		repo: repo,
	}
}

// CreateTribe creates a new tribe
func (s *TribeService) CreateTribe(ctx context.Context, userID string, req domain.CreateTribeRequest) (*domain.Tribe, error) {
	// Generate ID and slug
	tribeID := uuid.New().String()
	slug := generateSlug(req.Name)

	// Check if slug already exists
	existing, _ := s.repo.FindBySlug(ctx, slug)
	if existing != nil {
		// Append random suffix to make it unique
		slug = fmt.Sprintf("%s-%s", slug, uuid.New().String()[:8])
	}

	// Marshal category
	var categoryJSON []byte
	if len(req.Category) > 0 {
		category := req.Category
		categoryJSON, _ = json.Marshal(category)
	}

	// Create tribe
	tribe := &domain.Tribe{
		ID:                tribeID,
		Name:              req.Name,
		Slug:              slug,
		Description:       req.Description,
		CreatorID:         userID,
		FastingSchedule:   req.FastingSchedule,
		PrimaryGoal:       req.PrimaryGoal,
		Category:          categoryJSON,
		Privacy:           req.Privacy,
		MemberCount:       1, // Creator is first member
		ActiveMemberCount: 1,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if req.AvatarURL != "" {
		tribe.AvatarURL = &req.AvatarURL
	}
	if req.CoverPhotoURL != "" {
		tribe.CoverPhotoURL = &req.CoverPhotoURL
	}
	if req.Rules != "" {
		tribe.Rules = &req.Rules
	}

	// Save tribe
	if err := s.repo.Save(ctx, tribe); err != nil {
		return nil, fmt.Errorf("failed to create tribe: %w", err)
	}

	// Add creator as first member with creator role
	membership := &domain.TribeMembership{
		ID:                   uuid.New().String(),
		TribeID:              tribeID,
		UserID:               userID,
		Role:                 "creator",
		Status:               "active",
		JoinedAt:             time.Now(),
		NotificationsEnabled: true,
	}

	if err := s.repo.SaveMembership(ctx, membership); err != nil {
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	return tribe, nil
}

// GetTribe retrieves a tribe by ID with membership info for current user
func (s *TribeService) GetTribe(ctx context.Context, tribeID string, currentUserID *string) (*domain.Tribe, error) {
	tribe, err := s.repo.FindByID(ctx, tribeID)
	if err != nil {
		return nil, err
	}

	// If user is logged in, check membership status
	if currentUserID != nil && *currentUserID != "" {
		membership, _ := s.repo.FindMembership(ctx, tribeID, *currentUserID)
		if membership != nil && membership.Status == "active" {
			tribe.IsJoined = true
			tribe.UserRole = membership.Role
		}
	}

	return tribe, nil
}

// UpdateTribe updates tribe information (only creator can update)
func (s *TribeService) UpdateTribe(ctx context.Context, tribeID, userID string, req domain.UpdateTribeRequest) (*domain.Tribe, error) {
	// Get existing tribe
	tribe, err := s.repo.FindByID(ctx, tribeID)
	if err != nil {
		return nil, err
	}

	// Check authorization - only creator can update
	if tribe.CreatorID != userID {
		return nil, fmt.Errorf("unauthorized: only tribe creator can update tribe")
	}

	// Update fields
	if req.Description != nil {
		tribe.Description = *req.Description
	}
	if req.Category != nil {
		categoryJSON, _ := json.Marshal(req.Category)
		tribe.Category = categoryJSON
	}
	if req.Privacy != nil {
		tribe.Privacy = *req.Privacy
	}
	if req.Rules != nil {
		tribe.Rules = req.Rules
	}
	if req.AvatarURL != nil {
		tribe.AvatarURL = req.AvatarURL
	}
	if req.CoverPhotoURL != nil {
		tribe.CoverPhotoURL = req.CoverPhotoURL
	}

	// Save updates
	if err := s.repo.Update(ctx, tribe); err != nil {
		return nil, fmt.Errorf("failed to update tribe: %w", err)
	}

	return tribe, nil
}

// DeleteTribe soft deletes a tribe (only creator can delete)
func (s *TribeService) DeleteTribe(ctx context.Context, tribeID, userID string) error {
	// Get tribe
	tribe, err := s.repo.FindByID(ctx, tribeID)
	if err != nil {
		return err
	}

	// Check authorization
	if tribe.CreatorID != userID {
		return fmt.Errorf("unauthorized: only tribe creator can delete tribe")
	}

	return s.repo.Delete(ctx, tribeID)
}

// ListTribes retrieves a list of tribes with optional filters
func (s *TribeService) ListTribes(ctx context.Context, query domain.ListTribesQuery, currentUserID *string) ([]domain.Tribe, int, error) {
	tribes, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// If user is logged in, check membership status for each tribe
	if currentUserID != nil && *currentUserID != "" {
		for i := range tribes {
			membership, _ := s.repo.FindMembership(ctx, tribes[i].ID, *currentUserID)
			if membership != nil && membership.Status == "active" {
				tribes[i].IsJoined = true
				tribes[i].UserRole = membership.Role
			}
		}
	}

	return tribes, total, nil
}

// JoinTribe adds a user as a member of a tribe
func (s *TribeService) JoinTribe(ctx context.Context, tribeID, userID string) error {
	// Get tribe
	tribe, err := s.repo.FindByID(ctx, tribeID)
	if err != nil {
		return err
	}

	// Check if already a member
	existing, _ := s.repo.FindMembership(ctx, tribeID, userID)
	if existing != nil {
		if existing.Status == "active" {
			return fmt.Errorf("already a member of this tribe")
		}
		// Reactivate if previously left
		existing.Status = "active"
		existing.LeftAt = nil
		return s.repo.UpdateMembership(ctx, existing)
	}

	// For private tribes, would need approval (Phase 2 feature)
	// For now, allow immediate join for all tribes
	status := "active"
	if tribe.Privacy == "private" {
		status = "pending" // Will need approval
	}

	// Create membership
	membership := &domain.TribeMembership{
		ID:                   uuid.New().String(),
		TribeID:              tribeID,
		UserID:               userID,
		Role:                 "member",
		Status:               status,
		JoinedAt:             time.Now(),
		NotificationsEnabled: true,
	}

	if err := s.repo.SaveMembership(ctx, membership); err != nil {
		return fmt.Errorf("failed to join tribe: %w", err)
	}

	// Update member counts
	if err := s.repo.UpdateMemberCounts(ctx, tribeID); err != nil {
		return fmt.Errorf("failed to update member counts: %w", err)
	}

	return nil
}

// LeaveTribe removes a user from a tribe
func (s *TribeService) LeaveTribe(ctx context.Context, tribeID, userID string) error {
	// Get membership
	membership, err := s.repo.FindMembership(ctx, tribeID, userID)
	if err != nil {
		return err
	}
	if membership == nil {
		return fmt.Errorf("not a member of this tribe")
	}

	// Don't allow creator to leave (must delete tribe instead)
	if membership.Role == "creator" {
		return fmt.Errorf("tribe creator cannot leave; delete the tribe instead")
	}

	// Remove membership
	if err := s.repo.DeleteMembership(ctx, tribeID, userID); err != nil {
		return fmt.Errorf("failed to leave tribe: %w", err)
	}

	// Update member counts
	if err := s.repo.UpdateMemberCounts(ctx, tribeID); err != nil {
		return fmt.Errorf("failed to update member counts: %w", err)
	}

	return nil
}

// GetTribeMembers retrieves members of a tribe
func (s *TribeService) GetTribeMembers(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.GetMembersByTribeID(ctx, tribeID, limit, offset)
}

// GetMyTribes retrieves all tribes the user is a member of
func (s *TribeService) GetMyTribes(ctx context.Context, userID string) ([]domain.Tribe, error) {
	return s.repo.GetUserTribes(ctx, userID, "active")
}

// GetTribeStats retrieves statistics for a tribe
func (s *TribeService) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	return s.repo.GetTribeStats(ctx, tribeID)
}

// Helper functions

// generateSlug creates a URL-friendly slug from a tribe name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	// Remove consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Limit length to 60 characters
	if len(slug) > 60 {
		slug = slug[:60]
	}

	return slug
}
