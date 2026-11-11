package request

type Repository interface {
	Save(r *Request) error
	ById(id int) (*Request, error) 
	List(status string) ([]*Request, error) 
	ByMother(motherId int) ([]*Request, error)
	Update(r *Request) error
}

type MotherChecker interface {
	Exists(id int) (bool, error)
}

type HelpChecker interface {
	MotherOf(helpId int) (motherId int, found bool, err error)
}

type Service struct {
	repo    Repository
	mothers MotherChecker
	helps   HelpChecker
}

func NewService(repo Repository, mothers MotherChecker, helps HelpChecker) *Service {
	return &Service{repo: repo, mothers: mothers, helps: helps}
}

func (s *Service) Create(motherId int, kind, description, source string) (*Request, error) {
	r, err := AddNew(motherId, kind, description, source)
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
	if err := s.repo.Save(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) SetStatus(id int, to, comment string) (*Request, error) {
	r, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	if err := r.ChangeStatus(to, comment); err != nil {
		return nil, err
	}
	if err := s.repo.Update(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) Fulfill(id, helpId int) (*Request, error) {
	r, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	motherId, found, err := s.helps.MotherOf(helpId)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrorNoHelp
	}
	if motherId != r.MotherId {
		return nil, ErrorHelpMismatch
	}
	if err := r.Fulfill(helpId); err != nil {
		return nil, err
	}
	if err := s.repo.Update(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) List(status string) ([]*Request, error) { return s.repo.List(status) }

func (s *Service) ByMother(motherId int) ([]*Request, error) { return s.repo.ByMother(motherId) }