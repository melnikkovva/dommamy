package mysql

import (
	"database/sql"

	"backend/internal/child"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(c *child.Child) error {
	res, err := r.db.Exec(
		"INSERT INTO child (mother_id, first_name, birth_date, created_at) VALUES (?, ?, ?, ?)",
		c.MotherId, c.FirstName, c.BirthDate, c.CreatedAt,
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

func (r *Repo) SearchByMother(motherId int) ([]*child.Child, error) {
	rows, err := r.db.Query(
		"SELECT id, mother_id, first_name, birth_date, created_at FROM child WHERE mother_id = ? ORDER BY birth_date",
		motherId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*child.Child{}
	for rows.Next() {
		c := &child.Child{}
		err := rows.Scan(&c.Id, &c.MotherId, &c.FirstName, &c.BirthDate, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}