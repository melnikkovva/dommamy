package child

import (
	"time"
)

type Repository interface {
	Save(c *Child) error
	SearchByMother(motherId int) ([]*Child, error)
}

type MotherChecker interface {
	Exists(id int) (bool, error)
}

type Service struct {
	repo Repository
	mothers MotherChecker
}

func NewService(repo Repository, mothers MotherChecker) *Service {
	return &Service{repo: repo, mothers: mothers}
}

func (s *Service) Add(motherId int, firstName string, birthDate time.Time) (*Child, error) {
	c, err := AddNew(motherId, firstName, birthDate)
	if err != nil {
		return nil, err
	}
	ok, err := s.mothers.Exists(motherId)
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
