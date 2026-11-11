package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/secretInfo"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Save(d *secretInfo.Data) error {
	_, err := r.db.Exec(
		`INSERT INTO mother_secret_info
		 (mother_id, criminal_record, nd_pnd, drugs, chronic_diseases, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		 criminal_record = VALUES(criminal_record),
		 nd_pnd = VALUES(nd_pnd),
		 drugs = VALUES(drugs),
		 chronic_diseases = VALUES(chronic_diseases),
		 updated_at = VALUES(updated_at)`,
		d.MotherId, d.CriminalRecord, d.NdPnd, d.Drugs, d.ChronicDiseases, d.UpdatedAt,
	)
	return err
}

func (r *Repo) SearchByMother(motherId int) (*secretInfo.Data, error) {
	d := &secretInfo.Data{}
	err := r.db.QueryRow(
		`SELECT mother_id, criminal_record, nd_pnd, drugs, chronic_diseases, updated_at
		 FROM mother_secret_info WHERE mother_id = ?`, motherId,
	).Scan(&d.MotherId, &d.CriminalRecord, &d.NdPnd, &d.Drugs, &d.ChronicDiseases, &d.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, secretInfo.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}