package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"backend/internal/profile/domain"
)

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(ctx context.Context, p domain.Profile) error {
	s := p.Snapshot()
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO mother_profile
		 (mother_id, birth_date, birth_place, citizenship, registration, living, email,
		  religion, education, profession, income_source, income_amount, marital_status,
		  father_info, pregnancy_weeks, situation, expected_help, source, status, updated_at)
		 VALUES (?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		 birth_date = VALUES(birth_date), birth_place = VALUES(birth_place),
		 citizenship = VALUES(citizenship), registration = VALUES(registration),
		 living = VALUES(living), email = VALUES(email), religion = VALUES(religion),
		 education = VALUES(education), profession = VALUES(profession),
		 income_source = VALUES(income_source), income_amount = VALUES(income_amount),
		 marital_status = VALUES(marital_status), father_info = VALUES(father_info),
		 pregnancy_weeks = VALUES(pregnancy_weeks), situation = VALUES(situation),
		 expected_help = VALUES(expected_help), source = VALUES(source),
		 status = VALUES(status), updated_at = VALUES(updated_at)`,
		s.MotherID, s.BirthDate, s.BirthPlace, s.Citizenship, s.Registration, s.Living, s.Email,
		s.Religion, s.Education, s.Profession, s.IncomeSource, s.IncomeAmount, s.MaritalStatus,
		s.FatherInfo, s.PregnancyWeeks, s.Situation, s.ExpectedHelp, s.Source, string(s.Status), s.UpdatedAt,
	)
	return err
}

func (r *repository) ByMother(ctx context.Context, motherID int) (domain.Profile, error) {
	var s domain.Snapshot
	row := r.db.QueryRowContext(
		ctx,
		`SELECT mother_id, COALESCE(DATE_FORMAT(birth_date, '%Y-%m-%d'), ''), birth_place,
		        citizenship, registration, living, email, religion, education, profession,
		        income_source, income_amount, marital_status, father_info, pregnancy_weeks,
		        situation, expected_help, source, status, updated_at
		 FROM mother_profile WHERE mother_id = ?`,
		motherID,
	)

	var statusStr string
	var updatedAt sql.NullTime
	err := row.Scan(
		&s.MotherID, &s.BirthDate, &s.BirthPlace, &s.Citizenship, &s.Registration,
		&s.Living, &s.Email, &s.Religion, &s.Education, &s.Profession, &s.IncomeSource,
		&s.IncomeAmount, &s.MaritalStatus, &s.FatherInfo, &s.PregnancyWeeks,
		&s.Situation, &s.ExpectedHelp, &s.Source, &statusStr, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Profile{}, domain.ErrNotFilled
	}
	if err != nil {
		return domain.Profile{}, err
	}

	s.Status = domain.Status(statusStr)
	if updatedAt.Valid {
		s.UpdatedAt = updatedAt.Time
	}

	return domain.Restore(s), nil
}

func (r *repository) UpdateStatus(ctx context.Context, motherID int, status domain.Status) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE mother_profile SET status = ?, updated_at = ? WHERE mother_id = ?",
		string(status), time.Now(), motherID,
	)
	return err
}
