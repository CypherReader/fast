package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type PostgresGamificationRepository struct {
	db *sql.DB
}

var _ ports.GamificationRepository = (*PostgresGamificationRepository)(nil)

func NewPostgresGamificationRepository(db *sql.DB) *PostgresGamificationRepository {
	return &PostgresGamificationRepository{db: db}
}

func (r *PostgresGamificationRepository) SaveUserBadge(ctx context.Context, badge *domain.UserBadge) error {
	query := `
		INSERT INTO user_badges (user_id, badge_id, earned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, badge_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, badge.UserID, badge.BadgeID, badge.EarnedAt)
	return err
}

func (r *PostgresGamificationRepository) GetUserBadges(ctx context.Context, userID uuid.UUID) ([]domain.UserBadge, error) {
	query := `
		SELECT user_id, badge_id, earned_at
		FROM user_badges
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []domain.UserBadge
	for rows.Next() {
		var b domain.UserBadge
		var badgeIDStr string
		if err := rows.Scan(&b.UserID, &badgeIDStr, &b.EarnedAt); err != nil {
			return nil, err
		}
		b.BadgeID = domain.BadgeID(badgeIDStr)
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *PostgresGamificationRepository) UpdateUserStreak(ctx context.Context, streak *domain.UserStreak) error {
	query := `
		INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_activity_date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			current_streak = $2,
			longest_streak = $3,
			last_activity_date = $4
	`
	_, err := r.db.ExecContext(ctx, query, streak.UserID, streak.CurrentStreak, streak.LongestStreak, streak.LastActivityDate)
	return err
}

func (r *PostgresGamificationRepository) GetUserStreak(ctx context.Context, userID uuid.UUID) (*domain.UserStreak, error) {
	query := `
		SELECT user_id, current_streak, longest_streak, last_activity_date
		FROM user_streaks
		WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	var s domain.UserStreak
	if err := row.Scan(&s.UserID, &s.CurrentStreak, &s.LongestStreak, &s.LastActivityDate); err != nil {
		if err == sql.ErrNoRows {
			// Return empty streak if not found, or nil
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
