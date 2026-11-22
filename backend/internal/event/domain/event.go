package domain

import (
	"slices"
	"strings"
	"time"
)

type ParticipationStatus string

const (
	StatusRegistered ParticipationStatus = "registered"
	StatusCame       ParticipationStatus = "came"
	StatusAbsent     ParticipationStatus = "absent"
)

func (s ParticipationStatus) Valid() bool {
	return slices.Contains([]ParticipationStatus{StatusRegistered, StatusCame, StatusAbsent}, s)
}

type Event struct {
	id        int
	typeID    int
	title     string
	startsAt  time.Time
	place     string
	maxPeople int
}

func NewEvent(typeID int, title string, startsAt time.Time, place string, maxPeople int) (Event, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Event{}, ErrTitle
	}
	if maxPeople < 0 {
		return Event{}, ErrLimit
	}
	return Event{
		typeID:    typeID,
		title:     title,
		startsAt:  startsAt,
		place:     strings.TrimSpace(place),
		maxPeople: maxPeople,
	}, nil
}

func (e *Event) SetID(id int) {
	e.id = id
}

func (e Event) ID() int        { return e.id }
func (e Event) TypeID() int    { return e.typeID }
func (e Event) Title() string  { return e.title }
func (e Event) StartsAt() time.Time { return e.startsAt }
func (e Event) Place() string  { return e.place }
func (e Event) MaxPeople() int { return e.maxPeople }

type Snapshot struct {
	ID        int
	TypeID    int
	Title     string
	StartsAt  time.Time
	Place     string
	MaxPeople int
}

func Restore(s Snapshot) Event {
	return Event{
		id:        s.ID,
		typeID:    s.TypeID,
		title:     s.Title,
		startsAt:  s.StartsAt,
		place:     s.Place,
		maxPeople: s.MaxPeople,
	}
}

func (e Event) Snapshot() Snapshot {
	return Snapshot{
		ID:        e.id,
		TypeID:    e.typeID,
		Title:     e.title,
		StartsAt:  e.startsAt,
		Place:     e.place,
		MaxPeople: e.maxPeople,
	}
}
