package command

import (
	"context"
)

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type HelpDirectory interface {
	MotherOf(ctx context.Context, helpID int) (int, error)
}
