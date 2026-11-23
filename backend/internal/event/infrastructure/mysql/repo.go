package mysql

import (
	"context"
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"

	"backend/internal/event/domain"
)

const eventColumns = "e.id, e.type_id, e.title, e.starts_at, e.place, e.max_people"

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewRepository(q DBTX) domain.Repository { return &repository{q: q} }

type repository struct{ q DBTX }

func (r *repository) Add(ctx context.Context, e domain.Event) (id int, err error) {
	res, err := r.q.ExecContext(ctx,
		"INSERT INTO event (type_id, title, starts_at, place, max_people) VALUES (?, ?, ?, ?, ?)",
		e.TypeID(), e.Title(), e.StartsAt(), e.Place(), e.MaxPeople(),
	)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastID), nil
}

func (r *repository) ByID(ctx context.Context, id int) (domain.Event, error) {
	return r.byID(ctx, id, "")
}

func (r *repository) ByIDForUpdate(ctx context.Context, id int) (domain.Event, error) {
	return r.byID(ctx, id, " FOR UPDATE")
}

func (r *repository) byID(ctx context.Context, id int, lock string) (domain.Event, error) {
	var s domain.Snapshot
	err := r.q.QueryRowContext(ctx, "SELECT "+eventColumns+" FROM event e WHERE e.id = ?"+lock, id).
		Scan(&s.ID, &s.TypeID, &s.Title, &s.StartsAt, &s.Place, &s.MaxPeople)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Event{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Event{}, err
	}
	return domain.Restore(s), nil
}

func (r *repository) CountParticipants(ctx context.Context, eventID int) (int, error) {
	var n int
	err := r.q.QueryRowContext(ctx, "SELECT COUNT(*) FROM participation WHERE event_id = ?", eventID).Scan(&n)
	return n, err
}

func (r *repository) AddParticipant(ctx context.Context, eventID int, motherID int) error {
	_, err := r.q.ExecContext(ctx,
		"INSERT INTO participation (event_id, mother_id, status) VALUES (?, ?, 'registered')",
		eventID, motherID,
	)
	if err != nil {
		var me *mysqldriver.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return domain.ErrAlreadyJoined
		}
		return err
	}
	return nil
}

func (r *repository) HasParticipation(ctx context.Context, eventID int, motherID int) (bool, error) {
	var exists int
	err := r.q.QueryRowContext(ctx,
		"SELECT 1 FROM participation WHERE event_id = ? AND mother_id = ? LIMIT 1",
		eventID, motherID,
	).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *repository) GetParticipation(ctx context.Context, eventID int, motherID int) (domain.ParticipationStatus, error) {
	var status string
	err := r.q.QueryRowContext(ctx,
		"SELECT status FROM participation WHERE event_id = ? AND mother_id = ?",
		eventID, motherID,
	).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrNotRegistered
	}
	if err != nil {
		return "", err
	}
	return domain.ParticipationStatus(status), nil
}

func (r *repository) SetParticipationStatus(ctx context.Context, eventID int, motherID int, status domain.ParticipationStatus) error {
	_, err := r.q.ExecContext(ctx,
		"UPDATE participation SET status = ? WHERE event_id = ? AND mother_id = ?",
		status, eventID, motherID,
	)
	return err
}
