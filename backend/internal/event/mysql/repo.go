package mysql 

import (
	"database/sql"
	"errors"

	"backend/internal/event"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Types() ([]*event.Type, error) {
	rows, err := r.db.Query("SELECT id, name FROM event_type ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*event.Type{}
	for rows.Next() {
		t := &event.Type{}
		if err := rows.Scan(&t.Id, &t.Name); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

func (r *Repo) SaveEvent(e *event.Event) error {
	res, err := r.db.Exec(
		"INSERT INTO event (type_id, title, starts_at, place, max_people) VALUES (?, ?, ?, ?, ?)",
		e.TypeId, e.Title, e.StartsAt, e.Place, e.MaxPeople,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.Id = int(id)
	return nil
}

const eventColumns = "e.id, e.type_id, e.title, e.starts_at, e.place, e.max_people"

func scanEvents(rows *sql.Rows) ([]*event.Event, error) {
	defer rows.Close()
	result := []*event.Event{}
	for rows.Next() {
		e := &event.Event{}
		err := rows.Scan(&e.Id, &e.TypeId, &e.Title, &e.StartsAt, &e.Place, &e.MaxPeople)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *Repo) Events() ([]*event.Event, error) {
	rows, err := r.db.Query("SELECT " + eventColumns + " FROM event e ORDER BY e.starts_at DESC")
	if err != nil {
		return nil, err
	}
	return scanEvents(rows)
}

func (r *Repo) EventById(id int) (*event.Event, error) {
	e := &event.Event{}
	err := r.db.QueryRow(
		"SELECT "+eventColumns+" FROM event e WHERE e.id = ?", id,
	).Scan(&e.Id, &e.TypeId, &e.Title, &e.StartsAt, &e.Place, &e.MaxPeople)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, event.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *Repo) EventsByMother(motherId int) ([]*event.Event, error) {
	rows, err := r.db.Query(
		"SELECT "+eventColumns+` FROM event e
		 JOIN participation p ON p.event_id = e.id
		 WHERE p.mother_id = ? ORDER BY e.starts_at DESC`, motherId)
	if err != nil {
		return nil, err
	}
	return scanEvents(rows)
}

func (r *Repo) Participation(eventId, motherId int) (*event.Participation, error) {
	p := &event.Participation{}
	err := r.db.QueryRow(
		"SELECT event_id, mother_id, status FROM participation WHERE event_id = ? AND mother_id = ?",
		eventId, motherId,
	).Scan(&p.EventId, &p.MotherId, &p.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, event.ErrorNotRegistered
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repo) CountParticipants(eventId int) (int, error) {
	var n int
	err := r.db.QueryRow("SELECT COUNT(*) FROM participation WHERE event_id = ?", eventId).Scan(&n)
	return n, err
}

func (r *Repo) AddParticipant(eventId, motherId int) error {
	_, err := r.db.Exec(
		"INSERT INTO participation (event_id, mother_id, status) VALUES (?, ?, 'registered')",
		eventId, motherId,
	)
	return err
}

func (r *Repo) SetStatus(eventId, motherId int, status string) error {
	_, err := r.db.Exec(
		"UPDATE participation SET status = ? WHERE event_id = ? AND mother_id = ?",
		status, eventId, motherId,
	)
	return err
}

func (r *Repo) Participants(eventId int) ([]*event.Participation, error) {
	rows, err := r.db.Query(
		"SELECT event_id, mother_id, status FROM participation WHERE event_id = ? ORDER BY created_at",
		eventId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*event.Participation{}
	for rows.Next() {
		p := &event.Participation{}
		if err := rows.Scan(&p.EventId, &p.MotherId, &p.Status); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, rows.Err()
}
