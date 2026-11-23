package command

import (
	"context"

	"backend/internal/mother/domain"
)

type SoftDelete struct {
	MotherID int
}

type SoftDeleteHandler interface {
	Handle(ctx context.Context, cmd SoftDelete) error
}

func NewSoftDeleteHandler(repo domain.Repository) SoftDeleteHandler {
	return &softDeleteHandler{repo: repo}
}

type softDeleteHandler struct {
	repo domain.Repository
}

func (h *softDeleteHandler) Handle(ctx context.Context, cmd SoftDelete) error {
	return h.repo.SoftDelete(ctx, cmd.MotherID)
}
