package command

import (
	"context"

	"backend/internal/document/domain"
)

type Sign struct {
	DocumentID int
	SignedAt   string
}

type SignHandler interface {
	Handle(ctx context.Context, cmd Sign) (int, error)
}

func NewSignHandler(repo domain.Repository) SignHandler {
	return &signHandler{repo: repo}
}

type signHandler struct {
	repo domain.Repository
}

func (h *signHandler) Handle(ctx context.Context, cmd Sign) (int, error) {
	d, err := h.repo.ByID(ctx, cmd.DocumentID)
	if err != nil {
		return 0, err
	}

	if err := d.Sign(cmd.SignedAt); err != nil {
		return 0, err
	}

	if err := h.repo.Update(ctx, d); err != nil {
		return 0, err
	}

	return d.ID(), nil
}
