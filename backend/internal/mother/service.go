package mother

import (
	"errors"

	"backend/internal/learn"
)

type LearnStatsRepository interface {
	GetCourseStats(learnUserID string) (*learn.CourseStats, error)
	GetCourseResults(learnUserID string) ([]learn.CourseResult, error)
}

type Service struct {
	repo Repository
	learnStatsRepo LearnStatsRepository
}

type Repository interface {
	Save(m *Mother) error
	Search(q string) ([]*Mother, error)
	ById(id int) (*Mother, error)
	GetAll() ([]*Mother, error)
	UpdateLearnUserID(id int, learnUserID string) error
	UpdatePhotoURL(id int, photoURL string) error
}

func NewService(repo Repository, learnStatsRepo LearnStatsRepository) *Service {
	return &Service{repo: repo, learnStatsRepo: learnStatsRepo}
}

func (s *Service) Create(firstName string, lastName string, phone string, email string) (*Mother, error) {
	m, err := CreateNew(firstName, lastName, phone, email)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ById(id int) (*Mother, error) {
	m, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}
	if m.LearnUserID != "" && s.learnStatsRepo != nil {
		stats, err := s.learnStatsRepo.GetCourseStats(m.LearnUserID)
		if err != nil {
			return nil, err
		}
		m.CourseStats = stats
	}
	return m, nil
}

func (s *Service) Search(q string) ([]*Mother, error) {
	mothers, err := s.repo.Search(q)
	if err != nil {
		return nil, err
	}
	for _, m := range mothers {
		if m.LearnUserID != "" && s.learnStatsRepo != nil {
			stats, err := s.learnStatsRepo.GetCourseStats(m.LearnUserID)
			if err == nil {
				m.CourseStats = stats
			}
		}
	}
	return mothers, nil
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

func (s *Service) GetCourses(motherId int) ([]learn.CourseResult, error) {
	m, err := s.repo.ById(motherId)
	if err != nil {
		return nil, err
	}
	if m.LearnUserID == "" || s.learnStatsRepo == nil {
		return []learn.CourseResult{}, nil
	}
	return s.learnStatsRepo.GetCourseResults(m.LearnUserID)
}

func (s *Service) UpdatePhotoURL(id int, photoURL string) error {
	_, err := s.repo.ById(id)
	if err != nil {
		return err
	}
	return s.repo.UpdatePhotoURL(id, photoURL)
}