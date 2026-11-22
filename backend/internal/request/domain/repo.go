package domain

import "context"

type Repository interface {
	Add(ctx context.Context, r Request) (id int, err error)
	ByID(ctx context.Context, id int) (Request, error)
	ByIDForUpdate(ctx context.Context, id int) (Request, error)
	Update(ctx context.Context, r Request) error
}