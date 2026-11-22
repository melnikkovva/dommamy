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
		has, err := repo.HasParticipation(ctx, cmd.EventID, cmd.MotherID)
		if err != nil {
			return err
		}
		if has {
			return domain.ErrAlreadyJoined
		}

		e, err := repo.ByID(ctx, cmd.EventID)
		if err != nil {
			return err
		}

		if e.MaxPeople() > 0 {
			count, err := repo.CountParticipants(ctx, cmd.EventID)
			if err != nil {
				return err
			}
			if count >= e.MaxPeople() {
				return domain.ErrFull
			}
		}

		return repo.AddParticipant(ctx, cmd.EventID, cmd.MotherID)
	})
}
