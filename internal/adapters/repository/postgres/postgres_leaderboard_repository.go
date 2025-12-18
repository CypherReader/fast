package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type PostgresLeaderboardRepository struct {
	db *sql.DB
}

var _ ports.LeaderboardRepository = (*PostgresLeaderboardRepository)(nil)

func NewPostgresLeaderboardRepository(db *sql.DB) *PostgresLeaderboardRepository {
	return &PostgresLeaderboardRepository{db: db}
}

func (r *PostgresLeaderboardRepository) GetGlobalLeaderboard(ctx context.Context, limit int) ([]domain.LeaderboardEntry, error) {
	// Aggregate total fasting hours from fasting_sessions
	// Join with users to get name
	// This is a simplified query; in production, you might want a materialized view or cached table
	query := `
		SELECT 
			u.id, 
			u.email, -- Using email as name for now, should be display_name
			COALESCE(SUM(EXTRACT(EPOCH FROM (fs.end_time - fs.start_time))/3600), 0) as total_hours,
			COALESCE(u.discipline_index, 0) as discipline_score
		FROM users u
		LEFT JOIN fasting_sessions fs ON u.id = fs.user_id AND fs.end_time IS NOT NULL
		GROUP BY u.id, u.email, u.discipline_index
		ORDER BY total_hours DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e domain.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.UserName, &e.TotalFastingHours, &e.DisciplineScore); err != nil {
			return nil, err
		}
		e.Rank = rank
		entries = append(entries, e)
		rank++
	}
	return entries, nil
}

func (r *PostgresLeaderboardRepository) GetTribeLeaderboard(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.LeaderboardEntry, error) {
	query := `
		SELECT 
			u.id, 
			u.email,
			COALESCE(SUM(EXTRACT(EPOCH FROM (fs.end_time - fs.start_time))/3600), 0) as total_hours,
			COALESCE(u.discipline_index, 0) as discipline_score
		FROM users u
		LEFT JOIN fasting_sessions fs ON u.id = fs.user_id AND fs.end_time IS NOT NULL
		WHERE u.tribe_id = $1
		GROUP BY u.id, u.email, u.discipline_index
		ORDER BY total_hours DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, tribeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e domain.LeaderboardEntry
		if err := rows.Scan(&e.UserID, &e.UserName, &e.TotalFastingHours, &e.DisciplineScore); err != nil {
			return nil, err
		}
		e.Rank = rank
		entries = append(entries, e)
		rank++
	}
	return entries, nil
}
