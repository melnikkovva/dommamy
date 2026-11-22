package learn

import (
	"log"
	"strings"
)

type Repository interface {
	SaveLearners(learners []MotherLearnAccount) error
	SaveCourseResults(results []CourseResult) error
	SaveIPRPlans(plans []IPRPlan) error
}

type LearnerReader interface {
	GetAllLearners() ([]Learner, error)
}

type MotherReader interface {
	GetAll() ([]MotherRef, error)
	UpdateLearnUserID(id int, learnUserID string) error
}

type MotherRef struct {
	Id        int
	FirstName string
	LastName  string
	Phone     string
	Email     string
}

type SyncService struct {
	repo          Repository
	api           LearnApiClient
	learnerReader LearnerReader
	motherReader  MotherReader
}

func NewSyncService(repo Repository, api LearnApiClient) *SyncService {
	return &SyncService{repo: repo, api: api}
}

func (s *SyncService) SetLinker(learnerReader LearnerReader, motherReader MotherReader) {
	s.learnerReader = learnerReader
	s.motherReader = motherReader
}

func (s *SyncService) SyncAll() error {
	log.Println("Начало синхронизации с iSpring...")

	learners, err := s.api.FetchAllLearnAccounts()
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

	iprPlans, err := s.api.FetchAllIPR()
	if err != nil {
		log.Printf("FetchAllIPR ошибка: %v", err)
	} else {
		if err := s.repo.SaveIPRPlans(iprPlans); err != nil {
			log.Printf("SaveIPRPlans ошибка: %v", err)
		} else {
			log.Printf("Синхронизировано ИПР: %d", len(iprPlans))
		}
	}

	log.Println("Синхронизация успешно завершена")
	return nil
}

func (s *SyncService) SyncUsers() error {
	if s.learnerReader == nil || s.motherReader == nil {
		log.Println("SyncUsers: репозитории не настроены, пропускаем")
		return nil
	}

	learners, err := s.learnerReader.GetAllLearners()
	if err != nil {
		return err
	}

	mothers, err := s.motherReader.GetAll()
	if err != nil {
		return err
	}

	emailMap := make(map[string]string)
	for _, l := range learners {
		email := strings.TrimSpace(strings.ToLower(l.Email))
		if email != "" {
			emailMap[email] = l.UserId
		}
	}

	loginMap := make(map[string]string)
	for _, l := range learners {
		login := strings.TrimSpace(strings.ToLower(l.Login))
		if login != "" {
			loginMap[login] = l.UserId
		}
	}

	linked := 0
	for _, m := range mothers {
		motherEmail := strings.TrimSpace(strings.ToLower(m.Email))
		if motherEmail != "" {
			if userId, ok := emailMap[motherEmail]; ok {
				if err := s.motherReader.UpdateLearnUserID(m.Id, userId); err != nil {
					log.Printf("SyncUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("SyncUsers: mother %d (%s %s, email=%s) → learner %s (by email)", m.Id, m.FirstName, m.LastName, m.Email, userId)
				continue
			}
		}

		if motherEmail != "" {
			if userId, ok := loginMap[motherEmail]; ok {
				if err := s.motherReader.UpdateLearnUserID(m.Id, userId); err != nil {
					log.Printf("SyncUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("SyncUsers: mother %d (%s %s, email=%s) → learner %s (by login)", m.Id, m.FirstName, m.LastName, m.Email, userId)
				continue
			}
		}

		motherPhone := strings.TrimSpace(strings.ToLower(m.Phone))
		if motherPhone != "" {
			if userId, ok := loginMap[motherPhone]; ok {
				if err := s.motherReader.UpdateLearnUserID(m.Id, userId); err != nil {
					log.Printf("SyncUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("SyncUsers: mother %d (%s %s, phone=%s) → learner %s (by login)", m.Id, m.FirstName, m.LastName, m.Phone, userId)
			}
		}
	}

	log.Printf("SyncUsers: связано %d мам с пользователями Learn", linked)
	return nil
}
