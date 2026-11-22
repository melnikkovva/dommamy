package query

import (
	"context"
	"time"
)

type RequestView struct {
	ID          int
	MotherID    int
	Kind        string
	Description string
	Status      string
	Comment     string
	Source      string
	HelpID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReadModel interface {
	Get(ctx context.Context, id int) (RequestView, error) 
	List(ctx context.Context, status string) ([]RequestView, error)
	ByMother(ctx context.Context, motherID int) ([]RequestView, error)
}