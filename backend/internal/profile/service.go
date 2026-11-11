package profile

import "time"

type Repository interface {
	Save(p *Profile) error
	ByMother(motherId int) (*Profile, error)
	SetStatus(motherId int, status string) error
	ListByStatus(status string) ([]*Profile, error)
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

func (s *Service) Save(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}
	ok, err := s.mothers.Exists(p.MotherId)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorNoMother
	}

	old, err := s.repo.ByMother(p.MotherId)
	if err != nil {
		return err
	}
	if old != nil {
		p.Status = old.Status
	} else {
		p.Status = StatusDraft
	}
	p.UpdatedAt = time.Now()
	return s.repo.Save(p)
}

func (s *Service) Get(motherId int) (*Profile, error) {
	p, err := s.repo.ByMother(motherId)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return &Profile{MotherId: motherId, Status: StatusDraft}, nil
	}
	return p, nil
}

func (s *Service) SetStatus(motherId int, status string) error {
	if !ValidStatus(status) {
		return ErrorStatus
	}
	p, err := s.repo.ByMother(motherId)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrorNotFilled
	}
	if !CanChange(p.Status, status) {
		return ErrorTransition
	}
	if status == StatusSent {
		if err := p.Complete(); err != nil {
			return err
		}
	}
	return s.repo.SetStatus(motherId, status)
}

func (s *Service) List(status string) ([]*Profile, error) {
	return s.repo.ListByStatus(status)
}