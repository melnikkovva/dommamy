package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/mother/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, m domain.Mother) error {
	var learnUserID sql.NullString
	if m.LearnUserID() != "" {
		learnUserID = sql.NullString{String: m.LearnUserID(), Valid: true}
	}
	var photoURL sql.NullString
	if m.PhotoURL() != "" {
		photoURL = sql.NullString{String: m.PhotoURL(), Valid: true}
	}
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO mother (first_name, last_name, phone, email, learn_user_id, photo_url, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		m.FirstName(), m.LastName(), m.Phone(), m.Email(), learnUserID, photoURL, m.CreatedAt(),
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.SetID(int(id))
	return nil
}

func (r *Repository) ByID(ctx context.Context, id int) (domain.MotherSnapshot, error) {
	var s domain.MotherSnapshot
	var learnUserID sql.NullString
	var photoURL sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, first_name, last_name, phone, email, learn_user_id, photo_url, created_at
		 FROM mother WHERE id = ?`, id,
	).Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &learnUserID, &photoURL, &s.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.MotherSnapshot{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.MotherSnapshot{}, err
	}
	if learnUserID.Valid {
		s.LearnUserID = learnUserID.String
	}
	if photoURL.Valid {
		s.PhotoURL = photoURL.String
	}
	return s, nil
}

func (r *Repository) Search(ctx context.Context, q string) ([]domain.MotherSnapshot, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, first_name, last_name, phone, email, learn_user_id, photo_url, created_at
		 FROM mother
		 WHERE CONCAT(first_name, ' ', last_name) LIKE CONCAT('%', ?, '%')
		 ORDER BY last_name, first_name`,
		q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.MotherSnapshot
	for rows.Next() {
		var s domain.MotherSnapshot
		var learnUserID sql.NullString
		var photoURL sql.NullString
		if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &learnUserID, &photoURL, &s.CreatedAt); err != nil {
			return nil, err
		}
		if learnUserID.Valid {
			s.LearnUserID = learnUserID.String
		}
		if photoURL.Valid {
			s.PhotoURL = photoURL.String
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repository) GetAll(ctx context.Context) ([]domain.MotherSnapshot, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT m.id, m.first_name, m.last_name, m.phone, mp.email, m.learn_user_id, m.photo_url, m.created_at
		 FROM mother m
		 LEFT JOIN mother_profile mp ON mp.mother_id = m.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.MotherSnapshot
	for rows.Next() {
		var s domain.MotherSnapshot
		var email sql.NullString
		var learnUserID sql.NullString
		var photoURL sql.NullString
		if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &email, &learnUserID, &photoURL, &s.CreatedAt); err != nil {
			return nil, err
		}
		if email.Valid {
			s.Email = email.String
		}
		if learnUserID.Valid {
			s.LearnUserID = learnUserID.String
		}
		if photoURL.Valid {
			s.PhotoURL = photoURL.String
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repository) UpdateLearnUserID(ctx context.Context, id int, learnUserID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE mother SET learn_user_id = ? WHERE id = ?`, learnUserID, id)
	return err
}

func (r *Repository) UpdatePhotoURL(ctx context.Context, id int, photoURL string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE mother SET photo_url = ? WHERE id = ?`, photoURL, id)
	return err
}
