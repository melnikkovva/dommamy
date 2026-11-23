package command

import (
	"context"

	motherquery "backend/internal/mother/application/query"
	"backend/internal/mother/domain"
)

type UploadPhoto struct {
	MotherID    int
	Data        []byte
	ContentType string
}

type UploadPhotoHandler interface {
	Handle(ctx context.Context, cmd UploadPhoto) (string, error)
}

func NewUploadPhotoHandler(
	repo domain.Repository,
	storage domain.PhotoStorage,
	reads motherquery.ReadModel,
) UploadPhotoHandler {
	return &uploadPhotoHandler{
		repo:    repo,
		storage: storage,
		reads:   reads,
	}
}

type uploadPhotoHandler struct {
	repo    domain.Repository
	storage domain.PhotoStorage
	reads   motherquery.ReadModel
}

func (h *uploadPhotoHandler) Handle(ctx context.Context, cmd UploadPhoto) (string, error) {
	exists, err := h.reads.Exists(ctx, cmd.MotherID)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", domain.ErrNotFound
	}

	photoURL, err := h.storage.Upload(ctx, cmd.Data, cmd.ContentType)
	if err != nil {
		return "", err
	}

	if err := h.repo.UpdatePhotoURL(ctx, cmd.MotherID, photoURL); err != nil {
		return "", err
	}
	return photoURL, nil
}
