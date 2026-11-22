package command

import (
	"context"
	"log"
	"strings"

	"backend/internal/learn/domain"
)

type LinkUsersHandler struct {
	repo   domain.ReadRepository
	mothers domain.MotherDirectory
}

func NewLinkUsersHandler(repo domain.ReadRepository, mothers domain.MotherDirectory) *LinkUsersHandler {
	return &LinkUsersHandler{repo: repo, mothers: mothers}
}

func (h *LinkUsersHandler) LinkUsers(ctx context.Context) error {
	learners, err := h.repo.GetAllLearners(ctx)
	if err != nil {
		return err
	}

	mothers, err := h.mothers.GetAllMothers(ctx)
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
				if err := h.mothers.UpdateLearnUserID(ctx, m.Id, userId); err != nil {
					log.Printf("LinkUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("LinkUsers: mother %d (%s %s, email=%s) → learner %s (by email)", m.Id, m.FirstName, m.LastName, m.Email, userId)
				continue
			}
		}

		if motherEmail != "" {
			if userId, ok := loginMap[motherEmail]; ok {
				if err := h.mothers.UpdateLearnUserID(ctx, m.Id, userId); err != nil {
					log.Printf("LinkUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("LinkUsers: mother %d (%s %s, email=%s) → learner %s (by login)", m.Id, m.FirstName, m.LastName, m.Email, userId)
				continue
			}
		}

		motherPhone := strings.TrimSpace(strings.ToLower(m.Phone))
		if motherPhone != "" {
			if userId, ok := loginMap[motherPhone]; ok {
				if err := h.mothers.UpdateLearnUserID(ctx, m.Id, userId); err != nil {
					log.Printf("LinkUsers: ошибка обновления mother %d: %v", m.Id, err)
					continue
				}
				linked++
				log.Printf("LinkUsers: mother %d (%s %s, phone=%s) → learner %s (by login)", m.Id, m.FirstName, m.LastName, m.Phone, userId)
			}
		}
	}

	log.Printf("LinkUsers: связано %d мам с пользователями Learn", linked)
	return nil
}
