package memory

import (
	"context"
	"backend/internal/help"
)

type Repo struct {
	nextId int
	items []*help.Item
}

func NewRepo() *Repo {
	return &Repo{ nextId: 1 }
}

func (r *Repo) Save(c context.Context, item *help.Item) error {
	item.Id = r.nextId
	r.nextId++
	r.items = append(r.items, item)
	return nil
}

func (r *Repo) List(c context.Context) ([]*help.Item, error) {
	return r.items, nil
}