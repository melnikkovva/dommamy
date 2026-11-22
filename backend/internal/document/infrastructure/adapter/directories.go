package adapter

import (
	"context"
	"errors"

	"backend/internal/document/domain"
	"backend/internal/help"
)

type MotherChecker interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type MotherDirectory struct {
	checker MotherChecker
}

func NewMotherDirectory(checker MotherChecker) *MotherDirectory {
	return &MotherDirectory{checker: checker}
}

func (m *MotherDirectory) Exists(ctx context.Context, id int) (bool, error) {
	return m.checker.Exists(ctx, id)
}

type HelpDirectory struct {
	service *help.Service
}

func NewHelpDirectory(service *help.Service) *HelpDirectory {
	return &HelpDirectory{service: service}
}

func (a *HelpDirectory) MotherOf(ctx context.Context, helpID int) (int, error) {
	item, err := a.service.ById(ctx, helpID)
	if errors.Is(err, help.ErrorNotFound) {
		return 0, domain.ErrNoHelp
	}
	if err != nil {
		return 0, err
	}
	return item.MotherId, nil
}
