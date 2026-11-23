package command

import (
	"context"
	"io"

	"backend/internal/document/domain"
	"backend/internal/document/infrastructure/storage"
)

type AttachFile struct {
	DocumentID   int
	Content      io.Reader
	Size         int64
	ContentType  string
	OriginalName string
}

type AttachFileHandler interface {
	Handle(ctx context.Context, cmd AttachFile) (int, error)
}

func NewAttachFileHandler(repo domain.Repository, fileStorage domain.FileStorage) AttachFileHandler {
	return &attachFileHandler{repo: repo, fileStorage: fileStorage}
}

type attachFileHandler struct {
	repo        domain.Repository
	fileStorage domain.FileStorage
}

func (h *attachFileHandler) Handle(ctx context.Context, cmd AttachFile) (int, error) {
	d, err := h.repo.ByID(ctx, cmd.DocumentID)
	if err != nil {
		return 0, err
	}

	oldFile := d.FileName()

	stored, err := storage.GenerateStoredName(cmd.ContentType)
	if err != nil {
		return 0, err
	}

	if err := d.AttachFile(stored, cmd.OriginalName, cmd.ContentType, cmd.Size); err != nil {
		return 0, err
	}

	if err := h.fileStorage.Put(stored, cmd.Content); err != nil {
		return 0, err
	}

	if err := h.repo.Update(ctx, d); err != nil {
		h.fileStorage.Delete(stored)
		return 0, err
	}

	if oldFile != "" {
		h.fileStorage.Delete(oldFile) 
	}

	return d.ID(), nil
}
