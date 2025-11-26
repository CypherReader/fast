package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, subscription_tier, referral_code, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			subscription_tier = EXCLUDED.subscription_tier,
			referral_code = EXCLUDED.referral_code
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.SubscriptionTier, user.ReferralCode, user.CreatedAt)
	return err
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, subscription_tier, referral_code, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	var subTier string
	var refCode sql.NullString
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &subTier, &refCode, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	if refCode.Valid {
		user.ReferralCode = refCode.String
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password_hash, subscription_tier, referral_code, created_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	var subTier string
	var refCode sql.NullString
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &subTier, &refCode, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	if refCode.Valid {
		user.ReferralCode = refCode.String
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByReferralCode(ctx context.Context, code string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, subscription_tier, referral_code, created_at FROM users WHERE referral_code = $1`
	row := r.db.QueryRowContext(ctx, query, code)

	var user domain.User
	var subTier string
	var refCode sql.NullString
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &subTier, &refCode, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	if refCode.Valid {
		user.ReferralCode = refCode.String
	}
	return &user, nil
}

type PostgresFastingRepository struct {
	db *sql.DB
}

func NewPostgresFastingRepository(db *sql.DB) *PostgresFastingRepository {
	return &PostgresFastingRepository{db: db}
}

func (r *PostgresFastingRepository) Save(ctx context.Context, session *domain.FastingSession) error {
	query := `INSERT INTO fasting_sessions (id, user_id, start_time, end_time, goal_hours, plan_type, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, session.ID, session.UserID, session.StartTime, session.EndTime, session.GoalHours, session.PlanType, session.Status)
	return err
}

func (r *PostgresFastingRepository) Update(ctx context.Context, session *domain.FastingSession) error {
	query := `UPDATE fasting_sessions SET end_time = $1, status = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, session.EndTime, session.Status, session.ID)
	return err
}

func (r *PostgresFastingRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	query := `SELECT id, user_id, start_time, end_time, goal_hours, plan_type, status FROM fasting_sessions WHERE user_id = $1 AND status = 'active'`
	row := r.db.QueryRowContext(ctx, query, userID)

	var s domain.FastingSession
	var planType, status string
	var endTime *time.Time

	err := row.Scan(&s.ID, &s.UserID, &s.StartTime, &endTime, &s.GoalHours, &planType, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.EndTime = endTime
	s.PlanType = domain.FastingPlanType(planType)
	s.Status = domain.FastingStatus(status)
	return &s, nil
}

func (r *PostgresFastingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error) {
	query := `SELECT id, user_id, start_time, end_time, goal_hours, plan_type, status FROM fasting_sessions WHERE user_id = $1 ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.FastingSession
	for rows.Next() {
		var s domain.FastingSession
		var planType, status string
		var endTime *time.Time
		if err := rows.Scan(&s.ID, &s.UserID, &s.StartTime, &endTime, &s.GoalHours, &planType, &status); err != nil {
			return nil, err
		}
		s.EndTime = endTime
		s.PlanType = domain.FastingPlanType(planType)
		s.Status = domain.FastingStatus(status)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

type PostgresKetoRepository struct {
	db *sql.DB
}

func NewPostgresKetoRepository(db *sql.DB) *PostgresKetoRepository {
	return &PostgresKetoRepository{db: db}
}

func (r *PostgresKetoRepository) Save(ctx context.Context, entry *domain.KetoEntry) error {
	query := `INSERT INTO keto_entries (id, user_id, logged_at, ketone_level, acetone_level, source) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.UserID, entry.LoggedAt, entry.KetoneLevel, entry.AcetoneLevel, entry.Source)
	return err
}

func (r *PostgresKetoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KetoEntry, error) {
	query := `SELECT id, user_id, logged_at, ketone_level, acetone_level, source FROM keto_entries WHERE user_id = $1 ORDER BY logged_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.KetoEntry
	for rows.Next() {
		var e domain.KetoEntry
		var source string
		if err := rows.Scan(&e.ID, &e.UserID, &e.LoggedAt, &e.KetoneLevel, &e.AcetoneLevel, &source); err != nil {
			return nil, err
		}
		e.Source = domain.KetoSource(source)
		entries = append(entries, e)
	}
	return entries, nil
}
