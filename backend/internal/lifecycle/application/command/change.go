package command

import (
	"context"
	"time"

	"backend/internal/lifecycle/domain"
)

type Change struct {
	MotherID int
	Status   string
	Reason   string
}

type ChangeHandler interface {
	Handle(ctx context.Context, cmd Change) error
}

func NewChangeHandler(repo domain.Repository, mothers MotherDirectory) ChangeHandler {
	return &changeHandler{repo: repo, mothers: mothers}
}

type changeHandler struct {
	repo    domain.Repository
	mothers MotherDirectory
}

func (h *changeHandler) Handle(ctx context.Context, cmd Change) error {
	status := domain.Status(cmd.Status)
	if !status.Valid() {
		return domain.ErrStatus
	}

	ok, err := h.mothers.Exists(ctx, cmd.MotherID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNoMother
	}

	history, err := h.repo.History(ctx, cmd.MotherID)
	if err != nil {
		return err
	}

	aggregate := domain.RestoreHistory(cmd.MotherID, history)

	if err := aggregate.Transition(status, cmd.Reason, time.Now()); err != nil {
		return err
	}

	changes := aggregate.History()
	if len(changes) == 0 {
		return domain.ErrTransition
	}
	newest := changes[len(changes)-1]

	_, err = h.repo.Add(ctx, newest)
	return err
}
