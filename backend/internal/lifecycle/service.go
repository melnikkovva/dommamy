package lifecycle

import "errors"

type Repository interface {
	Add(c *Change) error
	Current(motherId int) (*Change, error)  
	History(motherId int) ([]*Change, error) 
	Currents() ([]*Change, error)  
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

func (s *Service) Status(motherId int) (string, error) {
	c, err := s.repo.Current(motherId)
	if err != nil {
		if errors.Is(err, ErrorNoMother) {
			return StatusReview, nil
		}
		return "", err
	}
	return c.Status, nil
}

func (s *Service) Change(motherId int, to, reason string) error {
	c, err := NewChange(motherId, to, reason)
	if err != nil {
		return err
	}
	ok, err := s.mothers.Exists(motherId)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorNoMother
	}
	current, err := s.Status(motherId)
	if err != nil {
		return err
	}
	if !CanChange(current, to) {
		return ErrorTransition
	}
	return s.repo.Add(c)
}

func (s *Service) History(motherId int) ([]*Change, error) {
	return s.repo.History(motherId)
}

func (s *Service) Currents() ([]*Change, error) {
	return s.repo.Currents()
}