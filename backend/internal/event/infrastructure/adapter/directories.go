package adapter

import (
	"context"

	"backend/internal/event/application/command"
)

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

func NewMotherDirectory(d MotherDirectory) command.MotherDirectory {
	return d
}
