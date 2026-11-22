package domain

import "context"

type Repository interface {
	Save(ctx context.Context, m Mother) error
	ByID(ctx context.Context, id int) (MotherSnapshot, error)
	Search(ctx context.Context, q string) ([]MotherSnapshot, error)
	GetAll(ctx context.Context) ([]MotherSnapshot, error)
	UpdateLearnUserID(ctx context.Context, id int, learnUserID string) error
	UpdatePhotoURL(ctx context.Context, id int, photoURL string) error
}
