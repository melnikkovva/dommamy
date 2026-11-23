package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	motherquery "backend/internal/mother/application/query"
	"backend/internal/mother/domain"
)

type ReadModel struct {
	db         *sql.DB
	learnStats motherquery.LearnStats
}

func NewReadModel(db *sql.DB, learnStats motherquery.LearnStats) *ReadModel {
	return &ReadModel{db: db, learnStats: learnStats}
}

func (r *ReadModel) GetByID(ctx context.Context, id int) (motherquery.MotherView, error) {
	s, err := r.snapshotByID(ctx, id)
	if err != nil {
		return motherquery.MotherView{}, err
	}
	view := toMotherView(s)
	if s.LearnUserID != "" && r.learnStats != nil {
		stats, err := r.learnStats.GetCourseStats(ctx, s.LearnUserID)
		if err == nil {
			view.CourseStats = &stats
		}
	}
	return view, nil
}

func (r *ReadModel) Search(ctx context.Context, f motherquery.MotherFilter) ([]motherquery.MotherView, error) {
	var conditions []string
	var args []interface{}

	if f.Q != "" {
		conditions = append(conditions, "(CONCAT(m.first_name, ' ', m.last_name) LIKE ? OR m.phone LIKE ?)")
		like := "%" + f.Q + "%"
		args = append(args, like, like)
	}

	if f.Phone != "" {
		conditions = append(conditions, "m.phone LIKE ?")
		args = append(args, "%"+f.Phone+"%")
	}

	if f.ProfileStatus != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, f.ProfileStatus)
	}

	if f.MotherStatus != "" {
		conditions = append(conditions, "ms.status = ?")
		args = append(args, f.MotherStatus)
	}

	where := "m.deleted_at IS NULL"
	if len(conditions) > 0 {
		where += " AND " + strings.Join(conditions, " AND ")
	}

	query := `SELECT m.id, m.first_name, m.last_name, m.phone, m.email,
		        mll.learn_user_id, m.photo_url, m.created_at,
		        COALESCE(p.status, ''),
		        COALESCE(ms.status, 'review')
		 FROM mother m
		 LEFT JOIN mother_learn_link mll ON mll.mother_id = m.id
		 LEFT JOIN mother_profile p ON p.mother_id = m.id
		 LEFT JOIN (
		   SELECT h.mother_id, h.status
		     FROM mother_status_history h
		     JOIN (SELECT MAX(id) AS last_id FROM mother_status_history GROUP BY mother_id) m2
		       ON h.id = m2.last_id
		 ) ms ON ms.mother_id = m.id
		 WHERE ` + where + `
		 ORDER BY m.last_name, m.first_name`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []motherquery.MotherView
	for rows.Next() {
		s, ps, ms, err := scanSearchRow(rows)
		if err != nil {
			return nil, err
		}
		view := toMotherView(s)
		view.ProfileStatus = ps
		view.MotherStatus = ms
		if s.LearnUserID != "" && r.learnStats != nil {
			stats, err := r.learnStats.GetCourseStats(ctx, s.LearnUserID)
			if err != nil {
				log.Printf("ReadModel.Search: ошибка получения статистики для mother %d: %v", s.ID, err)
			} else {
				view.CourseStats = &stats
			}
		}
		views = append(views, view)
	}
	return views, rows.Err()
}

func (r *ReadModel) Exists(ctx context.Context, id int) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM mother WHERE id = ? AND deleted_at IS NULL", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ReadModel) GetAllMothers(ctx context.Context) ([]motherquery.MotherRef, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, first_name, last_name, phone, email FROM mother WHERE deleted_at IS NULL ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refs []motherquery.MotherRef
	for rows.Next() {
		var ref motherquery.MotherRef
		var email sql.NullString
		if err := rows.Scan(&ref.ID, &ref.FirstName, &ref.LastName, &ref.Phone, &email); err != nil {
			return nil, err
		}
		if email.Valid {
			ref.Email = email.String
		}
		refs = append(refs, ref)
	}
	return refs, rows.Err()
}

func (r *ReadModel) snapshotByID(ctx context.Context, id int) (domain.MotherSnapshot, error) {
	var s domain.MotherSnapshot
	var learnUserID, photoURL sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT m.id, m.first_name, m.last_name, m.phone, m.email,
		        mll.learn_user_id, m.photo_url, m.created_at
		 FROM mother m
		 LEFT JOIN mother_learn_link mll ON mll.mother_id = m.id
		 WHERE m.id = ? AND m.deleted_at IS NULL`, id,
	).Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &learnUserID, &photoURL, &s.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.MotherSnapshot{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.MotherSnapshot{}, err
	}
	if learnUserID.Valid {
		s.LearnUserID = learnUserID.String
	}
	if photoURL.Valid {
		s.PhotoURL = photoURL.String
	}
	return s, nil
}

func scanSearchRow(rows *sql.Rows) (domain.MotherSnapshot, string, string, error) {
	var s domain.MotherSnapshot
	var learnUserID, photoURL sql.NullString
	var profileStatus, motherStatus string
	if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &learnUserID, &photoURL, &s.CreatedAt, &profileStatus, &motherStatus); err != nil {
		return domain.MotherSnapshot{}, "", "", err
	}
	if learnUserID.Valid {
		s.LearnUserID = learnUserID.String
	}
	if photoURL.Valid {
		s.PhotoURL = photoURL.String
	}
	return s, profileStatus, motherStatus, nil
}

func scanSnapshot(rows *sql.Rows) (domain.MotherSnapshot, error) {
	var s domain.MotherSnapshot
	var learnUserID sql.NullString
	var photoURL sql.NullString
	if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &learnUserID, &photoURL, &s.CreatedAt); err != nil {
		return domain.MotherSnapshot{}, err
	}
	if learnUserID.Valid {
		s.LearnUserID = learnUserID.String
	}
	if photoURL.Valid {
		s.PhotoURL = photoURL.String
	}
	return s, nil
}

func toMotherView(s domain.MotherSnapshot) motherquery.MotherView {
	return motherquery.MotherView{
		ID:          s.ID,
		FirstName:   s.FirstName,
		LastName:    s.LastName,
		Phone:       s.Phone,
		Email:       s.Email,
		LearnUserID: s.LearnUserID,
		PhotoURL:    s.PhotoURL,
		CreatedAt:   s.CreatedAt,
	}
}
