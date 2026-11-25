package mysql

import (
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"

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
		"INSERT INTO app_user (login, password, first_name, last_name, role, mother_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.Login, u.Password, u.FirstName, u.LastName, u.Role, u.MotherId, u.CreatedAt,
	)
	if err != nil {
		var me *mysqldriver.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return user.ErrorLoginTaken
		}
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
		"SELECT id, login, password, first_name, last_name, role, mother_id, created_at FROM app_user WHERE id = ? AND deleted_at IS NULL", id,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.MotherId, &u.CreatedAt)
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
		"SELECT id, login, password, first_name, last_name, role, mother_id, created_at FROM app_user WHERE login = ? AND deleted_at IS NULL", login,
	).Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.MotherId, &u.CreatedAt)
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
		"SELECT id, login, password, first_name, last_name, role, mother_id, created_at FROM app_user WHERE deleted_at IS NULL ORDER BY id",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*user.User{}
	for rows.Next() {
		u := &user.User{}
		if err := rows.Scan(&u.Id, &u.Login, &u.Password, &u.FirstName, &u.LastName, &u.Role, &u.MotherId, &u.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, rows.Err()
}

func (r *Repo) DeleteUser(id int) error {
	res, err := r.db.Exec(
		"UPDATE app_user SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL", id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return user.ErrorNotFound
	}
	return nil
}
