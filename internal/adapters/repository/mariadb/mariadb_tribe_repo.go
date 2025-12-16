package mariadb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
)

type tribeRepository struct {
	db *sql.DB
}

// NewTribeRepository creates a new MariaDB tribe repository
func NewTribeRepository(db *sql.DB) ports.TribeRepository {
	return &tribeRepository{db: db}
}

// Save creates a new tribe in the database
func (r *tribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	categoryJSON, err := json.Marshal(tribe.Category)
	if err != nil {
		return fmt.Errorf("failed to marshal category: %w", err)
	}

	query := `
		INSERT INTO tribes (
			id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		tribe.ID, tribe.Name, tribe.Slug, tribe.Description,
		tribe.AvatarURL, tribe.CoverPhotoURL, tribe.CreatorID,
		tribe.FastingSchedule, tribe.PrimaryGoal, categoryJSON,
		tribe.Privacy, tribe.Rules, tribe.MemberCount, tribe.ActiveMemberCount,
	)

	if err != nil {
		return fmt.Errorf("failed to save tribe: %w", err)
	}

	return nil
}

// Update updates an existing tribe
func (r *tribeRepository) Update(ctx context.Context, tribe *domain.Tribe) error {
	categoryJSON, err := json.Marshal(tribe.Category)
	if err != nil {
		return fmt.Errorf("failed to marshal category: %w", err)
	}

	query := `
		UPDATE tribes SET
			description = ?, avatar_url = ?, cover_photo_url = ?,
			category = ?, privacy = ?, rules = ?,
			member_count = ?, active_member_count = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		tribe.Description, tribe.AvatarURL, tribe.CoverPhotoURL,
		categoryJSON, tribe.Privacy, tribe.Rules,
		tribe.MemberCount, tribe.ActiveMemberCount, tribe.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update tribe: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tribe not found")
	}

	return nil
}

// FindByID retrieves a tribe by its ID
func (r *tribeRepository) FindByID(ctx context.Context, id string) (*domain.Tribe, error) {
	query := `
		SELECT id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count,
			created_at, updated_at
		FROM tribes
		WHERE id = ? AND deleted_at IS NULL
	`

	tribe := &domain.Tribe{}
	var categoryJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tribe.ID, &tribe.Name, &tribe.Slug, &tribe.Description,
		&tribe.AvatarURL, &tribe.CoverPhotoURL, &tribe.CreatorID,
		&tribe.FastingSchedule, &tribe.PrimaryGoal, &categoryJSON,
		&tribe.Privacy, &tribe.Rules, &tribe.MemberCount,
		&tribe.ActiveMemberCount, &tribe.CreatedAt, &tribe.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tribe not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find tribe: %w", err)
	}

	if len(categoryJSON) > 0 {
		tribe.Category = categoryJSON
	}

	return tribe, nil
}

// FindBySlug retrieves a tribe by its slug
func (r *tribeRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tribe, error) {
	query := `
		SELECT id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count,
			created_at, updated_at
		FROM tribes
		WHERE slug = ? AND deleted_at IS NULL
	`

	tribe := &domain.Tribe{}
	var categoryJSON []byte

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&tribe.ID, &tribe.Name, &tribe.Slug, &tribe.Description,
		&tribe.AvatarURL, &tribe.CoverPhotoURL, &tribe.CreatorID,
		&tribe.FastingSchedule, &tribe.PrimaryGoal, &categoryJSON,
		&tribe.Privacy, &tribe.Rules, &tribe.MemberCount,
		&tribe.ActiveMemberCount, &tribe.CreatedAt, &tribe.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tribe not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find tribe: %w", err)
	}

	if len(categoryJSON) > 0 {
		tribe.Category = categoryJSON
	}

	return tribe, nil
}

// List retrieves tribes based on query parameters
func (r *tribeRepository) List(ctx context.Context, query domain.ListTribesQuery) ([]domain.Tribe, int, error) {
	// Build WHERE clause
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	if query.Search != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		searchTerm := "%" + query.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if query.FastingSchedule != "" {
		conditions = append(conditions, "fasting_schedule = ?")
		args = append(args, query.FastingSchedule)
	}

	if query.PrimaryGoal != "" {
		conditions = append(conditions, "primary_goal = ?")
		args = append(args, query.PrimaryGoal)
	}

	if query.Privacy != "" {
		conditions = append(conditions, "privacy = ?")
		args = append(args, query.Privacy)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tribes WHERE %s", whereClause)
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tribes: %w", err)
	}

	// Build ORDER BY clause
	orderBy := "created_at DESC"
	switch query.SortBy {
	case "popular":
		orderBy = "member_count DESC"
	case "active":
		orderBy = "active_member_count DESC"
	case "members":
		orderBy = "member_count DESC"
	case "newest":
		orderBy = "created_at DESC"
	}

	// Set defaults for limit and offset
	if query.Limit <= 0 {
		query.Limit = 20
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	// Build main query
	mainQuery := fmt.Sprintf(`
		SELECT id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count,
			created_at, updated_at
		FROM tribes
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tribes: %w", err)
	}
	defer rows.Close()

	tribes := []domain.Tribe{}
	for rows.Next() {
		tribe := domain.Tribe{}
		var categoryJSON []byte

		err := rows.Scan(
			&tribe.ID, &tribe.Name, &tribe.Slug, &tribe.Description,
			&tribe.AvatarURL, &tribe.CoverPhotoURL, &tribe.CreatorID,
			&tribe.FastingSchedule, &tribe.PrimaryGoal, &categoryJSON,
			&tribe.Privacy, &tribe.Rules, &tribe.MemberCount,
			&tribe.ActiveMemberCount, &tribe.CreatedAt, &tribe.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tribe: %w", err)
		}

		if len(categoryJSON) > 0 {
			tribe.Category = categoryJSON
		}

		tribes = append(tribes, tribe)
	}

	return tribes, totalCount, nil
}

// Delete soft deletes a tribe
func (r *tribeRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tribes SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tribe: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tribe not found")
	}

	return nil
}

// SaveMembership creates a new tribe membership
func (r *tribeRepository) SaveMembership(ctx context.Context, membership *domain.TribeMembership) error {
	query := `
		INSERT INTO tribe_memberships (
			id, tribe_id, user_id, role, status, notifications_enabled
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		membership.ID, membership.TribeID, membership.UserID,
		membership.Role, membership.Status, membership.NotificationsEnabled,
	)

	if err != nil {
		return fmt.Errorf("failed to save membership: %w", err)
	}

	return nil
}

// UpdateMembership updates an existing membership
func (r *tribeRepository) UpdateMembership(ctx context.Context, membership *domain.TribeMembership) error {
	query := `
		UPDATE tribe_memberships SET
			role = ?, status = ?, left_at = ?, notifications_enabled = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		membership.Role, membership.Status, membership.LeftAt,
		membership.NotificationsEnabled, membership.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update membership: %w", err)
	}

	return nil
}

// FindMembership finds a specific membership
func (r *tribeRepository) FindMembership(ctx context.Context, tribeID, userID string) (*domain.TribeMembership, error) {
	query := `
		SELECT id, tribe_id, user_id, role, status, joined_at, left_at, notifications_enabled
		FROM tribe_memberships
		WHERE tribe_id = ? AND user_id = ?
	`

	membership := &domain.TribeMembership{}
	err := r.db.QueryRowContext(ctx, query, tribeID, userID).Scan(
		&membership.ID, &membership.TribeID, &membership.UserID,
		&membership.Role, &membership.Status, &membership.JoinedAt,
		&membership.LeftAt, &membership.NotificationsEnabled,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}

	return membership, nil
}

// GetMembersByTribeID retrieves members of a tribe
func (r *tribeRepository) GetMembersByTribeID(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	query := `
		SELECT 
			tm.id, tm.tribe_id, tm.user_id, tm.role, tm.status,
			tm.joined_at, tm.left_at, tm.notifications_enabled,
			u.name as user_name, u.email as user_avatar, 0 as user_streak
		FROM tribe_memberships tm
		JOIN users u ON tm.user_id = u.id
		WHERE tm.tribe_id = ? AND tm.status = 'active'
		ORDER BY tm.joined_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, tribeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	members := []domain.TribeMember{}
	for rows.Next() {
		member := domain.TribeMember{}
		err := rows.Scan(
			&member.ID, &member.TribeID, &member.UserID, &member.Role, &member.Status,
			&member.JoinedAt, &member.LeftAt, &member.NotificationsEnabled,
			&member.UserName, &member.UserAvatar, &member.UserStreak,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}

	return members, nil
}

// GetUserTribes retrieves all tribes a user is a member of
func (r *tribeRepository) GetUserTribes(ctx context.Context, userID string, status string) ([]domain.Tribe, error) {
	query := `
		SELECT 
			t.id, t.name, t.slug, t.description, t.avatar_url, t.cover_photo_url,
			t.creator_id, t.fasting_schedule, t.primary_goal, t.category,
			t.privacy, t.rules, t.member_count, t.active_member_count,
			t.created_at, t.updated_at
		FROM tribes t
		JOIN tribe_memberships tm ON t.id = tm.tribe_id
		WHERE tm.user_id = ? AND tm.status = ? AND t.deleted_at IS NULL
		ORDER BY tm.joined_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tribes: %w", err)
	}
	defer rows.Close()

	tribes := []domain.Tribe{}
	for rows.Next() {
		tribe := domain.Tribe{}
		var categoryJSON []byte

		err := rows.Scan(
			&tribe.ID, &tribe.Name, &tribe.Slug, &tribe.Description,
			&tribe.AvatarURL, &tribe.CoverPhotoURL, &tribe.CreatorID,
			&tribe.FastingSchedule, &tribe.PrimaryGoal, &categoryJSON,
			&tribe.Privacy, &tribe.Rules, &tribe.MemberCount,
			&tribe.ActiveMemberCount, &tribe.CreatedAt, &tribe.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tribe: %w", err)
		}

		if len(categoryJSON) > 0 {
			tribe.Category = categoryJSON
		}

		tribes = append(tribes, tribe)
	}

	return tribes, nil
}

// GetMembershipCount counts active members in a tribe
func (r *tribeRepository) GetMembershipCount(ctx context.Context, tribeID string) (int, error) {
	query := `SELECT COUNT(*) FROM tribe_memberships WHERE tribe_id = ? AND status = 'active'`

	var count int
	err := r.db.QueryRowContext(ctx, query, tribeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count members: %w", err)
	}

	return count, nil
}

// DeleteMembership removes a membership (soft delete via status)
func (r *tribeRepository) DeleteMembership(ctx context.Context, tribeID, userID string) error {
	query := `
		UPDATE tribe_memberships 
		SET status = 'left', left_at = CURRENT_TIMESTAMP
		WHERE tribe_id = ? AND user_id = ? AND status = 'active'
	`

	result, err := r.db.ExecContext(ctx, query, tribeID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete membership: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("membership not found")
	}

	return nil
}

// GetTribeStats retrieves statistics for a tribe
func (r *tribeRepository) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	// For MVP, return basic stats from tribe table
	// Future: aggregate from fasting sessions, posts, etc.
	tribe, err := r.FindByID(ctx, tribeID)
	if err != nil {
		return nil, err
	}

	stats := &domain.TribeStats{
		TribeID:              tribeID,
		TotalFasts:           0, // TODO: aggregate from fasting_sessions
		TotalFastingHours:    0, // TODO: aggregate from fasting_sessions
		AverageMemberStreak:  0, // TODO: calculate from user stats
		WeeklyGrowthPercent:  0, // TODO: calculate from memberships
		ActiveMembersPercent: float64(tribe.ActiveMemberCount) / float64(tribe.MemberCount) * 100,
	}

	return stats, nil
}

// UpdateMemberCounts recalculates and updates member counts for a tribe
func (r *tribeRepository) UpdateMemberCounts(ctx context.Context, tribeID string) error {
	// Count total active members
	totalQuery := `SELECT COUNT(*) FROM tribe_memberships WHERE tribe_id = ? AND status = 'active'`
	var totalCount int
	err := r.db.QueryRowContext(ctx, totalQuery, tribeID).Scan(&totalCount)
	if err != nil {
		return fmt.Errorf("failed to count total members: %w", err)
	}

	// Count active members (joined in last 7 days OR has activity in last 7 days)
	// For MVP, we'll just use members who joined recently
	activeQuery := `
		SELECT COUNT(*) FROM tribe_memberships 
		WHERE tribe_id = ? AND status = 'active' 
		AND joined_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
	`
	var activeCount int
	err = r.db.QueryRowContext(ctx, activeQuery, tribeID).Scan(&activeCount)
	if err != nil {
		return fmt.Errorf("failed to count active members: %w", err)
	}

	// Update tribe
	updateQuery := `
		UPDATE tribes 
		SET member_count = ?, active_member_count = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, updateQuery, totalCount, activeCount, tribeID)
	if err != nil {
		return fmt.Errorf("failed to update member counts: %w", err)
	}

	return nil
}
