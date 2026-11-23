package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"

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

func (r *ReadModel) Search(ctx context.Context, q string) ([]motherquery.MotherView, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT m.id, m.first_name, m.last_name, m.phone, m.email,
		        mll.learn_user_id, m.photo_url, m.created_at
		 FROM mother m
		 LEFT JOIN mother_learn_link mll ON mll.mother_id = m.id
		 WHERE CONCAT(m.first_name, ' ', m.last_name) LIKE CONCAT('%', ?, '%')
		 ORDER BY m.last_name, m.first_name`,
		q,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []motherquery.MotherView
	for rows.Next() {
		s, err := scanSnapshot(rows)
		if err != nil {
			return nil, err
		}
		view := toMotherView(s)
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
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM mother WHERE id = ?", id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ReadModel) GetAllMothers(ctx context.Context) ([]motherquery.MotherRef, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, first_name, last_name, phone, email FROM mother ORDER BY id`,
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
		 WHERE m.id = ?`, id,
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
