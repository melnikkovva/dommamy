package command

import (
	"context"
	"log"
	"strings"

	"backend/internal/learn/domain"
)

type LinkUsers interface {
	LinkUsers(ctx context.Context) error
}

type linkUsersHandler struct {
	repo    domain.ReadRepository
	mothers domain.MotherDirectory
	links   domain.MotherLearnLinkRepository
}

func NewLinkUsersHandler(repo domain.ReadRepository, mothers domain.MotherDirectory, links domain.MotherLearnLinkRepository) LinkUsers {
	return &linkUsersHandler{repo: repo, mothers: mothers, links: links}
}

func (h *linkUsersHandler) LinkUsers(ctx context.Context) error {
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

	linked := 0
	for _, m := range mothers {
		alreadyLinked, err := h.links.IsLinked(ctx, m.Id)
		if err != nil {
			log.Printf("LinkUsers: ошибка проверки связи mother %d: %v", m.Id, err)
			continue
		}
		if alreadyLinked {
			continue
		}

		motherEmail := strings.TrimSpace(strings.ToLower(m.Email))
		if motherEmail == "" {
			continue
		}

		userId, ok := emailMap[motherEmail]
		if !ok {
			continue
		}

		if err := h.links.Link(ctx, m.Id, userId); err != nil {
			log.Printf("LinkUsers: ошибка создания связи mother %d: %v", m.Id, err)
			continue
		}
		linked++
		log.Printf("LinkUsers: mother %d связана с пользователем Learn", m.Id)
	}

	log.Printf("LinkUsers: создано %d новых связей", linked)
	return nil
}
