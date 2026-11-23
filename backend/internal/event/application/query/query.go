package query

import (
	"context"
	"time"

	"backend/internal/event/domain"
)

type EventView struct {
	ID        int       `json:"id"`
	TypeID    int       `json:"type_id"`
	TypeName  string    `json:"type_name"`
	Title     string    `json:"title"`
	StartsAt  time.Time `json:"starts_at"`
	Place     string    `json:"place"`
	MaxPeople int       `json:"max_people"`
}

type ParticipationView struct {
	MotherID int                         `json:"mother_id"`
	Status   domain.ParticipationStatus  `json:"status"`
}

type EventTypeView struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ReadModel interface {
	GetEvent(ctx context.Context, id int) (EventView, error)
	ListEvents(ctx context.Context) ([]EventView, error)
	EventsByMother(ctx context.Context, motherID int) ([]EventView, error)
	ListEventTypes(ctx context.Context) ([]EventTypeView, error)
	Participants(ctx context.Context, eventID int) ([]ParticipationView, error)
}
