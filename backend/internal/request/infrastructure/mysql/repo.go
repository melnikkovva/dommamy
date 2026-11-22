package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/request/domain"
)

const columns = "id, mother_id, kind, description, status, staff_comment, source, help_id, created_at, updated_at"

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewRepository(q DBTX) domain.Repository { return &repository{q: q} }

type repository struct{ q DBTX }

func (r *repository) Add(ctx context.Context, req domain.Request) (int, error) {
	s := req.Snapshot()
	res, err := r.q.ExecContext(ctx,
		`INSERT INTO help_request (mother_id, kind, description, status, staff_comment, source, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		s.MotherID, s.Kind, s.Description, s.Status, s.Comment, s.Source, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *repository) ByID(ctx context.Context, id int) (domain.Request, error) {
	return r.byID(ctx, id, "")
}

func (r *repository) ByIDForUpdate(ctx context.Context, id int) (domain.Request, error) {
	return r.byID(ctx, id, " FOR UPDATE")
}

func (r *repository) byID(ctx context.Context, id int, lock string) (domain.Request, error) {
	var (
		s      domain.Snapshot
		helpID sql.NullInt64
	)
	err := r.q.QueryRowContext(ctx, "SELECT "+columns+" FROM help_request WHERE id = ?"+lock, id).
		Scan(&s.ID, &s.MotherID, &s.Kind, &s.Description, &s.Status,
			&s.Comment, &s.Source, &helpID, &s.CreatedAt, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Request{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Request{}, err
	}
	s.HelpID = int(helpID.Int64)
	return domain.Restore(s), nil
}

func (r *repository) Update(ctx context.Context, req domain.Request) error {
	s := req.Snapshot()
	_, err := r.q.ExecContext(ctx,
		`UPDATE help_request
		 SET status = ?, staff_comment = ?, help_id = NULLIF(?, 0), updated_at = ?
		 WHERE id = ?`,
		s.Status, s.Comment, s.HelpID, s.UpdatedAt, s.ID,
	)
	return err
}