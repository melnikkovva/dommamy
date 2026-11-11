package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/mother"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(m *mother.Mother) error {
	res, err := r.db.Exec(
		`INSERT INTO mother (first_name, last_name, phone, created_at)
		 VALUES (?, ?, ?, ?)`,
		m.FirstName, m.LastName, m.Phone, m.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.Id = int(id)
	return nil
}

func (r *Repo) Search(q string) ([]*mother.Mother, error) {
	rows, err := r.db.Query(
		`SELECT id, first_name, last_name, phone, created_at
		 FROM mother
		 WHERE CONCAT(first_name, ' ', last_name) LIKE CONCAT('%', ?, '%')
		 ORDER BY last_name, first_name`,
		q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*mother.Mother{}
	for rows.Next() {
		m := &mother.Mother{}
		err := rows.Scan(&m.Id, &m.FirstName, &m.LastName, &m.Phone, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *Repo) ById(id int) (*mother.Mother, error) {
	m := &mother.Mother{}
	err := r.db.QueryRow(
		`SELECT id, first_name, last_name, phone, created_at
		 FROM mother WHERE id = ?`, id,
	).Scan(&m.Id, &m.FirstName, &m.LastName, &m.Phone, &m.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, mother.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}