package mother

import (
	"errors"
)

type Service struct {
	repo Repository
}

type Repository interface {
	Save(m *Mother) error
	Search(q string) ([]*Mother, error)
	ById(id int) (*Mother, error)
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(firstName string, lastName string, phone string) (*Mother, error) {
	m, err := CreateNew(firstName, lastName, phone)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) Search(q string) ([]*Mother, error) {
	return s.repo.Search(q)
}

func (s *Service) ById(id int) (*Mother, error) {
	return s.repo.ById(id)
}

func (s *Service) Exists(id int) (bool, error) {
	_, err := s.repo.ById(id)
	if errors.Is(err, ErrorNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}