package domain

import "context"

type MotherLearnLinkRepository interface {
	Link(ctx context.Context, motherID int, learnUserID string) error
	GetByMotherID(ctx context.Context, motherID int) (string, error)
	IsLinked(ctx context.Context, motherID int) (bool, error)
}
