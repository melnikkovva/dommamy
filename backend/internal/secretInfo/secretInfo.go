package secretInfo

import (
	"time"
)

type Data struct {
	MotherId int
	CriminalRecord string
	NdPnd string
	Drugs string
	ChronicDiseases string
	UpdatedAt time.Time
}

func AddNew(motherId int, criminal, ndPnd, drugs, chronic string, now time.Time) (*Data, error) {
	if motherId <= 0 {
		return nil, ErrorMother
	}
	return &Data{
		MotherId: motherId,
		CriminalRecord: criminal,
		NdPnd: ndPnd,
		Drugs: drugs,
		ChronicDiseases: chronic,
		UpdatedAt: now,
	}, nil
}