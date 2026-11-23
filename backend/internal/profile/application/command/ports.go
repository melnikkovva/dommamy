package command

import "context"

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}
