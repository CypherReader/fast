package postgres

import (
	"context"
	"database/sql"
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"

	"github.com/google/uuid"
)

type PostgresTribeRepository struct {
	db *sql.DB
}

var _ ports.TribeRepository = (*PostgresTribeRepository)(nil)

func NewPostgresTribeRepository(db *sql.DB) *PostgresTribeRepository {
	return &PostgresTribeRepository{db: db}
}

func (r *PostgresTribeRepository) Save(ctx context.Context, tribe *domain.Tribe) error {
	query := `
		INSERT INTO tribes (id, name, description, leader_id, member_count, total_discipline, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			member_count = EXCLUDED.member_count,
			total_discipline = EXCLUDED.total_discipline
	`
	_, err := r.db.ExecContext(ctx, query,
		tribe.ID, tribe.Name, tribe.Description, tribe.LeaderID, tribe.MemberCount, tribe.TotalDiscipline, tribe.CreatedAt,
	)
	return err
}

func (r *PostgresTribeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Tribe, error) {
	query := `SELECT id, name, description, leader_id, member_count, total_discipline, created_at FROM tribes WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var tribe domain.Tribe
	if err := row.Scan(
		&tribe.ID, &tribe.Name, &tribe.Description, &tribe.LeaderID, &tribe.MemberCount, &tribe.TotalDiscipline, &tribe.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &tribe, nil
}

func (r *PostgresTribeRepository) FindAll(ctx context.Context) ([]domain.Tribe, error) {
	query := `SELECT id, name, description, leader_id, member_count, total_discipline, created_at FROM tribes ORDER BY total_discipline DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tribes []domain.Tribe
	for rows.Next() {
		var t domain.Tribe
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &t.LeaderID, &t.MemberCount, &t.TotalDiscipline, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tribes = append(tribes, t)
	}
	return tribes, nil
}

func (r *PostgresTribeRepository) AddMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update User's TribeID
	_, err = tx.ExecContext(ctx, `UPDATE users SET tribe_id = $1 WHERE id = $2`, tribeID, userID)
	if err != nil {
		return err
	}

	// Increment Tribe Member Count
	_, err = tx.ExecContext(ctx, `UPDATE tribes SET member_count = member_count + 1 WHERE id = $1`, tribeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresTribeRepository) RemoveMember(ctx context.Context, tribeID, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear User's TribeID
	_, err = tx.ExecContext(ctx, `UPDATE users SET tribe_id = NULL WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	// Decrement Tribe Member Count
	_, err = tx.ExecContext(ctx, `UPDATE tribes SET member_count = member_count - 1 WHERE id = $1`, tribeID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
