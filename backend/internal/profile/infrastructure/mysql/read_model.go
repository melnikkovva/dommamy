package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/profile/application/query"
	"backend/internal/profile/domain"
)

func NewReadModel(db *sql.DB) query.ReadModel {
	return &readModel{db: db}
}

type readModel struct {
	db *sql.DB
}

func (r *readModel) GetByMother(ctx context.Context, motherID int) (query.ProfileView, error) {
	var v query.ProfileView
	var statusStr string
	var updatedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`SELECT mother_id, COALESCE(DATE_FORMAT(birth_date, '%Y-%m-%d'), ''), birth_place,
		        citizenship, registration, living, email, religion, education, profession,
		        income_source, income_amount, marital_status, father_info, pregnancy_weeks,
		        situation, expected_help, source, status, updated_at
		 FROM mother_profile WHERE mother_id = ?`,
		motherID,
	).Scan(
		&v.MotherID, &v.BirthDate, &v.BirthPlace, &v.Citizenship, &v.Registration,
		&v.Living, &v.Email, &v.Religion, &v.Education, &v.Profession, &v.IncomeSource,
		&v.IncomeAmount, &v.MaritalStatus, &v.FatherInfo, &v.PregnancyWeeks,
		&v.Situation, &v.ExpectedHelp, &v.Source, &statusStr, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return query.ProfileView{}, domain.ErrNotFilled
	}
	if err != nil {
		return query.ProfileView{}, err
	}

	v.Status = domain.Status(statusStr)
	if updatedAt.Valid {
		v.UpdatedAt = updatedAt.Time
	}

	return v, nil
}

func (r *readModel) ListByStatus(ctx context.Context, status string) ([]query.ProfileListItem, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT mother_id, status, updated_at FROM mother_profile
		 WHERE ? = '' OR status = ? ORDER BY updated_at DESC`,
		status, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []query.ProfileListItem
	for rows.Next() {
		var item query.ProfileListItem
		var statusStr string
		var updatedAt sql.NullTime
		if err := rows.Scan(&item.MotherID, &statusStr, &updatedAt); err != nil {
			return nil, err
		}
		item.Status = domain.Status(statusStr)
		if updatedAt.Valid {
			item.UpdatedAt = updatedAt.Time
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
