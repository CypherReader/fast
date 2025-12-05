package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"

	"github.com/google/uuid"
)

type PostgresNotificationRepository struct {
	db *sql.DB
}

func NewPostgresNotificationRepository(db *sql.DB) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{db: db}
}

func (r *PostgresNotificationRepository) SaveToken(ctx context.Context, token *domain.FCMToken) error {
	query := `
		INSERT INTO fcm_tokens (id, user_id, token, device_type, created_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (token) DO UPDATE SET
			last_used_at = EXCLUDED.last_used_at,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.Token, token.DeviceType, token.CreatedAt, token.LastUsedAt)
	return err
}

func (r *PostgresNotificationRepository) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error) {
	query := `SELECT id, user_id, token, device_type, created_at FROM fcm_tokens WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.FCMToken
	for rows.Next() {
		var t domain.FCMToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.DeviceType, &t.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func (r *PostgresNotificationRepository) DeleteToken(ctx context.Context, tokenString string) error {
	query := `DELETE FROM fcm_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, tokenString)
	return err
}

func (r *PostgresNotificationRepository) Save(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, read, link, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, n.ID, n.UserID, n.Type, n.Title, n.Message, n.Read, n.Link, n.CreatedAt)
	return err
}

func (r *PostgresNotificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, read, link, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		var link sql.NullString
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Message, &n.Read, &link, &n.CreatedAt); err != nil {
			return nil, err
		}
		if link.Valid {
			n.Link = link.String
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *PostgresNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET read = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresNotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET read = true WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
