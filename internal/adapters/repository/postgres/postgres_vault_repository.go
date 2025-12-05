package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type PostgresVaultRepository struct {
	db *sql.DB
}

func NewPostgresVaultRepository(db *sql.DB) *PostgresVaultRepository {
	return &PostgresVaultRepository{db: db}
}

func (r *PostgresVaultRepository) Save(ctx context.Context, vault *domain.VaultParticipation) error {
	query := `
		INSERT INTO vault_participations (
			id, user_id, month_start, month_end, deposit_amount, fasts_completed, amount_recovered, 
			refund_processed, refund_date, opted_in, forfeited_amount, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			fasts_completed = EXCLUDED.fasts_completed,
			amount_recovered = EXCLUDED.amount_recovered,
			refund_processed = EXCLUDED.refund_processed,
			refund_date = EXCLUDED.refund_date,
			opted_in = EXCLUDED.opted_in,
			forfeited_amount = EXCLUDED.forfeited_amount,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		vault.ID, vault.UserID, vault.MonthStart, vault.MonthEnd, vault.DepositAmount, vault.FastsCompleted,
		vault.AmountRecovered, vault.RefundProcessed, vault.RefundDate, vault.OptedIn, vault.ForfeitedAmount,
		vault.CreatedAt, time.Now(),
	)
	return err
}

func (r *PostgresVaultRepository) FindByUserIDAndMonth(ctx context.Context, userID uuid.UUID, monthStart time.Time) (*domain.VaultParticipation, error) {
	query := `
		SELECT id, user_id, month_start, month_end, deposit_amount, fasts_completed, amount_recovered, 
		refund_processed, refund_date, opted_in, forfeited_amount, created_at, updated_at
		FROM vault_participations WHERE user_id = $1 AND month_start = $2
	`
	row := r.db.QueryRowContext(ctx, query, userID, monthStart)

	var v domain.VaultParticipation
	var refundDate *time.Time

	err := row.Scan(
		&v.ID, &v.UserID, &v.MonthStart, &v.MonthEnd, &v.DepositAmount, &v.FastsCompleted, &v.AmountRecovered,
		&v.RefundProcessed, &refundDate, &v.OptedIn, &v.ForfeitedAmount, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	v.RefundDate = refundDate
	return &v, nil
}
