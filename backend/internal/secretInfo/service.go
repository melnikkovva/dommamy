package secretInfo

type Repository interface {
	Save(d *Data) error
	SearchByMother(motherId int) (*Data, error)
}

type MotherChecker interface {
	Exists(id int) (bool, error)
}

type Auditor interface {
	Log(action string, motherId int) error
}

type Service struct {
	repo Repository
	mothers MotherChecker
	audit Auditor
}

func NewService(repo Repository, mothers MotherChecker, audit Auditor) *Service {
	return &Service{repo: repo, mothers: mothers, audit: audit}
}

func (s *Service) Save(motherId int, criminal, ndPnd, drugs, chronic string) error {
	d, err := AddNew(motherId, criminal, ndPnd, drugs, chronic)
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