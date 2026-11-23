package command

import (
	"context"

	"backend/internal/event/domain"
)

type Register struct {
	EventID  int
	MotherID int
}

type RegisterHandler interface {
	Handle(ctx context.Context, cmd Register) error
}

func NewRegisterHandler(uow UnitOfWork, mothers MotherDirectory) RegisterHandler {
	return &registerHandler{uow: uow, mothers: mothers}
}

type registerHandler struct {
	uow     UnitOfWork
	mothers MotherDirectory
}

func (h *registerHandler) Handle(ctx context.Context, cmd Register) error {
	ok, err := h.mothers.Exists(ctx, cmd.MotherID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrMotherNotFound
	}

	return h.uow.Do(ctx, func(ctx context.Context, repo domain.Repository) error {
		e, err := repo.ByIDForUpdate(ctx, cmd.EventID)
		if err != nil {
			return err
		}

		count, err := repo.CountParticipants(ctx, cmd.EventID)
		if err != nil {
			return err
		}
		if err := e.EnsureHasPlace(count); err != nil {
			return err
		}

		return repo.AddParticipant(ctx, cmd.EventID, cmd.MotherID)
	})
}
