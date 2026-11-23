package mysql

import (
	"context"
	"database/sql"
	"errors"

	"backend/internal/event/application/query"
	"backend/internal/event/domain"
)

func NewReadModel(db *sql.DB) query.ReadModel { return &readModel{db: db} }

type readModel struct{ db *sql.DB }

func (r *readModel) GetEvent(ctx context.Context, id int) (query.EventView, error) {
	var v query.EventView
	err := r.db.QueryRowContext(ctx,
		"SELECT e.id, e.type_id, et.name, e.title, e.starts_at, e.place, e.max_people FROM event e LEFT JOIN event_type et ON et.id = e.type_id WHERE e.id = ?", id).
		Scan(&v.ID, &v.TypeID, &v.TypeName, &v.Title, &v.StartsAt, &v.Place, &v.MaxPeople)
	if errors.Is(err, sql.ErrNoRows) {
		return query.EventView{}, domain.ErrNotFound
	}
	if err != nil {
		return query.EventView{}, err
	}
	return v, nil
}

func (r *readModel) ListEvents(ctx context.Context) ([]query.EventView, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT e.id, e.type_id, et.name, e.title, e.starts_at, e.place, e.max_people FROM event e LEFT JOIN event_type et ON et.id = e.type_id ORDER BY e.starts_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.EventView
	for rows.Next() {
		var v query.EventView
		if err := rows.Scan(&v.ID, &v.TypeID, &v.TypeName, &v.Title, &v.StartsAt, &v.Place, &v.MaxPeople); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *readModel) EventsByMother(ctx context.Context, motherID int) ([]query.EventView, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT e.id, e.type_id, et.name, e.title, e.starts_at, e.place, e.max_people FROM event e JOIN participation p ON p.event_id = e.id LEFT JOIN event_type et ON et.id = e.type_id WHERE p.mother_id = ? ORDER BY e.starts_at DESC", motherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.EventView
	for rows.Next() {
		var v query.EventView
		if err := rows.Scan(&v.ID, &v.TypeID, &v.TypeName, &v.Title, &v.StartsAt, &v.Place, &v.MaxPeople); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *readModel) ListEventTypes(ctx context.Context) ([]query.EventTypeView, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM event_type ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.EventTypeView
	for rows.Next() {
		var v query.EventTypeView
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *readModel) Participants(ctx context.Context, eventID int) ([]query.ParticipationView, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT mother_id, status FROM participation WHERE event_id = ? ORDER BY created_at", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []query.ParticipationView
	for rows.Next() {
		var v query.ParticipationView
		var status string
		if err := rows.Scan(&v.MotherID, &status); err != nil {
			return nil, err
		}
		v.Status = domain.ParticipationStatus(status)
		result = append(result, v)
	}
	return result, rows.Err()
}
