package secretInfo

import (
	"errors"
	"time"
)

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorNotFound = errors.New("secret info not filled")

type Data struct {
	MotherId int
	CriminalRecord string
	NdPnd string
	Drugs string
	ChronicDiseases string
	UpdatedAt time.Time
}

func AddNew(motherId int, criminal, ndPnd, drugs, chronic string) (*Data, error) {
	if motherId <= 0 {
		return nil, ErrorMother
	}
	return &Data{
		MotherId: motherId,
		CriminalRecord: criminal,
		NdPnd: ndPnd,
		Drugs: drugs,
		ChronicDiseases: chronic,
		UpdatedAt: time.Now(),
	}, nil
}