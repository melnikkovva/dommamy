package domain

import "context"

type Repository interface {
	Add(ctx context.Context, c Change) (int, error)
	History(ctx context.Context, motherID int) ([]ChangeSnapshot, error)
	Currents(ctx context.Context) ([]ChangeSnapshot, error)
}
