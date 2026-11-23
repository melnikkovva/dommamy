package secretInfo

import (
	"context"
	"time"
)

type Repository interface {
	Save(d *Data) error
	SearchByMother(motherId int) (*Data, error)
}

type MotherDirectory interface {
	Exists(ctx context.Context, id int) (bool, error)
}

type Auditor interface {
	Log(action string, motherId int) error
}

type Service struct {
	repo Repository
	mothers MotherDirectory
	audit Auditor
}

func NewService(repo Repository, mothers MotherDirectory, audit Auditor) *Service {
	return &Service{repo: repo, mothers: mothers, audit: audit}
}

func (s *Service) Save(ctx context.Context, motherId int, criminal, ndPnd, drugs, chronic string) error {
	d, err := AddNew(motherId, criminal, ndPnd, drugs, chronic, time.Now())
	if err != nil {
		return err
	}
	ok, err := s.mothers.Exists(ctx, motherId)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorNoMother
	}
	if err := s.repo.Save(d); err != nil {
		return err
	}
	return s.audit.Log("edit_secret_info", motherId)
}

func (s *Service) Reveal(motherId int) (*Data, error) {
	if err := s.audit.Log("reveal_secret_info", motherId); err != nil {
		return nil, err
	}
	return s.repo.SearchByMother(motherId)
}