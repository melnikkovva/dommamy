package command

import (
	"context"
	"time"

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

	today := time.Now().Format("2006-01-02")
	if err := d.Sign(cmd.SignedAt, today); err != nil {
		return 0, err
	}

	if err := h.repo.Update(ctx, d); err != nil {
		return 0, err
	}

	return d.ID(), nil
}
