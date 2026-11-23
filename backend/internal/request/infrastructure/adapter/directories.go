package adapter

import (
	"context"
	"errors"

	"backend/internal/help"
	"backend/internal/request/application/command"
	"backend/internal/request/domain"
)

func NewHelpDirectory(s *help.Service) command.HelpDirectory {
	return &helpDirectory{s: s}
}

type helpDirectory struct{ s *help.Service }

func (d *helpDirectory) MotherOf(ctx context.Context, helpID int) (int, error) {
	item, err := d.s.ById(ctx, helpID)
	if errors.Is(err, help.ErrorNotFound) {
		return 0, domain.ErrHelpNotFound
	}
	if err != nil {
		return 0, err
	}
	if item == nil { 
		return 0, domain.ErrHelpNotFound
	}
	return item.MotherId, nil
}