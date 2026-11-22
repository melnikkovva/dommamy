package command

import (
	"context"

	"backend/internal/request/domain"
)

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, repo domain.Repository) error) error
}

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type HelpDirectory interface {
	MotherOf(ctx context.Context, helpID int) (int, error)
}