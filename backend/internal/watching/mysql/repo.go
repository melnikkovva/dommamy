package mysql

import (
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Log(action string, motherId int) error {
	_, err := r.db.Exec(
		"INSERT INTO watching_log (action, mother_id) VALUES (?, ?)",
		action, motherId,
	)
	return err
}