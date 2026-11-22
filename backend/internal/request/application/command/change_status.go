package command

import (
	"context"
	"time"

	"backend/internal/request/domain"
)

type ChangeStatus struct {
	RequestID int
	Status    string
	Comment   string
}

type ChangeStatusHandler interface {
	Handle(ctx context.Context, cmd ChangeStatus) error
}

func NewChangeStatusHandler(uow UnitOfWork, now func() time.Time) ChangeStatusHandler {
	return &changeStatusHandler{uow: uow, now: now}
}

type changeStatusHandler struct {
	uow UnitOfWork
	now func() time.Time
}

func (h *changeStatusHandler) Handle(ctx context.Context, cmd ChangeStatus) error {
	return h.uow.Do(ctx, func(ctx context.Context, repo domain.Repository) error {
		req, err := repo.ByIDForUpdate(ctx, cmd.RequestID)
		if err != nil {
			return err
		}
		if err := req.ChangeStatus(domain.Status(cmd.Status), cmd.Comment, h.now()); err != nil {
			return err
		}
		return repo.Update(ctx, req)
	})
}