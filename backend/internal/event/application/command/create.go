package command

import (
	"context"
	"time"

	"backend/internal/event/domain"
)

type CreateEvent struct {
	TypeID    int
	Title     string
	StartsAt  time.Time
	Place     string
	MaxPeople int
}

type CreateEventHandler interface {
	Handle(ctx context.Context, cmd CreateEvent) (id int, err error)
}

func NewCreateEventHandler(repo domain.Repository) CreateEventHandler {
	return &createEventHandler{repo: repo}
}

type createEventHandler struct {
	repo domain.Repository
}

func (h *createEventHandler) Handle(ctx context.Context, cmd CreateEvent) (id int, err error) {
	e, err := domain.NewEvent(cmd.TypeID, cmd.Title, cmd.StartsAt, cmd.Place, cmd.MaxPeople)
	if err != nil {
		return 0, err
	}
	return h.repo.Add(ctx, e)
}
