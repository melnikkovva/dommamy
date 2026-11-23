package command

import (
	"context"
	"time"

	"backend/internal/mother/domain"
)

type Create struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

type CreateHandler interface {
	Handle(ctx context.Context, cmd Create) (domain.MotherSnapshot, error)
}

func NewCreateHandler(repo domain.Repository) CreateHandler {
	return &createHandler{repo: repo}
}

type createHandler struct {
	repo domain.Repository
}

func (h *createHandler) Handle(ctx context.Context, cmd Create) (domain.MotherSnapshot, error) {
	m, err := domain.NewMother(cmd.FirstName, cmd.LastName, cmd.Phone, cmd.Email, time.Now())
	if err != nil {
		return domain.MotherSnapshot{}, err
	}
	id, err := h.repo.Add(ctx, m)
	if err != nil {
		return domain.MotherSnapshot{}, err
	}
	snapshot := m.Snapshot()
	snapshot.ID = id
	return snapshot, nil
}
