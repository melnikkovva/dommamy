package mysql

import (
	"context"
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"

	"backend/internal/document/domain"
)

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}

const columns = `id, mother_id, doc_type, help_id, status,
	COALESCE(DATE_FORMAT(signed_at, '%Y-%m-%d'), ''),
	file_name, original_name, content_type, file_size, created_at`

func (r *repository) Add(ctx context.Context, d domain.Document) (int, error) {
	s := d.Snapshot()
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO document (mother_id, doc_type, help_id, status, created_at)
		 VALUES (?, ?, NULLIF(?, 0), ?, ?)`,
		s.MotherID, string(s.DocType), s.HelpID, string(s.Status), s.CreatedAt,
	)
	if err != nil {
		var me *mysqldriver.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return 0, domain.ErrDuplicate
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *repository) ByID(ctx context.Context, id int) (domain.Document, error) {
	var s domain.Snapshot
	var helpID sql.NullInt64
	var signedAt sql.NullString
	var fileName sql.NullString
	var originalName sql.NullString
	var contentType sql.NullString
	var fileSize sql.NullInt64

	row := r.db.QueryRowContext(
		ctx,
		"SELECT "+columns+" FROM document WHERE id = ?", id,
	)
	err := row.Scan(
		&s.ID, &s.MotherID, &s.DocType, &helpID, &s.Status, &signedAt,
		&fileName, &originalName, &contentType, &fileSize, &s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Document{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Document{}, err
	}

	if helpID.Valid {
		s.HelpID = int(helpID.Int64)
	}
	if signedAt.Valid {
		s.SignedAt = signedAt.String
	}
	if fileName.Valid {
		s.FileName = fileName.String
	}
	if originalName.Valid {
		s.OriginalName = originalName.String
	}
	if contentType.Valid {
		s.ContentType = contentType.String
	}
	if fileSize.Valid {
		s.Size = fileSize.Int64
	}

	return domain.Restore(s), nil
}

func (r *repository) ByMother(ctx context.Context, motherID int) ([]domain.Document, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+columns+" FROM document WHERE mother_id = ? ORDER BY id DESC", motherID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []domain.Document{}
	for rows.Next() {
		d, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, rows.Err()
}

func (r *repository) Update(ctx context.Context, d domain.Document) error {
	s := d.Snapshot()
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE document
		 SET status = ?, signed_at = NULLIF(?, ''), file_name = ?, original_name = ?,
		     content_type = ?, file_size = ?
		 WHERE id = ?`,
		string(s.Status), s.SignedAt, s.FileName, s.OriginalName, s.ContentType, s.Size, s.ID,
	)
	return err
}

func scanRow(s rowScanner) (domain.Document, error) {
	var sDoc domain.Snapshot
	var helpID sql.NullInt64
	var signedAt sql.NullString
	var fileName sql.NullString
	var originalName sql.NullString
	var contentType sql.NullString
	var fileSize sql.NullInt64

	err := s.Scan(
		&sDoc.ID, &sDoc.MotherID, &sDoc.DocType, &helpID, &sDoc.Status, &signedAt,
		&fileName, &originalName, &contentType, &fileSize, &sDoc.CreatedAt,
	)
	if err != nil {
		return domain.Document{}, err
	}

	if helpID.Valid {
		sDoc.HelpID = int(helpID.Int64)
	}
	if signedAt.Valid {
		sDoc.SignedAt = signedAt.String
	}
	if fileName.Valid {
		sDoc.FileName = fileName.String
	}
	if originalName.Valid {
		sDoc.OriginalName = originalName.String
	}
	if contentType.Valid {
		sDoc.ContentType = contentType.String
	}
	if fileSize.Valid {
		sDoc.Size = fileSize.Int64
	}

	return domain.Restore(sDoc), nil
}

type rowScanner interface {
	Scan(dest ...any) error
}
