package domain

import (
	"slices"
	"strings"
	"time"
)

type Kind string

const (
	KindFood      Kind = "food"
	KindClothes   Kind = "clothes"
	KindFurniture Kind = "furniture"
	KindOther     Kind = "other"
)

func (k Kind) Valid() bool {
	return slices.Contains([]Kind{KindFood, KindClothes, KindFurniture, KindOther}, k)
}

type Source string

const (
	SourceStaff  Source = "staff"
	SourceMother Source = "mother"
)

func (s Source) Valid() bool { return s == SourceStaff || s == SourceMother }

type Status string

const (
	StatusNew      Status = "new"
	StatusInWork   Status = "in_work"
	StatusIssued   Status = "issued"
	StatusRejected Status = "rejected"
)

var transitions = map[Status][]Status{
	StatusNew:      {StatusInWork, StatusRejected},
	StatusInWork:   {StatusRejected},
	StatusRejected: {StatusNew},
}

func (s Status) Valid() bool {
	return slices.Contains([]Status{StatusNew, StatusInWork, StatusIssued, StatusRejected}, s)
}

func (s Status) CanChangeTo(to Status) bool {
	return slices.Contains(transitions[s], to)
}

type Request struct {
	id, motherID, helpID int
	kind                 Kind
	description          string
	status               Status
	comment              string
	source               Source
	createdAt, updatedAt time.Time
}

func NewRequest(motherID int, kind Kind, description string, source Source, now time.Time) (Request, error) {
	if motherID <= 0 {
		return Request{}, ErrMotherNotSpecified
	}
	if !kind.Valid() {
		return Request{}, ErrInvalidKind
	}
	if !source.Valid() {
		return Request{}, ErrInvalidSource
	}
	return Request{
		motherID:    motherID,
		kind:        kind,
		description: strings.TrimSpace(description),
		status:      StatusNew,
		source:      source,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (r *Request) ChangeStatus(to Status, comment string, now time.Time) error {
	if to == StatusIssued {
		return ErrUseFulfill
	}
	if !to.Valid() {
		return ErrInvalidStatus
	}
	if !r.status.CanChangeTo(to) {
		return ErrTransition
	}
	comment = strings.TrimSpace(comment)
	if to == StatusRejected && comment == "" {
		return ErrCommentRequired
	}
	r.status = to
	if comment != "" {
		r.comment = comment
	}
	r.updatedAt = now
	return nil
}

func (r *Request) Fulfill(helpID int, now time.Time) error {
	if helpID <= 0 {
		return ErrHelpNotFound
	}
	if r.status != StatusNew && r.status != StatusInWork {
		return ErrTransition
	}
	r.status = StatusIssued
	r.helpID = helpID
	r.updatedAt = now
	return nil
}

func (r Request) MotherID() int { return r.motherID }

type Snapshot struct {
	ID, MotherID, HelpID                      int
	Kind, Description, Status, Comment, Source string
	CreatedAt, UpdatedAt                      time.Time
}

func Restore(s Snapshot) Request {
	return Request{
		id: s.ID, motherID: s.MotherID, helpID: s.HelpID,
		kind: Kind(s.Kind), description: s.Description,
		status: Status(s.Status), comment: s.Comment, source: Source(s.Source),
		createdAt: s.CreatedAt, updatedAt: s.UpdatedAt,
	}
}

func (r Request) Snapshot() Snapshot {
	return Snapshot{
		ID: r.id, MotherID: r.motherID, HelpID: r.helpID,
		Kind: string(r.kind), Description: r.description,
		Status: string(r.status), Comment: r.comment, Source: string(r.source),
		CreatedAt: r.createdAt, UpdatedAt: r.updatedAt,
	}
}