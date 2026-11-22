package query

import (
	"context"
	"time"

	"backend/internal/document/domain"
)

type DocumentView struct {
	ID           int            `json:"id"`
	MotherID     int            `json:"mother_id"`
	Type         domain.Type    `json:"type"`
	HelpID       int            `json:"help_id"`
	Status       domain.Status  `json:"status"`
	SignedAt     string         `json:"signed_at"`
	HasFile      bool           `json:"has_file"`
	OriginalName string         `json:"original_name"`
	ContentType  string         `json:"content_type"`
	Size         int64          `json:"size"`
	CreatedAt    time.Time      `json:"created_at"`
}

type ReadModel interface {
	GetByID(ctx context.Context, id int) (DocumentView, error)
	ByMother(ctx context.Context, motherID int) ([]DocumentView, error)
}
