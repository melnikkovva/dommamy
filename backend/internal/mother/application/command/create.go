package command

import (
	"context"

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
	m, err := domain.NewMother(cmd.FirstName, cmd.LastName, cmd.Phone, cmd.Email)
	if err != nil {
		return domain.MotherSnapshot{}, err
	}
	if err := h.repo.Save(ctx, m); err != nil {
		return domain.MotherSnapshot{}, err
	}
	return m.Snapshot(), nil
}
