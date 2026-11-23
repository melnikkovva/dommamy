package adapter

import (
	"context"

	"backend/internal/learn/domain"
	motherquery "backend/internal/mother/application/query"
)

type MotherDirectoryAdapter struct {
	reader motherquery.MotherListReader
}

func NewMotherDirectoryAdapter(reader motherquery.MotherListReader) *MotherDirectoryAdapter {
	return &MotherDirectoryAdapter{reader: reader}
}

func (a *MotherDirectoryAdapter) GetAllMothers(ctx context.Context) ([]domain.MotherRef, error) {
	refs, err := a.reader.GetAllMothers(ctx)
	if err != nil {
		return nil, err
	}
	domainRefs := make([]domain.MotherRef, len(refs))
	for i, r := range refs {
		domainRefs[i] = domain.MotherRef{
			Id:        r.ID,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Phone:     r.Phone,
			Email:     r.Email,
		}
	}
	return domainRefs, nil
}
