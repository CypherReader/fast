package mariadb

import (
	"context"
	"database/sql"
	"errors"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash, subscription_tier, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, user.ID.String(), user.Email, user.PasswordHash, user.SubscriptionTier, user.CreatedAt)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, subscription_tier, created_at FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	var idStr string
	var subTier string
	err := row.Scan(&idStr, &user.Email, &user.PasswordHash, &subTier, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.ID = uuid.MustParse(idStr)
	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `SELECT id, email, password_hash, subscription_tier, created_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var user domain.User
	var idStr string
	var subTier string
	err := row.Scan(&idStr, &user.Email, &user.PasswordHash, &subTier, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	user.ID = uuid.MustParse(idStr)
	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	return &user, nil
}

type FastingRepository struct {
	db *sql.DB
}

func NewFastingRepository(db *sql.DB) *FastingRepository {
	return &FastingRepository{db: db}
}

func (r *FastingRepository) Save(ctx context.Context, session *domain.FastingSession) error {
	query := `INSERT INTO fasting_sessions (id, user_id, start_time, end_time, goal_hours, plan_type, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, session.ID.String(), session.UserID.String(), session.StartTime, session.EndTime, session.GoalHours, session.PlanType, session.Status)
	return err
}

func (r *FastingRepository) Update(ctx context.Context, session *domain.FastingSession) error {
	query := `UPDATE fasting_sessions SET end_time = ?, status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, session.EndTime, session.Status, session.ID.String())
	return err
}

func (r *FastingRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.FastingSession, error) {
	query := `SELECT id, user_id, start_time, end_time, goal_hours, plan_type, status FROM fasting_sessions WHERE user_id = ? AND status = 'active'`
	row := r.db.QueryRowContext(ctx, query, userID.String())

	var s domain.FastingSession
	var idStr, userIDStr, planType, status string
	var endTime *time.Time

	err := row.Scan(&idStr, &userIDStr, &s.StartTime, &endTime, &s.GoalHours, &planType, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	s.ID = uuid.MustParse(idStr)
	s.UserID = uuid.MustParse(userIDStr)
	s.EndTime = endTime
	s.PlanType = domain.FastingPlanType(planType)
	s.Status = domain.FastingStatus(status)
	return &s, nil
}

func (r *FastingRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FastingSession, error) {
	query := `SELECT id, user_id, start_time, end_time, goal_hours, plan_type, status FROM fasting_sessions WHERE user_id = ? ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.FastingSession
	for rows.Next() {
		var s domain.FastingSession
		var idStr, userIDStr, planType, status string
		var endTime *time.Time
		if err := rows.Scan(&idStr, &userIDStr, &s.StartTime, &endTime, &s.GoalHours, &planType, &status); err != nil {
			return nil, err
		}
		s.ID = uuid.MustParse(idStr)
		s.UserID = uuid.MustParse(userIDStr)
		s.EndTime = endTime
		s.PlanType = domain.FastingPlanType(planType)
		s.Status = domain.FastingStatus(status)
		sessions = append(sessions, s)
	}
	return sessions, nil
}

type KetoRepository struct {
	db *sql.DB
}

func NewKetoRepository(db *sql.DB) *KetoRepository {
	return &KetoRepository{db: db}
}

func (r *KetoRepository) Save(ctx context.Context, entry *domain.KetoEntry) error {
	query := `INSERT INTO keto_entries (id, user_id, logged_at, ketone_level, acetone_level, source) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, entry.ID.String(), entry.UserID.String(), entry.LoggedAt, entry.KetoneLevel, entry.AcetoneLevel, entry.Source)
	return err
}

func (r *KetoRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.KetoEntry, error) {
	query := `SELECT id, user_id, logged_at, ketone_level, acetone_level, source FROM keto_entries WHERE user_id = ? ORDER BY logged_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.KetoEntry
	for rows.Next() {
		var e domain.KetoEntry
		var idStr, userIDStr, source string
		if err := rows.Scan(&idStr, &userIDStr, &e.LoggedAt, &e.KetoneLevel, &e.AcetoneLevel, &source); err != nil {
			return nil, err
		}
		e.ID = uuid.MustParse(idStr)
		e.UserID = uuid.MustParse(userIDStr)
		e.Source = domain.KetoSource(source)
		entries = append(entries, e)
	}
	return entries, nil
}
