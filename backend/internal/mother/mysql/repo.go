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
	var learnUserID interface{}
	if m.LearnUserID != "" {
		learnUserID = m.LearnUserID
	}
	res, err := r.db.Exec(
		`INSERT INTO mother (first_name, last_name, phone, email, learn_user_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		m.FirstName, m.LastName, m.Phone, m.Email, learnUserID, m.CreatedAt,
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
		`SELECT id, first_name, last_name, phone, email, learn_user_id, created_at
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
		var learnUserID sql.NullString
		err := rows.Scan(&m.Id, &m.FirstName, &m.LastName, &m.Phone, &m.Email, &learnUserID, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		if learnUserID.Valid {
			m.LearnUserID = learnUserID.String
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *Repo) ById(id int) (*mother.Mother, error) {
	m := &mother.Mother{}
	var learnUserID sql.NullString
	err := r.db.QueryRow(
		`SELECT id, first_name, last_name, phone, email, learn_user_id, created_at
		 FROM mother WHERE id = ?`, id,
	).Scan(&m.Id, &m.FirstName, &m.LastName, &m.Phone, &m.Email, &learnUserID, &m.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, mother.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	if learnUserID.Valid {
		m.LearnUserID = learnUserID.String
	}
	return m, nil
}

func (r *Repo) GetAll() ([]*mother.Mother, error) {
	rows, err := r.db.Query(
		`SELECT m.id, m.first_name, m.last_name, m.phone, mp.email, m.learn_user_id, m.created_at
		 FROM mother m
		 LEFT JOIN mother_profile mp ON mp.mother_id = m.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*mother.Mother{}
	for rows.Next() {
		m := &mother.Mother{}
		var email sql.NullString
		var learnUserID sql.NullString
		err := rows.Scan(&m.Id, &m.FirstName, &m.LastName, &m.Phone, &email, &learnUserID, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		if email.Valid {
			m.Email = email.String
		}
		if learnUserID.Valid {
			m.LearnUserID = learnUserID.String
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *Repo) UpdateLearnUserID(id int, learnUserID string) error {
	_, err := r.db.Exec(`UPDATE mother SET learn_user_id = ? WHERE id = ?`, learnUserID, id)
	return err
}