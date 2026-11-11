package event

import "time"

type Repository interface {
	Types() ([]*Type, error)
	SaveEvent(e *Event) error
	Events() ([]*Event, error)
	EventById(id int) (*Event, error) 
	EventsByMother(motherId int) ([]*Event, error)

	Participation(eventId, motherId int) (*Participation, error) 
	CountParticipants(eventId int) (int, error)
	AddParticipant(eventId, motherId int) error
	SetStatus(eventId, motherId int, status string) error
	Participants(eventId int) ([]*Participation, error)
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

func (s *Service) Types() ([]*Type, error) { return s.repo.Types() }

func (s *Service) Create(typeId int, title string, startsAt time.Time, place string, maxPeople int) (*Event, error) {
	e, err := AddNew(typeId, title, startsAt, place, maxPeople)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveEvent(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) List() ([]*Event, error) { return s.repo.Events() }

func (s *Service) ByMother(motherId int) ([]*Event, error) {
	return s.repo.EventsByMother(motherId)
}

func (s *Service) Register(eventId, motherId int) error {
	e, err := s.repo.EventById(eventId)
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
	p, err := s.repo.Participation(eventId, motherId)
	if err != nil {
		return err
	}
	if p != nil {
		return ErrorAlreadyJoined
	}
	if e.MaxPeople > 0 {
		n, err := s.repo.CountParticipants(eventId)
		if err != nil {
			return err
		}
		if n >= e.MaxPeople {
			return ErrorFull
		}
	}
	return s.repo.AddParticipant(eventId, motherId)
}

func (s *Service) Mark(eventId, motherId int, status string) error {
	p, err := s.repo.Participation(eventId, motherId)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrorNotRegistered
	}
	return s.repo.SetStatus(eventId, motherId, status)
}

func (s *Service) Participants(eventId int) ([]*Participation, error) {
	return s.repo.Participants(eventId)
}