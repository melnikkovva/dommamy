package domain

import "context"

type Repository interface {
	Add(ctx context.Context, m Mother) (int, error)
	UpdatePhotoURL(ctx context.Context, id int, photoURL string) error
	SoftDelete(ctx context.Context, id int) error
}
