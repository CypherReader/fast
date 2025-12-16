package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"time"

	"github.com/google/uuid"
)

type PostgresSocialRepository struct {
	db *sql.DB
}

func NewPostgresSocialRepository(db *sql.DB) *PostgresSocialRepository {
	return &PostgresSocialRepository{db: db}
}

// Friend Network
func (r *PostgresSocialRepository) SaveFriendNetwork(ctx context.Context, fn *domain.FriendNetwork) error {
	query := `
		INSERT INTO friend_networks (id, user_id, friend_id, status, connected_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, friend_id) DO UPDATE SET
			status = EXCLUDED.status,
			connected_at = EXCLUDED.connected_at
	`
	_, err := r.db.ExecContext(ctx, query, fn.ID, fn.UserID, fn.FriendID, fn.Status, fn.ConnectedAt, fn.CreatedAt)
	return err
}

func (r *PostgresSocialRepository) FindFriends(ctx context.Context, userID uuid.UUID) ([]domain.FriendNetwork, error) {
	query := `SELECT id, user_id, friend_id, status, connected_at, created_at FROM friend_networks WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []domain.FriendNetwork
	for rows.Next() {
		var fn domain.FriendNetwork
		var connectedAt sql.NullTime
		if err := rows.Scan(&fn.ID, &fn.UserID, &fn.FriendID, &fn.Status, &connectedAt, &fn.CreatedAt); err != nil {
			return nil, err
		}
		if connectedAt.Valid {
			fn.ConnectedAt = connectedAt.Time
		}
		friends = append(friends, fn)
	}
	return friends, nil
}

// Tribes
func (r *PostgresSocialRepository) SaveTribe(ctx context.Context, tribe *domain.Tribe) error {
	query := `
		INSERT INTO tribes (id, name, description, creator_id, member_count, privacy, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			member_count = EXCLUDED.member_count,
			privacy = EXCLUDED.privacy,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		tribe.ID, tribe.Name, tribe.Description, tribe.CreatorID, tribe.MemberCount, tribe.Privacy,
		tribe.CreatedAt, time.Now(),
	)
	return err
}

func (r *PostgresSocialRepository) FindTribeByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error) {
	query := `SELECT id, name, description, creator_id, member_count, privacy, created_at, updated_at FROM tribes WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id.String())

	var t domain.Tribe
	err := row.Scan(&t.ID, &t.Name, &t.Description, &t.CreatorID, &t.MemberCount, &t.Privacy, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresSocialRepository) FindAllTribes(ctx context.Context, limit, offset int) ([]domain.Tribe, error) {
	query := `
		SELECT id, name, description, creator_id, member_count, privacy, created_at, updated_at
		FROM tribes
		ORDER BY member_count DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tribes []domain.Tribe
	for rows.Next() {
		var t domain.Tribe
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatorID, &t.MemberCount, &t.Privacy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tribes = append(tribes, t)
	}
	return tribes, nil
}

// Challenges
func (r *PostgresSocialRepository) SaveChallenge(ctx context.Context, c *domain.FriendChallenge) error {
	query := `
		INSERT INTO friend_challenges (
			id, creator_id, name, challenge_type, goal, start_date, end_date, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			goal = EXCLUDED.goal,
			end_date = EXCLUDED.end_date,
			status = EXCLUDED.status
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.CreatorID, c.Name, c.ChallengeType, c.Goal, c.StartDate, c.EndDate, c.Status, c.CreatedAt,
	)
	return err
}

func (r *PostgresSocialRepository) FindChallengesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.FriendChallenge, error) {
	query := `
		SELECT id, creator_id, name, challenge_type, goal, start_date, end_date, status, created_at
		FROM friend_challenges WHERE creator_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []domain.FriendChallenge
	for rows.Next() {
		var c domain.FriendChallenge
		if err := rows.Scan(
			&c.ID, &c.CreatorID, &c.Name, &c.ChallengeType, &c.Goal, &c.StartDate, &c.EndDate, &c.Status, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}
	return challenges, nil
}

// Feed
func (r *PostgresSocialRepository) SaveEvent(ctx context.Context, event *domain.SocialEvent) error {
	query := `
		INSERT INTO social_events (id, user_id, event_type, data, created_at, likes, comments)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.UserID, event.EventType, event.Data, event.CreatedAt, event.Likes, event.Comments,
	)
	return err
}

func (r *PostgresSocialRepository) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.SocialEvent, error) {
	// For now, global feed. In future, filter by friends/tribes.
	query := `
		SELECT e.id, e.user_id, u.name, e.event_type, e.data, e.created_at, e.likes, e.comments
		FROM social_events e
		JOIN users u ON e.user_id = u.id
		ORDER BY e.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.SocialEvent
	for rows.Next() {
		var e domain.SocialEvent
		var userName sql.NullString
		if err := rows.Scan(
			&e.ID, &e.UserID, &userName, &e.EventType, &e.Data, &e.CreatedAt, &e.Likes, &e.Comments,
		); err != nil {
			return nil, err
		}
		if userName.Valid {
			e.UserName = userName.String
		} else {
			e.UserName = "Anonymous"
		}
		events = append(events, e)
	}
	return events, nil
}
