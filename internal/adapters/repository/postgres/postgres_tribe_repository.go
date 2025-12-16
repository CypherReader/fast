package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
)

// PostgresTribeRepository implements TribeRepository for PostgreSQL
type PostgresTribeRepository struct {
	db *sql.DB
}

// NewPostgresTribeRepository creates a new PostgreSQL tribe repository
func NewPostgresTribeRepository(db *sql.DB) ports.TribeRepository {
	return &PostgresTribeRepository{db: db}
}

// Save creates a new tribe
func (r *PostgresTribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	categoryJSON, _ := json.Marshal(tribe.Category)

	query := `
		INSERT INTO tribes (
			id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		tribe.ID, tribe.Name, tribe.Slug, tribe.Description,
		tribe.AvatarURL, tribe.CoverPhotoURL, tribe.CreatorID,
		tribe.FastingSchedule, tribe.PrimaryGoal, categoryJSON,
		tribe.Privacy, tribe.Rules, tribe.MemberCount, tribe.ActiveMemberCount,
	)

	return err
}

// Update updates an existing tribe
func (r *PostgresTribeRepository) Update(ctx context.Context, tribe *domain.Tribe) error {
	categoryJSON, _ := json.Marshal(tribe.Category)

	query := `
		UPDATE tribes SET
			description = $1, avatar_url = $2, cover_photo_url = $3,
			category = $4, privacy = $5, rules = $6,
			member_count = $7, active_member_count = $8,
			updated_at = NOW()
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		tribe.Description, tribe.AvatarURL, tribe.CoverPhotoURL,
		categoryJSON, tribe.Privacy, tribe.Rules,
		tribe.MemberCount, tribe.ActiveMemberCount, tribe.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tribe not found")
	}

	return nil
}

// FindByID retrieves a tribe by ID
func (r *PostgresTribeRepository) FindByID(ctx context.Context, id string) (*domain.Tribe, error) {
	query := `
		SELECT id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count,
			created_at, updated_at
		FROM tribes
		WHERE id = $1 AND deleted_at IS NULL
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
		return nil, err
	}

	if len(categoryJSON) > 0 {
		tribe.Category = categoryJSON
	}

	return tribe, nil
}

// FindBySlug retrieves a tribe by slug
func (r *PostgresTribeRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tribe, error) {
	query := `
		SELECT id, name, slug, description, avatar_url, cover_photo_url,
			creator_id, fasting_schedule, primary_goal, category,
			privacy, rules, member_count, active_member_count,
			created_at, updated_at
		FROM tribes
		WHERE slug = $1 AND deleted_at IS NULL
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
		return nil, err
	}

	if len(categoryJSON) > 0 {
		tribe.Category = categoryJSON
	}

	return tribe, nil
}

// List retrieves tribes with filtering and pagination
func (r *PostgresTribeRepository) List(ctx context.Context, query domain.ListTribesQuery) ([]domain.Tribe, int, error) {
	// Build WHERE clause
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argCount := 1

	if query.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argCount, argCount+1))
		searchTerm := "%" + query.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argCount += 2
	}

	if query.FastingSchedule != "" {
		conditions = append(conditions, fmt.Sprintf("fasting_schedule = $%d", argCount))
		args = append(args, query.FastingSchedule)
		argCount++
	}

	if query.PrimaryGoal != "" {
		conditions = append(conditions, fmt.Sprintf("primary_goal = $%d", argCount))
		args = append(args, query.PrimaryGoal)
		argCount++
	}

	if query.Privacy != "" {
		conditions = append(conditions, fmt.Sprintf("privacy = $%d", argCount))
		args = append(args, query.Privacy)
		argCount++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tribes WHERE %s", whereClause)
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Build ORDER BY
	orderBy := "created_at DESC"
	switch query.SortBy {
	case "popular", "members":
		orderBy = "member_count DESC"
	case "active":
		orderBy = "active_member_count DESC"
	case "newest":
		orderBy = "created_at DESC"
	}

	// Defaults
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
		LIMIT $%d OFFSET $%d
	`, whereClause, orderBy, argCount, argCount+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.db.QueryContext(ctx, mainQuery, args...)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}

		if len(categoryJSON) > 0 {
			tribe.Category = categoryJSON
		}

		tribes = append(tribes, tribe)
	}

	return tribes, totalCount, nil
}

// Delete soft deletes a tribe
func (r *PostgresTribeRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tribes SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tribe not found")
	}

	return nil
}

// SaveMembership creates a new membership
func (r *PostgresTribeRepository) SaveMembership(ctx context.Context, membership *domain.TribeMembership) error {
	query := `
		INSERT INTO tribe_memberships (
			id, tribe_id, user_id, role, status, notifications_enabled
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		membership.ID, membership.TribeID, membership.UserID,
		membership.Role, membership.Status, membership.NotificationsEnabled,
	)

	return err
}

// UpdateMembership updates an existing membership
func (r *PostgresTribeRepository) UpdateMembership(ctx context.Context, membership *domain.TribeMembership) error {
	query := `
		UPDATE tribe_memberships SET
			role = $1, status = $2, left_at = $3, notifications_enabled = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		membership.Role, membership.Status, membership.LeftAt,
		membership.NotificationsEnabled, membership.ID,
	)

	return err
}

// FindMembership finds a specific membership
func (r *PostgresTribeRepository) FindMembership(ctx context.Context, tribeID, userID string) (*domain.TribeMembership, error) {
	query := `
		SELECT id, tribe_id, user_id, role, status, joined_at, left_at, notifications_enabled
		FROM tribe_memberships
		WHERE tribe_id = $1 AND user_id = $2
	`

	membership := &domain.TribeMembership{}
	err := r.db.QueryRowContext(ctx, query, tribeID, userID).Scan(
		&membership.ID, &membership.TribeID, &membership.UserID,
		&membership.Role, &membership.Status, &membership.JoinedAt,
		&membership.LeftAt, &membership.NotificationsEnabled,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return membership, nil
}

// GetMembersByTribeID retrieves members of a tribe
func (r *PostgresTribeRepository) GetMembersByTribeID(ctx context.Context, tribeID string, limit, offset int) ([]domain.TribeMember, error) {
	query := `
		SELECT 
			tm.id, tm.tribe_id, tm.user_id, tm.role, tm.status,
			tm.joined_at, tm.left_at, tm.notifications_enabled,
			u.name as user_name, u.email as user_avatar, 0 as user_streak
		FROM tribe_memberships tm
		JOIN users u ON tm.user_id = u.id
		WHERE tm.tribe_id = $1 AND tm.status = 'active'
		ORDER BY tm.joined_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, tribeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []domain.TribeMember{}
	for rows.Next() {
		member := domain.TribeMember{}
		var userName, userAvatar sql.NullString

		err := rows.Scan(
			&member.ID, &member.TribeID, &member.UserID, &member.Role, &member.Status,
			&member.JoinedAt, &member.LeftAt, &member.NotificationsEnabled,
			&userName, &userAvatar, &member.UserStreak,
		)
		if err != nil {
			return nil, err
		}

		if userName.Valid {
			member.UserName = userName.String
		}
		if userAvatar.Valid {
			member.UserAvatar = userAvatar.String
		}

		members = append(members, member)
	}

	return members, nil
}

// GetUserTribes retrieves all tribes a user is a member of
func (r *PostgresTribeRepository) GetUserTribes(ctx context.Context, userID string, status string) ([]domain.Tribe, error) {
	query := `
		SELECT 
			t.id, t.name, t.slug, t.description, t.avatar_url, t.cover_photo_url,
			t.creator_id, t.fasting_schedule, t.primary_goal, t.category,
			t.privacy, t.rules, t.member_count, t.active_member_count,
			t.created_at, t.updated_at
		FROM tribes t
		JOIN tribe_memberships tm ON t.id = tm.tribe_id
		WHERE tm.user_id = $1 AND tm.status = $2 AND t.deleted_at IS NULL
		ORDER BY tm.joined_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if len(categoryJSON) > 0 {
			tribe.Category = categoryJSON
		}

		tribes = append(tribes, tribe)
	}

	return tribes, nil
}

// GetMembershipCount counts active members
func (r *PostgresTribeRepository) GetMembershipCount(ctx context.Context, tribeID string) (int, error) {
	query := `SELECT COUNT(*) FROM tribe_memberships WHERE tribe_id = $1 AND status = 'active'`

	var count int
	err := r.db.QueryRowContext(ctx, query, tribeID).Scan(&count)
	return count, err
}

// DeleteMembership removes a membership
func (r *PostgresTribeRepository) DeleteMembership(ctx context.Context, tribeID, userID string) error {
	query := `
		UPDATE tribe_memberships 
		SET status = 'left', left_at = NOW()
		WHERE tribe_id = $1 AND user_id = $2 AND status = 'active'
	`

	result, err := r.db.ExecContext(ctx, query, tribeID, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("membership not found")
	}

	return nil
}

// GetTribeStats retrieves statistics for a tribe
func (r *PostgresTribeRepository) GetTribeStats(ctx context.Context, tribeID string) (*domain.TribeStats, error) {
	tribe, err := r.FindByID(ctx, tribeID)
	if err != nil {
		return nil, err
	}

	stats := &domain.TribeStats{
		TribeID:              tribeID,
		TotalFasts:           0, // TODO: aggregate from fasting_sessions
		TotalFastingHours:    0, // TODO: aggregate
		AverageMemberStreak:  0, // TODO: calculate
		WeeklyGrowthPercent:  0, // TODO: calculate
		ActiveMembersPercent: 0,
	}

	if tribe.MemberCount > 0 {
		stats.ActiveMembersPercent = float64(tribe.ActiveMemberCount) / float64(tribe.MemberCount) * 100
	}

	return stats, nil
}

// UpdateMemberCounts recalculates member counts
func (r *PostgresTribeRepository) UpdateMemberCounts(ctx context.Context, tribeID string) error {
	// Count total active members
	totalQuery := `SELECT COUNT(*) FROM tribe_memberships WHERE tribe_id = $1 AND status = 'active'`
	var totalCount int
	err := r.db.QueryRowContext(ctx, totalQuery, tribeID).Scan(&totalCount)
	if err != nil {
		return err
	}

	// Count active members (joined recently)
	activeQuery := `
		SELECT COUNT(*) FROM tribe_memberships 
		WHERE tribe_id = $1 AND status = 'active' 
		AND joined_at >= NOW() - INTERVAL '7 days'
	`
	var activeCount int
	err = r.db.QueryRowContext(ctx, activeQuery, tribeID).Scan(&activeCount)
	if err != nil {
		return err
	}

	// Update tribe
	updateQuery := `
		UPDATE tribes 
		SET member_count = $1, active_member_count = $2, updated_at = NOW()
		WHERE id = $3
	`
	_, err = r.db.ExecContext(ctx, updateQuery, totalCount, activeCount, tribeID)
	return err
}
