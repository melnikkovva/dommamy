package domain

import "context"

type Repository interface {
	Save(ctx context.Context, p Profile) error
	ByMother(ctx context.Context, motherID int) (Profile, error)
	UpdateStatus(ctx context.Context, motherID int, status Status) error
}
