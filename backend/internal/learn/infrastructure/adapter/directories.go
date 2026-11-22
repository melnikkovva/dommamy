package adapter

import (
	"context"

	"backend/internal/learn/domain"
	motherdomain "backend/internal/mother/domain"
)

type MotherDirectoryAdapter struct {
	repo MotherRepo
}

type MotherRepo interface {
	GetAll(ctx context.Context) ([]motherdomain.MotherSnapshot, error)
	UpdateLearnUserID(ctx context.Context, id int, learnUserID string) error
}

func NewMotherDirectoryAdapter(repo MotherRepo) *MotherDirectoryAdapter {
	return &MotherDirectoryAdapter{repo: repo}
}

func (a *MotherDirectoryAdapter) GetAllMothers(ctx context.Context) ([]domain.MotherRef, error) {
	snapshots, err := a.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	refs := make([]domain.MotherRef, len(snapshots))
	for i, s := range snapshots {
		refs[i] = domain.MotherRef{
			Id:        s.ID,
			FirstName: s.FirstName,
			LastName:  s.LastName,
			Phone:     s.Phone,
			Email:     s.Email,
		}
	}
	return refs, nil
}

func (a *MotherDirectoryAdapter) UpdateLearnUserID(ctx context.Context, id int, learnUserID string) error {
	return a.repo.UpdateLearnUserID(ctx, id, learnUserID)
}
