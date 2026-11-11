package lms

import "log"

type SyncService struct {
	repo LearnerRepository
	api LearnApiClient
}

func NewSyncService(repo LearnerRepository, api LearnApiClient) *SyncService {
	return &SyncService{repo: repo, api: api}
}

func (s *SyncService) SyncAll() error {
	log.Println("Начало синхронизации с iSpring...")

	learners, err := s.api.FetchAllLearners()
	if err != nil {
		return err
	}
	if err := s.repo.SaveLearners(learners); err != nil {
		return err
	}
	log.Printf("Синхронизировано пользователей: %d", len(learners))

	results, err := s.api.FetchAllCourseResults()
	if err != nil {
		return err
	}
	if err := s.repo.SaveCourseResults(results); err != nil {
		return err
	}
	log.Printf("Синхронизировано результатов обучения: %d", len(results))

	log.Println("Синхронизация успешно завершена")
	return nil
}