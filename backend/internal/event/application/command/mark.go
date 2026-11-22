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
	if !domain.ParticipationStatus(cmd.Status).Valid() {
		return domain.ErrInvalidStatus
	}

	return h.uow.Do(ctx, func(ctx context.Context, repo domain.Repository) error {
		_, err := repo.ByID(ctx, cmd.EventID)
		if err != nil {
			return err
		}

		_, err = repo.GetParticipation(ctx, cmd.EventID, cmd.MotherID)
		if err != nil {
			return err
		}

		return repo.SetParticipationStatus(ctx, cmd.EventID, cmd.MotherID, domain.ParticipationStatus(cmd.Status))
	})
}
