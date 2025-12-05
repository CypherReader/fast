package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (r *PostgresSubscriptionRepository) Save(ctx context.Context, sub *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, user_id, plan_type, subscription_price, stripe_subscription_id, status, 
			current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			plan_type = EXCLUDED.plan_type,
			subscription_price = EXCLUDED.subscription_price,
			stripe_subscription_id = EXCLUDED.stripe_subscription_id,
			status = EXCLUDED.status,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end,
			cancel_at_period_end = EXCLUDED.cancel_at_period_end,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.UserID, sub.PlanType, sub.SubscriptionPrice, sub.StripeSubscriptionID, sub.Status,
		sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CancelAtPeriodEnd, sub.CreatedAt, time.Now(),
	)
	return err
}

func (r *PostgresSubscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, subscription_price, stripe_subscription_id, status, 
		current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		FROM subscriptions WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var sub domain.Subscription
	var status string
	var start, end *time.Time
	var stripeSubID sql.NullString

	err := row.Scan(
		&sub.ID, &sub.UserID, &sub.PlanType, &sub.SubscriptionPrice, &stripeSubID, &status,
		&start, &end, &sub.CancelAtPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no subscription found
		}
		return nil, err
	}

	sub.Status = domain.SubscriptionStatus(status)
	sub.CurrentPeriodStart = start
	sub.CurrentPeriodEnd = end
	if stripeSubID.Valid {
		sub.StripeSubscriptionID = stripeSubID.String
	}

	return &sub, nil
}

func (r *PostgresSubscriptionRepository) FindByStripeSubscriptionID(ctx context.Context, stripeSubID string) (*domain.Subscription, error) {
	query := `
		SELECT id, user_id, plan_type, subscription_price, stripe_subscription_id, status, 
		current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		FROM subscriptions WHERE stripe_subscription_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, stripeSubID)

	var sub domain.Subscription
	var status string
	var start, end *time.Time
	var dbStripeSubID sql.NullString

	err := row.Scan(
		&sub.ID, &sub.UserID, &sub.PlanType, &sub.SubscriptionPrice, &dbStripeSubID, &status,
		&start, &end, &sub.CancelAtPeriodEnd, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	sub.Status = domain.SubscriptionStatus(status)
	sub.CurrentPeriodStart = start
	sub.CurrentPeriodEnd = end
	if dbStripeSubID.Valid {
		sub.StripeSubscriptionID = dbStripeSubID.String
	}

	return &sub, nil
}
