package child

import (
	"context"
	"time"
)

type Repository interface {
	Save(c *Child) error
	SearchByMother(motherId int) ([]*Child, error)
}

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type Service struct {
	repo    Repository
	mothers MotherDirectory
}

func NewService(repo Repository, mothers MotherDirectory) *Service {
	return &Service{repo: repo, mothers: mothers}
}

func (s *Service) Add(ctx context.Context, motherId int, firstName string, birthDate time.Time) (*Child, error) {
	c, err := AddNew(motherId, firstName, birthDate, time.Now())
	if err != nil {
		return nil, err
	}
	ok, err := s.mothers.Exists(ctx, motherId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrorNoMother
	}
	if err := s.repo.Save(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) SearchByMother(motherId int) ([]*Child, error) {
	return s.repo.SearchByMother(motherId)
}
