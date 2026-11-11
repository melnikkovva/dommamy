package mysql

import (
	"database/sql"

	"backend/internal/report"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) AidRows(from, to string) ([]report.AidRow, error) {
	rows, err := r.db.Query(
		"SELECT DATE_FORMAT(h.created_at, '%Y-%m-%d'), h.mother_id, CONCAT(m.last_name, ' ', m.first_name), h.kind, h.`sum`, h.check_url "+
			"FROM help_items h JOIN mother m ON m.id = h.mother_id "+
			"WHERE h.created_at >= ? AND h.created_at < DATE_ADD(?, INTERVAL 1 DAY) "+
			"ORDER BY h.created_at DESC, h.id DESC", from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []report.AidRow{}
	for rows.Next() {
		var a report.AidRow
		if err := rows.Scan(&a.Date, &a.MotherId, &a.Mother, &a.Kind, &a.Sum, &a.CheckUrl); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

func (r *Repo) AidTotals(from, to string) ([]report.KindTotal, error) {
	rows, err := r.db.Query(
		"SELECT kind, COUNT(*), COALESCE(SUM(`sum`), 0) FROM help_items "+
			"WHERE created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY) GROUP BY kind ORDER BY kind", from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []report.KindTotal{}
	for rows.Next() {
		var t report.KindTotal
		if err := rows.Scan(&t.Kind, &t.Count, &t.Sum); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *Repo) EventRows(from, to string) ([]report.EventRow, error) {
	rows, err := r.db.Query(
		"SELECT DATE_FORMAT(e.starts_at, '%Y-%m-%d %H:%i'), e.title, t.name, e.place, p.mother_id, CONCAT(m.last_name, ' ', m.first_name), p.status "+
			"FROM participation p JOIN event e ON e.id = p.event_id JOIN event_type t ON t.id = e.type_id JOIN mother m ON m.id = p.mother_id "+
			"WHERE e.starts_at >= ? AND e.starts_at < DATE_ADD(?, INTERVAL 1 DAY) "+
			"ORDER BY e.starts_at DESC, m.last_name", from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []report.EventRow{}
	for rows.Next() {
		var e report.EventRow
		if err := rows.Scan(&e.When, &e.Title, &e.Type, &e.Place, &e.MotherId, &e.Mother, &e.Status); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *Repo) MotherRows(from, to string) ([]report.MotherRow, error) {
	rows, err := r.db.Query(
		"SELECT m.id, CONCAT(m.last_name, ' ', m.first_name), m.phone, DATE_FORMAT(m.created_at, '%Y-%m-%d'), "+
			"(SELECT COUNT(*) FROM child c WHERE c.mother_id = m.id), "+
			"COALESCE(p.situation, ''), COALESCE(p.income_source, ''), COALESCE(p.income_amount, 0), COALESCE(p.status, ''), "+
			"COALESCE((SELECT s.status FROM mother_status_history s WHERE s.mother_id = m.id ORDER BY s.id DESC LIMIT 1), 'review'), "+
			"(SELECT COUNT(*) FROM help_items i WHERE i.mother_id = m.id) "+
			"FROM mother m LEFT JOIN mother_profile p ON p.mother_id = m.id "+
			"WHERE m.created_at >= ? AND m.created_at < DATE_ADD(?, INTERVAL 1 DAY) ORDER BY m.id", from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []report.MotherRow{}
	for rows.Next() {
		var m report.MotherRow
		err := rows.Scan(&m.Id, &m.Name, &m.Phone, &m.Registered, &m.Children, &m.Situation,
			&m.IncomeSource, &m.Income, &m.ProfileStatus, &m.Status, &m.HelpCount)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}