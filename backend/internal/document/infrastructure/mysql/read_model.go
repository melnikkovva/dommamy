package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/document/application/query"
	"backend/internal/document/domain"
)

func NewReadModel(db *sql.DB) query.ReadModel {
	return &readModel{db: db}
}

type readModel struct {
	db *sql.DB
}

const viewColumns = `id, mother_id, doc_type, help_id, status,
	COALESCE(DATE_FORMAT(signed_at, '%Y-%m-%d'), ''),
	file_name, original_name, content_type, file_size, created_at`

func (r *readModel) GetByID(ctx context.Context, id int) (query.DocumentView, error) {
	v, err := r.scanView(ctx, r.db.QueryRowContext(ctx, "SELECT "+viewColumns+" FROM document WHERE id = ?", id))
	if err != nil {
		return query.DocumentView{}, err
	}
	return v, nil
}

func (r *readModel) ByMother(ctx context.Context, motherID int) ([]query.DocumentView, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+viewColumns+" FROM document WHERE mother_id = ? ORDER BY id DESC", motherID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []query.DocumentView{}
	for rows.Next() {
		v, err := r.scanView(ctx, rows)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *readModel) scanView(ctx context.Context, s viewScanner) (query.DocumentView, error) {
	var v query.DocumentView
	var helpID sql.NullInt64
	var fileName sql.NullString
	var originalName sql.NullString
	var contentType sql.NullString
	var fileSize sql.NullInt64

	err := s.Scan(
		&v.ID, &v.MotherID, &v.Type, &helpID, &v.Status, &v.SignedAt,
		&fileName, &originalName, &contentType, &fileSize, &v.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return query.DocumentView{}, domain.ErrNotFound
	}
	if err != nil {
		return query.DocumentView{}, err
	}

	if helpID.Valid {
		v.HelpID = int(helpID.Int64)
	}
	if fileName.Valid {
		v.HasFile = fileName.String != ""
		v.OriginalName = originalName.String
		v.ContentType = contentType.String
		if fileSize.Valid {
			v.Size = fileSize.Int64
		}
	}

	return v, nil
}

type viewScanner interface {
	Scan(dest ...any) error
}
