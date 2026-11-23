package adapter

import (
	"context"

	"backend/internal/learn/domain"
	motherhttp "backend/internal/mother/interfaces/httpapi"
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

type LearnerReader struct {
	reader domain.ReadRepository
}

func NewLearnerReader(reader domain.ReadRepository) *LearnerReader {
	return &LearnerReader{reader: reader}
}

func (a *LearnerReader) GetAllLearners(ctx context.Context) ([]motherhttp.LearnerDTO, error) {
	learners, err := a.reader.GetAllLearners(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]motherhttp.LearnerDTO, len(learners))
	for i, l := range learners {
		dtos[i] = motherhttp.LearnerDTO{
			UserID:    l.UserId,
			FirstName: l.FirstName,
			LastName:  l.LastName,
			Email:     l.Email,
		}
	}
	return dtos, nil
}
