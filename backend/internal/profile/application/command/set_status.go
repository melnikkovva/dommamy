package command

import (
	"context"
	"errors"

	"backend/internal/profile/domain"
)

type SetStatus struct {
	MotherID int
	Status   string
}

type SetStatusHandler interface {
	Handle(ctx context.Context, cmd SetStatus) error
}

func NewSetStatusHandler(repo domain.Repository) SetStatusHandler {
	return &setStatusHandler{repo: repo}
}

type setStatusHandler struct {
	repo domain.Repository
}

func (h *setStatusHandler) Handle(ctx context.Context, cmd SetStatus) error {
	status := domain.Status(cmd.Status)
	if !status.Valid() {
		return domain.ErrStatus
	}

	p, err := h.repo.ByMother(ctx, cmd.MotherID)
	if errors.Is(err, domain.ErrNotFilled) {
		return domain.ErrNotFilled
	}
	if err != nil {
		return err
	}

	if status == domain.StatusSent {
		if err := p.Send(); err != nil {
			return err
		}
	} else {
		if err := p.Transition(status); err != nil {
			return err
		}
	}

	return h.repo.UpdateStatus(ctx, cmd.MotherID, p.Status())
}
