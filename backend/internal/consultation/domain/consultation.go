package domain

import (
	"errors"
	"time"
)

var (
	ErrorDate    = errors.New("дата консультации обязательна")
	ErrorType    = errors.New("тип специалиста обязателен")
	ErrorUnknown = errors.New("неизвестный тип специалиста")
)

type SpecialistType string

const (
	Psychologist  SpecialistType = "psychologist"
	Lawyer        SpecialistType = "lawyer"
	Priest        SpecialistType = "priest"
	Employment    SpecialistType = "employment"
)

var validSpecialistTypes = map[SpecialistType]bool{
	Psychologist: true,
	Lawyer:       true,
	Priest:       true,
	Employment:   true,
}

func (s SpecialistType) Valid() bool {
	return validSpecialistTypes[s]
}

type Status string

const (
	Scheduled   Status = "scheduled"
	Completed   Status = "completed"
	Cancelled   Status = "cancelled"
)

type Consultation struct {
	Id              int
	MotherId        int
	SpecialistType  SpecialistType
	ConsultationDate time.Time
	Notes           string
	Status          Status
	CreatedAt       time.Time
}

type ConsultationInput struct {
	MotherId        int
	SpecialistType  SpecialistType
	ConsultationDate time.Time
	Notes           string
	Status          Status
}

func NewConsultation(input ConsultationInput) (*Consultation, error) {
	if input.MotherId <= 0 {
		return nil, ErrorMum
	}
	if !input.SpecialistType.Valid() {
		return nil, ErrorUnknown
	}
	if input.ConsultationDate.IsZero() {
		return nil, ErrorDate
	}
	if input.Status == "" {
		input.Status = Completed
	}
	return &Consultation{
		MotherId:        input.MotherId,
		SpecialistType:  input.SpecialistType,
		ConsultationDate: input.ConsultationDate,
		Notes:           input.Notes,
		Status:          input.Status,
		CreatedAt:       time.Now(),
	}, nil
}

var ErrorMum = errors.New("id мамы обязателен")

type Row struct {
	Id              int
	MotherId        int
	MotherName      string
	SpecialistType  SpecialistType
	ConsultationDate time.Time
	Notes           string
	Status          Status
	CreatedAt       time.Time
}
