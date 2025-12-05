package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type PostgresProgressRepository struct {
	db *sql.DB
}

func NewPostgresProgressRepository(db *sql.DB) *PostgresProgressRepository {
	return &PostgresProgressRepository{db: db}
}

// Weight Logs
func (r *PostgresProgressRepository) SaveWeightLog(ctx context.Context, log *domain.WeightLog) error {
	query := `
		INSERT INTO weight_logs (id, user_id, weight_lbs, weight_kg, logged_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.WeightLbs, log.WeightKg, log.LoggedAt, log.CreatedAt)
	return err
}

func (r *PostgresProgressRepository) GetWeightHistory(ctx context.Context, userID uuid.UUID, days int) ([]domain.WeightLog, error) {
	query := `
		SELECT id, user_id, weight_lbs, weight_kg, logged_at, created_at
		FROM weight_logs
		WHERE user_id = $1 AND logged_at >= NOW() - make_interval(days => $2)
		ORDER BY logged_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.WeightLog
	for rows.Next() {
		var l domain.WeightLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.WeightLbs, &l.WeightKg, &l.LoggedAt, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// Hydration Logs
func (r *PostgresProgressRepository) SaveHydrationLog(ctx context.Context, log *domain.HydrationLog) error {
	query := `
		INSERT INTO hydration_logs (id, user_id, glasses_count, logged_date, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, logged_date) DO UPDATE SET
			glasses_count = EXCLUDED.glasses_count
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.GlassesCount, log.LoggedDate, log.CreatedAt)
	return err
}

func (r *PostgresProgressRepository) GetHydrationLog(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.HydrationLog, error) {
	query := `SELECT id, user_id, glasses_count, logged_date, created_at FROM hydration_logs WHERE user_id = $1 AND logged_date = $2`
	row := r.db.QueryRowContext(ctx, query, userID, date)

	var l domain.HydrationLog
	err := row.Scan(&l.ID, &l.UserID, &l.GlassesCount, &l.LoggedDate, &l.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

// Activity Logs
func (r *PostgresProgressRepository) SaveActivityLog(ctx context.Context, log *domain.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (id, user_id, steps, distance_km, calories_burned, logged_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, logged_date) DO UPDATE SET
			steps = EXCLUDED.steps,
			distance_km = EXCLUDED.distance_km,
			calories_burned = EXCLUDED.calories_burned
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.Steps, log.DistanceKm, log.CaloriesBurned, log.LoggedDate, log.CreatedAt)
	return err
}

func (r *PostgresProgressRepository) GetActivityStats(ctx context.Context, userID uuid.UUID, days int) ([]domain.ActivityLog, error) {
	query := `
		SELECT id, user_id, steps, distance_km, calories_burned, logged_date, created_at
		FROM activity_logs
		WHERE user_id = $1 AND logged_date >= CURRENT_DATE - make_interval(days => $2)
		ORDER BY logged_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.ActivityLog
	for rows.Next() {
		var l domain.ActivityLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Steps, &l.DistanceKm, &l.CaloriesBurned, &l.LoggedDate, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
