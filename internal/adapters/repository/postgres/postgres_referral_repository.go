package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type PostgresReferralRepository struct {
	db *sql.DB
}

func NewPostgresReferralRepository(db *sql.DB) *PostgresReferralRepository {
	return &PostgresReferralRepository{db: db}
}

func (r *PostgresReferralRepository) Save(ctx context.Context, referral *domain.Referral) error {
	query := `
		INSERT INTO referrals (id, referrer_id, referee_id, status, reward_value, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		referral.ID,
		referral.ReferrerID,
		referral.RefereeID,
		referral.Status,
		referral.RewardValue,
		referral.CreatedAt,
	)
	return err
}

func (r *PostgresReferralRepository) FindByRefereeID(ctx context.Context, refereeID uuid.UUID) (*domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referee_id, status, reward_value, created_at, completed_at
		FROM referrals
		WHERE referee_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, refereeID)

	var referral domain.Referral
	var completedAt sql.NullTime

	err := row.Scan(
		&referral.ID,
		&referral.ReferrerID,
		&referral.RefereeID,
		&referral.Status,
		&referral.RewardValue,
		&referral.CreatedAt,
		&completedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if completedAt.Valid {
		referral.CompletedAt = &completedAt.Time
	}

	return &referral, nil
}

func (r *PostgresReferralRepository) FindByReferrerID(ctx context.Context, referrerID uuid.UUID) ([]domain.Referral, error) {
	query := `
		SELECT id, referrer_id, referee_id, status, reward_value, created_at, completed_at
		FROM referrals
		WHERE referrer_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, referrerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var referrals []domain.Referral
	for rows.Next() {
		var referral domain.Referral
		var completedAt sql.NullTime
		if err := rows.Scan(
			&referral.ID,
			&referral.ReferrerID,
			&referral.RefereeID,
			&referral.Status,
			&referral.RewardValue,
			&referral.CreatedAt,
			&completedAt,
		); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			referral.CompletedAt = &completedAt.Time
		}
		referrals = append(referrals, referral)
	}
	return referrals, nil
}

func (r *PostgresReferralRepository) Update(ctx context.Context, referral *domain.Referral) error {
	query := `
		UPDATE referrals
		SET status = $1, completed_at = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, referral.Status, referral.CompletedAt, referral.ID)
	return err
}
