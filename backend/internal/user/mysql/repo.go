package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/user"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(u *user.User) error {
	res, err := r.db.Exec(
		"INSERT INTO app_user (login, password, first_name, last_name, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		u.Login, u.Password, u.FirstName, u.LastName, u.Role, u.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = int(id)
	return nil
}

func (r *Repo) UserById(id int) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRow(
		"SELECT id, login, password, first_name, last_name, role, created_at FROM app_user WHERE id = ?", id,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repo) UserByLogin(login string) (*user.User, error) {
	u := &user.User{}
	err := r.db.QueryRow(
		"SELECT id, login, password, first_name, last_name, role, created_at FROM app_user WHERE login = ?", login,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repo) ListUsers() ([]*user.User, error) {
	rows, err := r.db.Query(
		"SELECT id, login, password, first_name, last_name, role, created_at FROM app_user ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*user.User{}
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}
