package document

import "io"

type Repository interface {
	Save(d *Document) error 
	ById(id int) (*Document, error)
	ByMother(motherId int) ([]*Document, error)
	Update(d *Document) error
}

type FileStorage interface {
	Put(name string, content io.Reader) error
	Open(name string) (io.ReadCloser, error)
	Delete(name string) error
}

type MotherChecker interface {
	Exists(id int) (bool, error)
}

type HelpChecker interface {
	MotherOf(helpId int) (motherId int, found bool, err error)
}

type Service struct {
	repo Repository
	storage FileStorage
	mothers MotherChecker
	helps HelpChecker
}

func NewService(repo Repository, storage FileStorage, mothers MotherChecker, helps HelpChecker) *Service {
	return &Service{repo: repo, storage: storage, mothers: mothers, helps: helps}
}

func (s *Service) Create(motherId int, docType string, helpId int) (*Document, error) {
	d, err := AddNew(motherId, docType, helpId)
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
	if d.Type == TypeResponsibility {
		owner, found, err := s.helps.MotherOf(helpId)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, ErrorNoHelp
		}
		if owner != motherId {
			return nil, ErrorHelpMismatch
		}
	}
	if err := s.repo.Save(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ByMother(motherId int) ([]*Document, error) {
	return s.repo.ByMother(motherId)
}

func (s *Service) AttachFile(id int, content io.Reader, size int64, contentType, originalName string) (*Document, error) {
	d, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	oldFile := d.FileName

	stored, err := NewStoredName(contentType)
	if err != nil {
		return nil, err
	}
	if err := d.AttachFile(stored, originalName, contentType, size); err != nil {
		return nil, err 
	}
	if err := s.storage.Put(stored, content); err != nil {
		return nil, err
	}
	if err := s.repo.Update(d); err != nil {
		s.storage.Delete(stored) 
	}
	if oldFile != "" {
		s.storage.Delete(oldFile) 
	}
	return d, nil
}

func (s *Service) Sign(id int, date string) (*Document, error) {
	d, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	if err := d.Sign(date); err != nil {
		return nil, err
	}
	if err := s.repo.Update(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Open(id int) (*Document, io.ReadCloser, error) {
	d, err := s.repo.ById(id)
	if err != nil {
		return nil, nil, err
	}
	if !d.HasFile() {
		return nil, nil, ErrorNoFile
	}
	f, err := s.storage.Open(d.FileName)
	if err != nil {
		return nil, nil, err
	}
	return d, f, nil
}