package help

import (
	"context"
	"time"
)

type Repository interface {
	Save(c context.Context, item *Item) error
	List(c context.Context) ([]*Item, error)
	ById(ctx context.Context, id int) (*Item, error)
}

type Service struct {
	repo Repository
	mother MotherDirectory
}

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

func (s *Service) Item(c context.Context, motherId int, kind Kind, sum int, checkUrl string) (*Item, error) {
	item, err := AddItem(motherId, kind, sum, checkUrl, time.Now())
	if err != nil {
		return nil, err
	}
	ok, err := s.mother.Exists(c, motherId)
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

func NewService(repo Repository, mother MotherDirectory) *Service {
	return &Service{repo: repo, mother: mother}
}


func (s *Service) ById(ctx context.Context, id int) (*Item, error) {
	return s.repo.ById(ctx, id)
}