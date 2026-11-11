package mysql

import (
	"database/sql"
	"errors"

	"backend/internal/donor"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveDonor(d *donor.Donor) error {
	res, err := r.db.Exec(
		"INSERT INTO donor (name, kind, phone, email, comment, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		d.Name, d.Kind, d.Phone, d.Email, d.Comment, d.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.Id = int(id)
	return nil
}

func (r *Repo) UpdateDonor(d *donor.Donor) error {
	_, err := r.db.Exec(
		"UPDATE donor SET name = ?, kind = ?, phone = ?, email = ?, comment = ? WHERE id = ?",
		d.Name, d.Kind, d.Phone, d.Email, d.Comment, d.Id,
	)
	return err
}

func (r *Repo) DonorById(id int) (*donor.Donor, error) {
	d := &donor.Donor{}
	err := r.db.QueryRow(
		"SELECT id, name, kind, phone, email, comment, created_at FROM donor WHERE id = ?", id,
	).Scan(&d.Id, &d.Name, &d.Kind, &d.Phone, &d.Email, &d.Comment, &d.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, donor.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *Repo) ListDonors(q string) ([]*donor.Summary, error) {
	rows, err := r.db.Query(
		`SELECT d.id, d.name, d.kind, d.phone, d.email, d.comment, d.created_at,
		        COALESCE(SUM(n.amount), 0), COUNT(n.id)
		 FROM donor d
		 LEFT JOIN donation n ON n.donor_id = d.id
		 WHERE d.name LIKE CONCAT('%', ?, '%') OR d.phone LIKE CONCAT('%', ?, '%')
		 GROUP BY d.id
		 ORDER BY d.name`, q, q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*donor.Summary{}
	for rows.Next() {
		s := &donor.Summary{}
		err := rows.Scan(&s.Id, &s.Name, &s.Kind, &s.Phone, &s.Email, &s.Comment, &s.CreatedAt,
			&s.Total, &s.Count)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repo) SaveDonation(n *donor.Donation) error {
	res, err := r.db.Exec(
		"INSERT INTO donation (donor_id, amount, donated_at, comment, created_at) VALUES (?, ?, ?, ?, ?)",
		n.DonorId, n.Amount, n.Date, n.Comment, n.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	n.Id = int(id)
	return nil
}

func (r *Repo) DonationsByDonor(donorId int) ([]*donor.Donation, error) {
	rows, err := r.db.Query(
		`SELECT id, donor_id, amount, DATE_FORMAT(donated_at, '%Y-%m-%d'), comment, created_at
		 FROM donation WHERE donor_id = ? ORDER BY donated_at DESC, id DESC`, donorId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*donor.Donation{}
	for rows.Next() {
		n := &donor.Donation{}
		if err := rows.Scan(&n.Id, &n.DonorId, &n.Amount, &n.Date, &n.Comment, &n.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, rows.Err()
}