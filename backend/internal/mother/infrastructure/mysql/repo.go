package mysql

import (
	"context"
	"database/sql"

	"backend/internal/mother/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, m domain.Mother) (int, error) {
	var photoURL sql.NullString
	if m.PhotoURL() != "" {
		photoURL = sql.NullString{String: m.PhotoURL(), Valid: true}
	}
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO mother (first_name, last_name, phone, email, photo_url, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		m.FirstName(), m.LastName(), m.Phone(), m.Email(), photoURL, m.CreatedAt(),
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

func (r *Repository) UpdatePhotoURL(ctx context.Context, id int, photoURL string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE mother SET photo_url = ? WHERE id = ?`, photoURL, id)
	return err
}
