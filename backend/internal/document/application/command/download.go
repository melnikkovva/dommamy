package command

import (
	"context"
	"io"

	"backend/internal/document/domain"
)

type Download struct {
	DocumentID int
}

type DownloadHandler interface {
	Handle(ctx context.Context, cmd Download) (DocumentFile, error)
}

type DocumentFile struct {
	FileName     string
	OriginalName string
	ContentType  string
	Reader       io.ReadCloser
}

func NewDownloadHandler(repo domain.Repository, storage domain.FileStorage) DownloadHandler {
	return &downloadHandler{repo: repo, storage: storage}
}

type downloadHandler struct {
	repo    domain.Repository
	storage domain.FileStorage
}

func (h *downloadHandler) Handle(ctx context.Context, cmd Download) (DocumentFile, error) {
	d, err := h.repo.ByID(ctx, cmd.DocumentID)
	if err != nil {
		return DocumentFile{}, err
	}

	if !d.HasFile() {
		return DocumentFile{}, domain.ErrNoFile
	}

	f, err := h.storage.Open(d.FileName())
	if err != nil {
		return DocumentFile{}, err
	}

	return DocumentFile{
		FileName:     d.FileName(),
		OriginalName: d.OriginalName(),
		ContentType:  d.ContentType(),
		Reader:       f,
	}, nil
}
