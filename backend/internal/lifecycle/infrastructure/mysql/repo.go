package mysql

import (
	"context"
	"database/sql"

	"backend/internal/lifecycle/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, c domain.Change) (int, error) {
	res, err := r.db.ExecContext(
		ctx,
		"INSERT INTO mother_status_history (mother_id, status, reason, created_at) VALUES (?, ?, ?, ?)",
		c.MotherID(), c.Status(), c.Reason(), c.CreatedAt(),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *Repository) History(ctx context.Context, motherID int) ([]domain.ChangeSnapshot, error) {
	rows, err := r.db.QueryContext(
		ctx,
 	"SELECT id, mother_id, status, reason, created_at FROM mother_status_history WHERE mother_id = ? ORDER BY id ASC",
		motherID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ChangeSnapshot
	for rows.Next() {
		var s domain.ChangeSnapshot
		if err := rows.Scan(&s.ID, &s.MotherID, &s.Status, &s.Reason, &s.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repository) Currents(ctx context.Context) ([]domain.ChangeSnapshot, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT h.id, h.mother_id, h.status, h.reason, h.created_at
		 FROM mother_status_history h
		 JOIN (SELECT MAX(id) AS last_id FROM mother_status_history GROUP BY mother_id) m
		   ON h.id = m.last_id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ChangeSnapshot
	for rows.Next() {
		var s domain.ChangeSnapshot
		if err := rows.Scan(&s.ID, &s.MotherID, &s.Status, &s.Reason, &s.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
