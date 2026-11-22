package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/lifecycle"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Add(c *lifecycle.Change) error {
	res, err := r.db.Exec(
		"INSERT INTO mother_status_history (mother_id, status, reason, created_at) VALUES (?, ?, ?, ?)",
		c.MotherId, c.Status, c.Reason, c.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.Id = int(id)
	return nil
}

func (r *Repo) Current(motherId int) (*lifecycle.Change, error) {
	c := &lifecycle.Change{}
	err := r.db.QueryRow(
		`SELECT id, mother_id, status, reason, created_at FROM mother_status_history
		 WHERE mother_id = ? ORDER BY id DESC LIMIT 1`, motherId,
	).Scan(&c.Id, &c.MotherId, &c.Status, &c.Reason, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, lifecycle.ErrorNoMother
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func scanChanges(rows *sql.Rows) ([]*lifecycle.Change, error) {
	defer rows.Close()
	result := []*lifecycle.Change{}
	for rows.Next() {
		c := &lifecycle.Change{}
		if err := rows.Scan(&c.Id, &c.MotherId, &c.Status, &c.Reason, &c.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *Repo) History(motherId int) ([]*lifecycle.Change, error) {
	rows, err := r.db.Query(
		`SELECT id, mother_id, status, reason, created_at FROM mother_status_history
		 WHERE mother_id = ? ORDER BY id DESC`, motherId)
	if err != nil {
		return nil, err
	}
	return scanChanges(rows)
}

func (r *Repo) Currents() ([]*lifecycle.Change, error) {
	rows, err := r.db.Query(
		`SELECT h.id, h.mother_id, h.status, h.reason, h.created_at
		 FROM mother_status_history h
		 JOIN (SELECT MAX(id) AS last_id FROM mother_status_history GROUP BY mother_id) m
		   ON h.id = m.last_id`)
	if err != nil {
		return nil, err
	}
	return scanChanges(rows)
}