package command

import (
	"context"
	"log"

	"backend/internal/learn/domain"
)

type SyncHandler struct {
	repo domain.WriteRepository
	api  domain.APIClient
}

func NewSyncHandler(repo domain.WriteRepository, api domain.APIClient) *SyncHandler {
	return &SyncHandler{repo: repo, api: api}
}

func (h *SyncHandler) SyncAll(ctx context.Context) error {
	log.Println("Начало синхронизации с iSpring...")

	learners, err := h.api.FetchAllLearnAccounts(ctx)
	if err != nil {
		return err
	}
	if err := h.repo.SaveLearners(ctx, learners); err != nil {
		return err
	}
	log.Printf("Синхронизировано пользователей: %d", len(learners))

	results, err := h.api.FetchAllCourseResults(ctx)
	if err != nil {
		return err
	}
	if err := h.repo.SaveCourseResults(ctx, results); err != nil {
		return err
	}
	log.Printf("Синхронизировано результатов обучения: %d", len(results))

	iprPlans, err := h.api.FetchAllIPR(ctx)
	if err != nil {
		log.Printf("FetchAllIPR ошибка: %v", err)
	} else {
		if err := h.repo.SaveIPRPlans(ctx, iprPlans); err != nil {
			log.Printf("SaveIPRPlans ошибка: %v", err)
		} else {
			log.Printf("Синхронизировано ИПР: %d", len(iprPlans))
		}
	}

	log.Println("Синхронизация успешно завершена")
	return nil
}
