package query

import (
	"context"
	"time"

	"backend/internal/lifecycle/domain"
)

type StatusView struct {
	MotherID  int                    `json:"mother_id"`
	Status    domain.Status          `json:"status"`
	History   []HistoryItemView      `json:"history"`
}

type HistoryItemView struct {
	Status    domain.Status `json:"status"`
	Reason    string        `json:"reason"`
	CreatedAt time.Time     `json:"created_at"`
}

type CurrentItem struct {
	MotherID int            `json:"mother_id"`
	Status   domain.Status  `json:"status"`
}

type ReadModel interface {
	GetStatus(ctx context.Context, motherID int) (StatusView, error)
	Currents(ctx context.Context) ([]CurrentItem, error)
}
