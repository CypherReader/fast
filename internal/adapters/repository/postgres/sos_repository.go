package postgres

import (
	"context"
	"database/sql"
	"time"

	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type PostgresSOSRepository struct {
	db *sql.DB
}

func NewPostgresSOSRepository(db *sql.DB) *PostgresSOSRepository {
	return &PostgresSOSRepository{db: db}
}

// Save creates a new SOS flare
func (r *PostgresSOSRepository) Save(ctx context.Context, sos *domain.SOSFlare) error {
	query := `
		INSERT INTO sos_flares (
			id, user_id, fasting_id, tribe_id, description, hours_fasted,
			status, hype_count, is_anonymous, cortex_responded, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		sos.ID,
		sos.UserID,
		sos.FastingID,
		sos.TribeID,
		sos.Description,
		sos.HoursFasted,
		sos.Status,
		sos.HypeCount,
		sos.IsAnonymous,
		sos.CortexResponded,
		sos.CreatedAt,
	)

	return err
}

// FindByID retrieves an SOS by ID
func (r *PostgresSOSRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SOSFlare, error) {
	query := `
		SELECT id, user_id, fasting_id, tribe_id, description, hours_fasted,
		       status, hype_count, is_anonymous, cortex_responded, created_at, resolved_at
		FROM sos_flares
		WHERE id = $1
	`

	var sos domain.SOSFlare
	var tribeID sql.NullString
	var resolvedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sos.ID,
		&sos.UserID,
		&sos.FastingID,
		&tribeID,
		&sos.Description,
		&sos.HoursFasted,
		&sos.Status,
		&sos.HypeCount,
		&sos.IsAnonymous,
		&sos.CortexResponded,
		&sos.CreatedAt,
		&resolvedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if tribeID.Valid {
		tid, _ := uuid.Parse(tribeID.String)
		sos.TribeID = &tid
	}

	if resolvedAt.Valid {
		sos.ResolvedAt = &resolvedAt.Time
	}

	return &sos, nil
}

// FindActiveByUserID finds the active SOS for a user
func (r *PostgresSOSRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.SOSFlare, error) {
	query := `
		SELECT id, user_id, fasting_id, tribe_id, description, hours_fasted,
		       status, hype_count, is_anonymous, cortex_responded, created_at, resolved_at
		FROM sos_flares
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var sos domain.SOSFlare
	var tribeID sql.NullString
	var resolvedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&sos.ID,
		&sos.UserID,
		&sos.FastingID,
		&tribeID,
		&sos.Description,
		&sos.HoursFasted,
		&sos.Status,
		&sos.HypeCount,
		&sos.IsAnonymous,
		&sos.CortexResponded,
		&sos.CreatedAt,
		&resolvedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if tribeID.Valid {
		tid, _ := uuid.Parse(tribeID.String)
		sos.TribeID = &tid
	}

	if resolvedAt.Valid {
		sos.ResolvedAt = &resolvedAt.Time
	}

	return &sos, nil
}

// UpdateStatus updates the status of an SOS
func (r *PostgresSOSRepository) UpdateStatus(ctx context.Context, sosID uuid.UUID, status domain.SOSStatus) error {
	query := `
		UPDATE sos_flares
		SET status = $1, resolved_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, status, now, sosID)
	return err
}

// UpdateCortexResponse marks that Cortex has responded to this SOS
func (r *PostgresSOSRepository) UpdateCortexResponse(ctx context.Context, sosID uuid.UUID) error {
	query := `
		UPDATE sos_flares
		SET cortex_responded = true
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, sosID)
	return err
}

// SaveHypeResponse saves a hype response
func (r *PostgresSOSRepository) SaveHypeResponse(ctx context.Context, hype *domain.HypeResponse) error {
	query := `
		INSERT INTO hype_responses (id, sos_id, from_user_id, from_name, message, emoji, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		hype.ID,
		hype.SOSID,
		hype.FromUserID,
		hype.FromName,
		hype.Message,
		hype.Emoji,
		hype.CreatedAt,
	)

	return err
}

// GetHypeResponses retrieves all hype responses for an SOS
func (r *PostgresSOSRepository) GetHypeResponses(ctx context.Context, sosID uuid.UUID) ([]domain.HypeResponse, error) {
	query := `
		SELECT id, sos_id, from_user_id, from_name, message, emoji, created_at
		FROM hype_responses
		WHERE sos_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, sosID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hypes []domain.HypeResponse
	for rows.Next() {
		var hype domain.HypeResponse
		var message sql.NullString

		err := rows.Scan(
			&hype.ID,
			&hype.SOSID,
			&hype.FromUserID,
			&hype.FromName,
			&message,
			&hype.Emoji,
			&hype.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if message.Valid {
			hype.Message = message.String
		}

		hypes = append(hypes, hype)
	}

	return hypes, rows.Err()
}

// GetUserHypeCount returns the number of hypes a user has sent since a given time
func (r *PostgresSOSRepository) GetUserHypeCount(ctx context.Context, userID uuid.UUID, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM hype_responses
		WHERE from_user_id = $1 AND created_at >= $2
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID, since).Scan(&count)
	return count, err
}

// FindAllActive returns all active SOS flares (for cron job)
func (r *PostgresSOSRepository) FindAllActive(ctx context.Context) ([]*domain.SOSFlare, error) {
	query := `
		SELECT id, user_id, fasting_id, tribe_id, description, hours_fasted,
		       status, hype_count, is_anonymous, cortex_responded, created_at, resolved_at
		FROM sos_flares
		WHERE status = 'active'
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sosFlares []*domain.SOSFlare
	for rows.Next() {
		var sos domain.SOSFlare
		var tribeID sql.NullString
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&sos.ID,
			&sos.UserID,
			&sos.FastingID,
			&tribeID,
			&sos.Description,
			&sos.HoursFasted,
			&sos.Status,
			&sos.HypeCount,
			&sos.IsAnonymous,
			&sos.CortexResponded,
			&sos.CreatedAt,
			&resolvedAt,
		)
		if err != nil {
			return nil, err
		}

		if tribeID.Valid {
			tid, _ := uuid.Parse(tribeID.String)
			sos.TribeID = &tid
		}

		if resolvedAt.Valid {
			sos.ResolvedAt = &resolvedAt.Time
		}

		sosFlares = append(sosFlares, &sos)
	}

	return sosFlares, rows.Err()
}

// IncrementHypeCount increments the hype count for an SOS
func (r *PostgresSOSRepository) IncrementHypeCount(ctx context.Context, sosID uuid.UUID) error {
	query := `
		UPDATE sos_flares
		SET hype_count = hype_count + 1
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, sosID)
	return err
}
