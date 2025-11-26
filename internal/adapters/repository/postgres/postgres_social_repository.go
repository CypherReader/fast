package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type PostgresSocialRepository struct {
	db *sql.DB
}

var _ ports.SocialRepository = (*PostgresSocialRepository)(nil)

func NewPostgresSocialRepository(db *sql.DB) *PostgresSocialRepository {
	return &PostgresSocialRepository{db: db}
}

func (r *PostgresSocialRepository) SaveEvent(ctx context.Context, event *domain.SocialEvent) error {
	query := `
		INSERT INTO social_events (id, user_id, user_name, type, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.UserID, event.UserName, event.Type, event.Data, event.CreatedAt,
	)
	return err
}

func (r *PostgresSocialRepository) GetFeed(ctx context.Context, limit int) ([]domain.SocialEvent, error) {
	query := `
		SELECT id, user_id, user_name, type, data, created_at
		FROM social_events
		ORDER BY created_at DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.SocialEvent
	for rows.Next() {
		var e domain.SocialEvent
		var data []byte
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.UserName, &e.Type, &data, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		e.Data = json.RawMessage(data)
		events = append(events, e)
	}
	return events, nil
}

func (r *PostgresSocialRepository) GetTribeFeed(ctx context.Context, tribeID uuid.UUID, limit int) ([]domain.SocialEvent, error) {
	// Join with users table to filter by tribe_id
	query := `
		SELECT e.id, e.user_id, e.user_name, e.type, e.data, e.created_at
		FROM social_events e
		JOIN users u ON e.user_id = u.id
		WHERE u.tribe_id = $1
		ORDER BY e.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, tribeID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.SocialEvent
	for rows.Next() {
		var e domain.SocialEvent
		var data []byte
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.UserName, &e.Type, &data, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		e.Data = json.RawMessage(data)
		events = append(events, e)
	}
	return events, nil
}
