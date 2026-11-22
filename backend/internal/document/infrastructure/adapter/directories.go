package adapter

import (
	"context"
	"errors"

	"backend/internal/document/domain"
	"backend/internal/help"
	"backend/internal/mother"
)

type MotherDirectory struct {
	service *mother.Service
}

func NewMotherDirectory(service *mother.Service) *MotherDirectory {
	return &MotherDirectory{service: service}
}

func (a *MotherDirectory) Exists(_ context.Context, id int) (bool, error) {
	return a.service.Exists(id)
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
