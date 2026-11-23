package domain

import (
	"slices"
	"strings"
	"time"
)

type Status string

const (
	StatusReview   Status = "review"
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

var validStatuses = []Status{StatusReview, StatusActive, StatusInactive}

func (s Status) Valid() bool {
	return slices.Contains(validStatuses, s)
}

var transitions = map[Status][]Status{
	StatusReview:   {StatusActive, StatusInactive},
	StatusActive:   {StatusInactive},
	StatusInactive: {StatusReview},
}

func CanTransition(from, to Status) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

type Change struct {
	id        int
	motherID  int
	status    Status
	reason    string
	createdAt time.Time
}

func newChange(motherID int, status Status, reason string, now time.Time) Change {
	return Change{
		motherID:  motherID,
		status:    status,
		reason:    strings.TrimSpace(reason),
		createdAt: now,
	}
}


func (c Change) ID() int             { return c.id }
func (c Change) MotherID() int       { return c.motherID }
func (c Change) Status() Status      { return c.status }
func (c Change) Reason() string      { return c.reason }
func (c Change) CreatedAt() time.Time { return c.createdAt }

type StatusHistory struct {
	motherID  int
	current   Status
	history   []Change
}

func NewStatusHistory(motherID int) (StatusHistory, error) {
	if motherID <= 0 {
		return StatusHistory{}, ErrMother
	}
	return StatusHistory{
		motherID: motherID,
		current:  StatusReview,
		history:  []Change{},
	}, nil
}

func (h *StatusHistory) Transition(to Status, reason string, now time.Time) error {
	if !to.Valid() {
		return ErrStatus
	}
	if !CanTransition(h.current, to) {
		return ErrTransition
	}
	if to == StatusInactive && strings.TrimSpace(reason) == "" {
		return ErrReason
	}
	h.current = to
	h.history = append(h.history, newChange(h.motherID, to, reason, now))
	return nil
}

func (h StatusHistory) MotherID() int       { return h.motherID }
func (h StatusHistory) Current() Status     { return h.current }
func (h StatusHistory) History() []Change   { return h.history }

type ChangeSnapshot struct {
	ID        int
	MotherID  int
	Status    Status
	Reason    string
	CreatedAt time.Time
}

func (c Change) Snapshot() ChangeSnapshot {
	return ChangeSnapshot{
		ID:        c.id,
		MotherID:  c.motherID,
		Status:    c.status,
		Reason:    c.reason,
		CreatedAt: c.createdAt,
	}
}

func restoreChange(s ChangeSnapshot) Change {
	return Change{
		id:        s.ID,
		motherID:  s.MotherID,
		status:    s.Status,
		reason:    s.Reason,
		createdAt: s.CreatedAt,
	}
}

func RestoreHistory(motherID int, changes []ChangeSnapshot) StatusHistory {
	h := StatusHistory{
		motherID: motherID,
		current:  StatusReview,
		history:  make([]Change, 0, len(changes)),
	}
	for _, s := range changes {
		c := restoreChange(s)
		h.current = c.status
		h.history = append(h.history, c)
	}
	return h
}
