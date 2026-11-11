package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/profile"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(p *profile.Profile) error {
	_, err := r.db.Exec(
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
		 updated_at = VALUES(updated_at)`,
		p.MotherId, p.BirthDate, p.BirthPlace, p.Citizenship, p.Registration, p.Living, p.Email,
		p.Religion, p.Education, p.Profession, p.IncomeSource, p.IncomeAmount, p.MaritalStatus,
		p.FatherInfo, p.PregnancyWeeks, p.Situation, p.ExpectedHelp, p.Source, p.Status, p.UpdatedAt,
	)
	return err
}

func (r *Repo) ByMother(motherId int) (*profile.Profile, error) {
	p := &profile.Profile{}
	err := r.db.QueryRow(
		`SELECT mother_id, COALESCE(DATE_FORMAT(birth_date, '%Y-%m-%d'), ''), birth_place,
		        citizenship, registration, living, email, religion, education, profession,
		        income_source, income_amount, marital_status, father_info, pregnancy_weeks,
		        situation, expected_help, source, status, updated_at
		 FROM mother_profile WHERE mother_id = ?`, motherId,
	).Scan(&p.MotherId, &p.BirthDate, &p.BirthPlace, &p.Citizenship, &p.Registration,
		&p.Living, &p.Email, &p.Religion, &p.Education, &p.Profession, &p.IncomeSource,
		&p.IncomeAmount, &p.MaritalStatus, &p.FatherInfo, &p.PregnancyWeeks, &p.Situation,
		&p.ExpectedHelp, &p.Source, &p.Status, &p.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repo) SetStatus(motherId int, status string) error {
	_, err := r.db.Exec(
		"UPDATE mother_profile SET status = ?, updated_at = NOW() WHERE mother_id = ?",
		status, motherId,
	)
	return err
}

func (r *Repo) ListByStatus(status string) ([]*profile.Profile, error) {
	rows, err := r.db.Query(
		`SELECT mother_id, status, updated_at FROM mother_profile
		 WHERE ? = '' OR status = ? ORDER BY updated_at DESC`, status, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*profile.Profile{}
	for rows.Next() {
		p := &profile.Profile{}
		if err := rows.Scan(&p.MotherId, &p.Status, &p.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}