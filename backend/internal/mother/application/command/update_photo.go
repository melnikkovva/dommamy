package command

import (
	"context"

	"backend/internal/mother/domain"
)

type UpdatePhoto struct {
	MotherID int
	PhotoURL string
}

type UpdatePhotoHandler interface {
	Handle(ctx context.Context, cmd UpdatePhoto) error
}

func NewUpdatePhotoHandler(repo domain.Repository) UpdatePhotoHandler {
	return &updatePhotoHandler{repo: repo}
}

type updatePhotoHandler struct {
	repo domain.Repository
}

func (h *updatePhotoHandler) Handle(ctx context.Context, cmd UpdatePhoto) error {
	if _, err := h.repo.ByID(ctx, cmd.MotherID); err != nil {
		return err
	}
	return h.repo.UpdatePhotoURL(ctx, cmd.MotherID, cmd.PhotoURL)
}
