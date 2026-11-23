package command

import (
	"context"
	"time"

	"backend/internal/document/domain"
)

type Create struct {
	MotherID int
	Type     string
	HelpID   int
}

type CreateHandler interface {
	Handle(ctx context.Context, cmd Create) (int, error)
}

func NewCreateHandler(repo domain.Repository, mothers MotherDirectory, helps HelpDirectory) CreateHandler {
	return &createHandler{repo: repo, mothers: mothers, helps: helps}
}

type createHandler struct {
	repo    domain.Repository
	mothers MotherDirectory
	helps   HelpDirectory
}

func (h *createHandler) Handle(ctx context.Context, cmd Create) (int, error) {
	d, err := domain.NewDocument(cmd.MotherID, domain.Type(cmd.Type), cmd.HelpID, time.Now())
	if err != nil {
		return 0, err
	}

	ok, err := h.mothers.Exists(ctx, cmd.MotherID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, domain.ErrNoMother
	}

	if d.Type() == domain.TypeResponsibility {
		owner, err := h.helps.MotherOf(ctx, cmd.HelpID)
		if err != nil {
			return 0, err
		}
		if owner != cmd.MotherID {
			return 0, domain.ErrHelpMismatch
		}
	}

	id, err := h.repo.Add(ctx, d)
	if err != nil {
		return 0, err
	}
	return id, nil
}
