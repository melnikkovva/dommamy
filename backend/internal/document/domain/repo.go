package domain

import "context"

type Repository interface {
	Save(ctx context.Context, d Document) error
	ByID(ctx context.Context, id int) (Document, error)
	ByMother(ctx context.Context, motherID int) ([]Document, error)
	Update(ctx context.Context, d Document) error
}
