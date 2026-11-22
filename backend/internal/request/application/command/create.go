package command

import (
	"context"
	"time"

	"backend/internal/request/domain"
)

type CreateRequest struct {
	MotherID    int
	Kind        string
	Description string
	Source      string
}

type CreateRequestHandler interface {
	Handle(ctx context.Context, cmd CreateRequest) (id int, err error)
}

func NewCreateRequestHandler(repo domain.Repository, mothers MotherDirectory, now func() time.Time) CreateRequestHandler {
	return &createRequestHandler{repo: repo, mothers: mothers, now: now}
}

type createRequestHandler struct {
	repo    domain.Repository
	mothers MotherDirectory
	now     func() time.Time
}

func (h *createRequestHandler) Handle(ctx context.Context, cmd CreateRequest) (int, error) {
	req, err := domain.NewRequest(cmd.MotherID, domain.Kind(cmd.Kind), cmd.Description, domain.Source(cmd.Source), h.now())
	if err != nil {
		return 0, err
	}
	ok, err := h.mothers.Exists(ctx, cmd.MotherID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, domain.ErrMotherNotFound
	}
	return h.repo.Add(ctx, req)
}