package memory

import (
	"strings"
	"backend/internal/mother"
)

type Repo struct {
	items  []*mother.Mother
	nextId int
}

func NewRepo() *Repo {
	return &Repo{nextId: 1}
}

func (r *Repo) Save(m *mother.Mother) error {
	m.Id = r.nextId
	r.nextId++
	r.items = append(r.items, m)
	return nil
}

func (r *Repo) Search(q string) ([]*mother.Mother, error) {
	q = strings.ToLower(q)
	result := []*mother.Mother{}
	for _, m := range r.items {
		name := strings.ToLower(m.FirstName + " " + m.LastName)
		if strings.Contains(name, q) {
			result = append(result, m)
		}
	}
	return result, nil
}

func (r *Repo) ById(id int) (*mother.Mother, error) {
	for _, m := range r.items {
		if m.Id == id {
			return m, nil
		}
	}
	return nil, mother.ErrorNotFound
}

func (r *Repo) GetAll() ([]*mother.Mother, error) {
	return r.items, nil
}

func (r *Repo) UpdateLearnUserID(id int, learnUserID string) error {
	for _, m := range r.items {
		if m.Id == id {
			m.LearnUserID = learnUserID
			return nil
		}
	}
	return mother.ErrorNotFound
}