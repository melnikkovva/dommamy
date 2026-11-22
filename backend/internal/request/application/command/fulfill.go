package command

import (
	"context"
	"time"

	"backend/internal/request/domain"
)

type FulfillRequest struct {
	RequestID int
	HelpID    int
}

type FulfillRequestHandler interface {
	Handle(ctx context.Context, cmd FulfillRequest) error
}

func NewFulfillRequestHandler(uow UnitOfWork, helps HelpDirectory, now func() time.Time) FulfillRequestHandler {
	return &fulfillRequestHandler{uow: uow, helps: helps, now: now}
}

type fulfillRequestHandler struct {
	uow   UnitOfWork
	helps HelpDirectory
	now   func() time.Time
}

func (h *fulfillRequestHandler) Handle(ctx context.Context, cmd FulfillRequest) error {
	return h.uow.Do(ctx, func(ctx context.Context, repo domain.Repository) error {
		req, err := repo.ByIDForUpdate(ctx, cmd.RequestID)
		if err != nil {
			return err
		}
		motherID, err := h.helps.MotherOf(ctx, cmd.HelpID)
		if err != nil {
			return err
		}
		if motherID != req.MotherID() {
			return domain.ErrHelpMismatch
		}
		if err := req.Fulfill(cmd.HelpID, h.now()); err != nil {
			return err
		}
		return repo.Update(ctx, req)
	})
}