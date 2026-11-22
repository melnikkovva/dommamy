package app

import (
	"context"
	"errors"

	"backend/internal/help"
	"backend/internal/mother/domain"
)

type helpAdapter struct{ s *help.Service }

func (a helpAdapter) MotherOf(helpId int) (int, bool, error) {
	item, err := a.s.ById(context.Background(), helpId)
	if errors.Is(err, help.ErrorNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return item.MotherId, true, nil
}

type MotherRepo interface {
	GetAll(ctx context.Context) ([]domain.MotherSnapshot, error)
	UpdateLearnUserID(ctx context.Context, id int, learnUserID string) error
}