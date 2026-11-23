package command

import (
	"context"

	"backend/internal/event/domain"
)

type Mark struct {
	EventID  int
	MotherID int
	Status   string
}

type MarkHandler interface {
	Handle(ctx context.Context, cmd Mark) error
}

func NewMarkHandler(uow UnitOfWork) MarkHandler {
	return &markHandler{uow: uow}
}

type markHandler struct {
	uow UnitOfWork
}

func (h *markHandler) Handle(ctx context.Context, cmd Mark) error {
	to := domain.ParticipationStatus(cmd.Status)
	if !to.Valid() {
		return domain.ErrInvalidStatus
	}

	return h.uow.Do(ctx, func(ctx context.Context, repo domain.Repository) error {
		_, err := repo.ByID(ctx, cmd.EventID)
		if err != nil {
			return err
		}

		current, err := repo.GetParticipation(ctx, cmd.EventID, cmd.MotherID)
		if err != nil {
			return err
		}

		if !domain.CanParticipationTransition(current, to) {
			return domain.ErrInvalidTransition
		}

		return repo.SetParticipationStatus(ctx, cmd.EventID, cmd.MotherID, to)
	})
}
