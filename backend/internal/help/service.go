package help

import (
	"context"
)

type Repository interface {
	Save(c context.Context, item *Item) error
	List(c context.Context) ([]*Item, error)
	ById(ctx context.Context, id int) (*Item, error)
}

type Service struct {
	repo Repository
	mother MotherChecker
}

type MotherChecker interface {
	Exists(id int) (bool, error)
}

func (s *Service) Item(c context.Context, motherId int, kind Kind, sum int, checkUrl string) (*Item, error) {
	item, err := AddItem(motherId, kind, sum, checkUrl)
	if err != nil {
		return nil, err
	}
	ok, err := s.mother.Exists(motherId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrorMum
	}
	if err := s.repo.Save(c, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) List(c context.Context) ([]*Item, error) {
	return s.repo.List(c)
}

func NewService(repo Repository, mother MotherChecker) *Service {
	return &Service{repo: repo, mother: mother}
}


func (s *Service) ById(ctx context.Context, id int) (*Item, error) {
	return s.repo.ById(ctx, id)
}