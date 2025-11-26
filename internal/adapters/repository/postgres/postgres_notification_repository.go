package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fastinghero/internal/core/domain"
	"time"

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
		ON CONFLICT (token) DO UPDATE 
		SET last_used_at = EXCLUDED.last_used_at
	`

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.DeviceType,
		token.CreatedAt,
		token.LastUsedAt,
	)

	return err
}

func (r *PostgresNotificationRepository) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]domain.FCMToken, error) {
	query := `
		SELECT id, user_id, token, device_type, created_at, last_used_at
		FROM fcm_tokens
		WHERE user_id = $1
		ORDER BY last_used_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.FCMToken
	for rows.Next() {
		var token domain.FCMToken
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.DeviceType,
			&token.CreatedAt,
			&token.LastUsedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}

func (r *PostgresNotificationRepository) DeleteToken(ctx context.Context, tokenString string) error {
	query := `DELETE FROM fcm_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, tokenString)
	return err
}

func (r *PostgresNotificationRepository) SaveNotification(ctx context.Context, notification *domain.Notification) error {
	dataJSON, err := json.Marshal(notification.Data)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO notifications (id, user_id, title, body, type, data, sent_at, read_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Title,
		notification.Body,
		notification.Type,
		dataJSON,
		notification.SentAt,
		notification.ReadAt,
	)

	return err
}

func (r *PostgresNotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Notification, error) {
	query := `
		SELECT id, user_id, title, body, type, data, sent_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY sent_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var notification domain.Notification
		var dataJSON []byte
		var readAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Body,
			&notification.Type,
			&dataJSON,
			&notification.SentAt,
			&readAt,
		)
		if err != nil {
			return nil, err
		}

		if readAt.Valid {
			notification.ReadAt = &readAt.Time
		}

		if len(dataJSON) > 0 {
			var data map[string]string
			if err := json.Unmarshal(dataJSON, &data); err == nil {
				notification.Data = data
			}
		}

		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (r *PostgresNotificationRepository) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	query := `UPDATE notifications SET read_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), notificationID)
	return err
}
