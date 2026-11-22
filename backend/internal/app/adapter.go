package app

import (
	"context"
	"errors"

	"backend/internal/help"
	"backend/internal/learn"
	learnmysql "backend/internal/learn/infrastructure/mysql"
	"backend/internal/mother"
)

type helpAdapter struct{ s *help.Service }

func (a helpAdapter) MotherOf(helpId int) (int, bool, error) {
	item, err := a.s.ById(context.Background(), helpId)
	if errors.Is(err, help.ErrorNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return item.MotherId, true, nil
}

type motherReaderAdapter struct{ repo mother.Repository }

func (a motherReaderAdapter) GetAll() ([]learn.MotherRef, error) {
	mothers, err := a.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]learn.MotherRef, len(mothers))
	for i, m := range mothers {
		result[i] = learn.MotherRef{Id: m.Id, FirstName: m.FirstName, LastName: m.LastName, Phone: m.Phone, Email: m.Email}
	}
	return result, nil
}

func (a motherReaderAdapter) UpdateLearnUserID(id int, learnUserID string) error {
	return a.repo.UpdateLearnUserID(id, learnUserID)
}

type learnerReaderAdapter struct{ repo *learnmysql.Repo }

func (a learnerReaderAdapter) GetAllLearners() ([]learn.Learner, error) {
	learners, err := a.repo.GetAllLearners()
	if err != nil {
		return nil, err
	}
	result := make([]learn.Learner, len(learners))
	for i, l := range learners {
		result[i] = learn.Learner{UserId: l.UserId, Login: l.Login, Email: l.Email, FirstName: l.FirstName, LastName: l.LastName}
	}
	return result, nil
}