package adapter

import (
	"context"

	"backend/internal/mother"
)

type MotherDirectory struct {
	service *mother.Service
}

func NewMotherDirectory(service *mother.Service) *MotherDirectory {
	return &MotherDirectory{service: service}
}

func (a *MotherDirectory) Exists(ctx context.Context, id int) (bool, error) {
	return a.service.Exists(id)
}
