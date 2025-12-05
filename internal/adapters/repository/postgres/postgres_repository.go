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
		INSERT INTO users (
			id, email, password_hash, name, onboarding_completed, goal, fasting_plan, sex, height_cm, 
			current_weight_lbs, target_weight_lbs, timezone, units, stripe_customer_id, subscription_tier, 
			subscription_status, subscription_id, vault_enabled, trial_ends_at, discipline_index, 
			current_price, vault_deposit, earned_refund, tribe_id, referral_code, signed_contract, 
			push_notifications_enabled, notification_token, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			password_hash = EXCLUDED.password_hash,
			name = EXCLUDED.name,
			onboarding_completed = EXCLUDED.onboarding_completed,
			goal = EXCLUDED.goal,
			fasting_plan = EXCLUDED.fasting_plan,
			sex = EXCLUDED.sex,
			height_cm = EXCLUDED.height_cm,
			current_weight_lbs = EXCLUDED.current_weight_lbs,
			target_weight_lbs = EXCLUDED.target_weight_lbs,
			timezone = EXCLUDED.timezone,
			units = EXCLUDED.units,
			stripe_customer_id = EXCLUDED.stripe_customer_id,
			subscription_tier = EXCLUDED.subscription_tier,
			subscription_status = EXCLUDED.subscription_status,
			subscription_id = EXCLUDED.subscription_id,
			vault_enabled = EXCLUDED.vault_enabled,
			trial_ends_at = EXCLUDED.trial_ends_at,
			discipline_index = EXCLUDED.discipline_index,
			current_price = EXCLUDED.current_price,
			vault_deposit = EXCLUDED.vault_deposit,
			earned_refund = EXCLUDED.earned_refund,
			tribe_id = EXCLUDED.tribe_id,
			referral_code = EXCLUDED.referral_code,
			signed_contract = EXCLUDED.signed_contract,
			push_notifications_enabled = EXCLUDED.push_notifications_enabled,
			notification_token = EXCLUDED.notification_token,
			updated_at = NOW()
	`
	var refCode sql.NullString
	if user.ReferralCode != "" {
		refCode = sql.NullString{String: user.ReferralCode, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.OnboardingCompleted, user.Goal, user.FastingPlan,
		user.Sex, user.HeightCm, user.CurrentWeightLbs, user.TargetWeightLbs, user.Timezone, user.Units,
		user.StripeCustomerID, user.SubscriptionTier, user.SubscriptionStatus, user.SubscriptionID, user.VaultEnabled,
		user.TrialEndsAt, user.DisciplineIndex, user.CurrentPrice, user.VaultDeposit, user.EarnedRefund,
		user.TribeID, refCode, user.SignedContract, user.PushNotificationsEnabled, user.NotificationToken,
		user.CreatedAt, time.Now(),
	)
	return err
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, onboarding_completed, goal, fasting_plan, sex, height_cm, 
		current_weight_lbs, target_weight_lbs, timezone, units, stripe_customer_id, subscription_tier, 
		subscription_status, subscription_id, vault_enabled, trial_ends_at, discipline_index, 
		current_price, vault_deposit, earned_refund, tribe_id, referral_code, signed_contract, 
		push_notifications_enabled, notification_token, created_at, updated_at
		FROM users WHERE email = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, email))
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, onboarding_completed, goal, fasting_plan, sex, height_cm, 
		current_weight_lbs, target_weight_lbs, timezone, units, stripe_customer_id, subscription_tier, 
		subscription_status, subscription_id, vault_enabled, trial_ends_at, discipline_index, 
		current_price, vault_deposit, earned_refund, tribe_id, referral_code, signed_contract, 
		push_notifications_enabled, notification_token, created_at, updated_at
		FROM users WHERE id = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresUserRepository) FindByReferralCode(ctx context.Context, code string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, name, onboarding_completed, goal, fasting_plan, sex, height_cm, 
		current_weight_lbs, target_weight_lbs, timezone, units, stripe_customer_id, subscription_tier, 
		subscription_status, subscription_id, vault_enabled, trial_ends_at, discipline_index, 
		current_price, vault_deposit, earned_refund, tribe_id, referral_code, signed_contract, 
		push_notifications_enabled, notification_token, created_at, updated_at
		FROM users WHERE referral_code = $1
	`
	return r.scanUser(r.db.QueryRowContext(ctx, query, code))
}

func (r *PostgresUserRepository) scanUser(row *sql.Row) (*domain.User, error) {
	var user domain.User
	var subTier string
	var subStatus string
	var refCode sql.NullString
	var tribeID uuid.NullUUID
	var trialEndsAt sql.NullTime
	var updatedAt sql.NullTime

	// Nullable strings
	var name, goal, fastingPlan, sex, timezone, units, stripeCustID, subID, notifToken sql.NullString
	// Nullable floats
	var height, curWeight, targetWeight sql.NullFloat64

	err := row.Scan(
		&user.ID, &user.Email, &user.PasswordHash, &name, &user.OnboardingCompleted, &goal, &fastingPlan,
		&sex, &height, &curWeight, &targetWeight, &timezone, &units, &stripeCustID, &subTier, &subStatus,
		&subID, &user.VaultEnabled, &trialEndsAt, &user.DisciplineIndex, &user.CurrentPrice, &user.VaultDeposit,
		&user.EarnedRefund, &tribeID, &refCode, &user.SignedContract, &user.PushNotificationsEnabled,
		&notifToken, &user.CreatedAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.SubscriptionTier = domain.SubscriptionTier(subTier)
	user.SubscriptionStatus = domain.SubscriptionStatus(subStatus)
	if refCode.Valid {
		user.ReferralCode = refCode.String
	}
	if tribeID.Valid {
		user.TribeID = &tribeID.UUID
	}
	if trialEndsAt.Valid {
		user.TrialEndsAt = &trialEndsAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	}

	// Handle nullable strings
	if name.Valid {
		user.Name = name.String
	}
	if goal.Valid {
		user.Goal = goal.String
	}
	if fastingPlan.Valid {
		user.FastingPlan = fastingPlan.String
	}
	if sex.Valid {
		user.Sex = sex.String
	}
	if timezone.Valid {
		user.Timezone = timezone.String
	}
	if units.Valid {
		user.Units = units.String
	}
	if stripeCustID.Valid {
		user.StripeCustomerID = stripeCustID.String
	}
	if subID.Valid {
		user.SubscriptionID = subID.String
	}
	if notifToken.Valid {
		user.NotificationToken = notifToken.String
	}

	// Handle nullable floats
	if height.Valid {
		user.HeightCm = height.Float64
	}
	if curWeight.Valid {
		user.CurrentWeightLbs = curWeight.Float64
	}
	if targetWeight.Valid {
		user.TargetWeightLbs = targetWeight.Float64
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
