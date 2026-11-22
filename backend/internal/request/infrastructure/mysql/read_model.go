package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/request/application/query"
	"backend/internal/request/domain"
)

func NewReadModel(db *sql.DB) query.ReadModel { return &readModel{db: db} }

type readModel struct{ db *sql.DB }

func (r *readModel) Get(ctx context.Context, id int) (query.RequestView, error) {
	var v query.RequestView
	var helpID sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		"SELECT id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at FROM help_request WHERE id = ?", id).
		Scan(&v.ID, &v.MotherID, &v.Kind, &v.Description, &v.Status,
			&v.Comment, &v.Source, &helpID, &v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return query.RequestView{}, domain.ErrNotFound
	}
	if err != nil {
		return query.RequestView{}, err
	}
	v.HelpID = int(helpID.Int64)
	return v, nil
}

func (r *readModel) List(ctx context.Context, status string) ([]query.RequestView, error) {
	var queryStr string
	var args []any
	if status != "" {
		queryStr = "SELECT id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at FROM help_request WHERE status = ? ORDER BY created_at DESC"
		args = []any{status}
	} else {
		queryStr = "SELECT id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at FROM help_request ORDER BY created_at DESC"
	}
	rows, err := r.db.QueryContext(ctx, queryStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.RequestView
	for rows.Next() {
		var v query.RequestView
		var helpID sql.NullInt64
		if err := rows.Scan(&v.ID, &v.MotherID, &v.Kind, &v.Description, &v.Status,
			&v.Comment, &v.Source, &helpID, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		v.HelpID = int(helpID.Int64)
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *readModel) ByMother(ctx context.Context, motherID int) ([]query.RequestView, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at FROM help_request WHERE mother_id = ? ORDER BY created_at DESC", motherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.RequestView
	for rows.Next() {
		var v query.RequestView
		var helpID sql.NullInt64
		if err := rows.Scan(&v.ID, &v.MotherID, &v.Kind, &v.Description, &v.Status,
			&v.Comment, &v.Source, &helpID, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		v.HelpID = int(helpID.Int64)
		result = append(result, v)
	}
	return result, rows.Err()
}
