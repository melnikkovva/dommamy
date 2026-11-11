package mysql

import (
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"

	"backend/internal/document"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

const columns = `id, mother_id, doc_type, help_id, status,
	COALESCE(DATE_FORMAT(signed_at, '%Y-%m-%d'), ''),
	file_name, original_name, content_type, file_size, created_at`

type scanner interface {
	Scan(dest ...any) error
}

func scan(s scanner) (*document.Document, error) {
	d := &document.Document{}
	var helpId sql.NullInt64
	err := s.Scan(&d.Id, &d.MotherId, &d.Type, &helpId, &d.Status, &d.SignedAt,
		&d.FileName, &d.OriginalName, &d.ContentType, &d.Size, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	d.HelpId = int(helpId.Int64)
	return d, nil
}

func (r *Repo) Save(d *document.Document) error {
	res, err := r.db.Exec(
		`INSERT INTO document (mother_id, doc_type, help_id, status, created_at)
		 VALUES (?, ?, NULLIF(?, 0), ?, ?)`,
		d.MotherId, d.Type, d.HelpId, d.Status, d.CreatedAt,
	)
	if err != nil {
		var me *mysqldriver.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return document.ErrorDuplicate
		}
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = int(id)
	return nil
}

func (r *Repo) ById(id int) (*document.Document, error) {
	d, err := scan(r.db.QueryRow("SELECT "+columns+" FROM document WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, document.ErrorNotFound
	}
	return d, err
}

func (r *Repo) ByMother(motherId int) ([]*document.Document, error) {
	rows, err := r.db.Query("SELECT "+columns+" FROM document WHERE mother_id = ? ORDER BY id DESC", motherId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*document.Document{}
	for rows.Next() {
		d, err := scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *Repo) Update(d *document.Document) error {
	_, err := r.db.Exec(
		`UPDATE document
		 SET status = ?, signed_at = NULLIF(?, ''), file_name = ?, original_name = ?,
		     content_type = ?, file_size = ?
		 WHERE id = ?`,
		d.Status, d.SignedAt, d.FileName, d.OriginalName, d.ContentType, d.Size, d.Id,
	)
	return err
}