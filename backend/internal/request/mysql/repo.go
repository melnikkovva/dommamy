package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/request"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

const columns = "id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at"

type scanner interface {
	Scan(dest ...any) error
}

func scan(s scanner) (*request.Request, error) {
	r := &request.Request{}
	var helpId sql.NullInt64
	err := s.Scan(&r.Id, &r.MotherId, &r.Kind, &r.Description, &r.Status,
		&r.Comment, &r.Source, &helpId, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}
	r.HelpId = int(helpId.Int64) 
	return r, nil
}

func scanAll(rows *sql.Rows) ([]*request.Request, error) {
	defer rows.Close()
	result := []*request.Request{}
	for rows.Next() {
		r, err := scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, rows.Err()
}

func (r *Repo) Save(req *request.Request) error {
	res, err := r.db.Exec(
		`INSERT INTO help_request (mother_id, kind, description, status, staff_comment, source, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		req.MotherId, req.Kind, req.Description, req.Status, req.Comment, req.Source, req.CreatedAt, req.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	req.Id = int(id)
	return nil
}

func (r *Repo) ById(id int) (*request.Request, error) {
	req, err := scan(r.db.QueryRow("SELECT "+columns+" FROM help_request WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, request.ErrorNotFound
	}
	return req, err
}

func (r *Repo) List(status string) ([]*request.Request, error) {
	rows, err := r.db.Query(
		"SELECT "+columns+" FROM help_request WHERE ? = '' OR status = ? ORDER BY id DESC",
		status, status,
	)
	if err != nil {
		return nil, err
	}
	return scanAll(rows)
}

func (r *Repo) ByMother(motherId int) ([]*request.Request, error) {
	rows, err := r.db.Query(
		"SELECT "+columns+" FROM help_request WHERE mother_id = ? ORDER BY id DESC", motherId,
	)
	if err != nil {
		return nil, err
	}
	return scanAll(rows)
}

func (r *Repo) Update(req *request.Request) error {
	_, err := r.db.Exec(
		`UPDATE help_request
		 SET status = ?, staff_comment = ?, help_id = NULLIF(?, 0), updated_at = ?
		 WHERE id = ?`,
		req.Status, req.Comment, req.HelpId, req.UpdatedAt, req.Id,
	)
	return err
}