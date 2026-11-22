package adapter

import (
	"context"

	"backend/internal/event/application/command"
	"backend/internal/mother"
)

func NewMotherDirectory(s *mother.Service) command.MotherDirectory {
	return &motherDirectory{s: s}
}

type motherDirectory struct{ s *mother.Service }

func (d *motherDirectory) Exists(_ context.Context, id int) (bool, error) {
	ok, err := d.s.Exists(id)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}
