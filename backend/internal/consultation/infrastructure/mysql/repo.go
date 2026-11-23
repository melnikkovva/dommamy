package mysql

import (
	"context"
	"database/sql"
	"time"

	"backend/internal/consultation/domain"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(ctx context.Context, c *domain.Consultation) error {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO consultation (mother_id, specialist_type, consultation_date, notes, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		c.MotherId, string(c.SpecialistType), c.ConsultationDate.Format("2006-01-02"),
		c.Notes, string(c.Status), c.CreatedAt,
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

func (r *Repo) ByMother(ctx context.Context, motherId int) ([]domain.Row, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT c.id, c.mother_id, CONCAT(m.last_name, ' ', m.first_name), c.specialist_type, "+
			"c.consultation_date, c.notes, c.status, c.created_at "+
			"FROM consultation c JOIN mother m ON m.id = c.mother_id "+
			"WHERE c.mother_id = ? ORDER BY c.consultation_date DESC, c.id DESC",
		motherId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *Repo) List(ctx context.Context) ([]domain.Row, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT c.id, c.mother_id, CONCAT(m.last_name, ' ', m.first_name), c.specialist_type, "+
			"c.consultation_date, c.notes, c.status, c.created_at "+
			"FROM consultation c JOIN mother m ON m.id = c.mother_id "+
			"ORDER BY c.consultation_date DESC, c.id DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *Repo) UpdateStatus(ctx context.Context, id int, status domain.Status) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE consultation SET status = ? WHERE id = ?",
		string(status), id,
	)
	return err
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM consultation WHERE id = ?",
		id,
	)
	return err
}

func (r *Repo) CountByMotherAndType(ctx context.Context, motherId int, specialistType domain.SpecialistType) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM consultation WHERE mother_id = ? AND specialist_type = ?",
		motherId, string(specialistType),
	).Scan(&count)
	return count, err
}

func scanRows(rows *sql.Rows) ([]domain.Row, error) {
	result := []domain.Row{}
	for rows.Next() {
		var r domain.Row
		var consultationDateStr string
		var createdAtStr string
		err := rows.Scan(&r.Id, &r.MotherId, &r.MotherName, &r.SpecialistType,
			&consultationDateStr, &r.Notes, &r.Status, &createdAtStr)
		if err != nil {
			return nil, err
		}
		r.ConsultationDate, _ = time.Parse("2006-01-02", consultationDateStr)
		r.CreatedAt, _ = time.Parse(time.DateTime, createdAtStr)
		result = append(result, r)
	}
	return result, rows.Err()
}
