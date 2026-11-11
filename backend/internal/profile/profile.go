package profile

import (
	"errors"
	"strings"
	"time"
)

var ErrorMother = errors.New("the mother is not specified")
var ErrorNoMother = errors.New("there is no such mother")
var ErrorBirth = errors.New("the date of birth is incorrect")
var ErrorIncome = errors.New("the income cannot be negative")
var ErrorWeeks = errors.New("pregnancy weeks must be from 0 to 42")
var ErrorStatus = errors.New("unknown status")
var ErrorTransition = errors.New("this status change is not allowed")
var ErrorNotFilled = errors.New("the profile is not filled yet")
var ErrorNotComplete = errors.New("fill the date of birth and the situation before sending")

const (
	StatusDraft = "draft"
	StatusSent = "sent"
	StatusChecked = "checked"
)

type Profile struct {
	MotherId int
	BirthDate string 
	BirthPlace string
	Citizenship string
	Registration string
	Living string
	Email string
	Religion string
	Education string
	Profession string
	IncomeSource string
	IncomeAmount int
	MaritalStatus string
	FatherInfo string
	PregnancyWeeks int
	Situation string
	ExpectedHelp string
	Source string
	Status string
	UpdatedAt time.Time
}

var transitions = map[string][]string{
	StatusDraft: {StatusSent},
	StatusSent: {StatusChecked, StatusDraft},
	StatusChecked: {StatusSent},
}

func CanChange(from, to string) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func ValidStatus(s string) bool {
	return s == StatusDraft || s == StatusSent || s == StatusChecked
}

func (p *Profile) Validate() error {
	if p.MotherId <= 0 {
		return ErrorMother
	}
	p.BirthDate = strings.TrimSpace(p.BirthDate)
	if p.BirthDate != "" {
		d, err := time.Parse("2006-01-02", p.BirthDate)
		if err != nil || d.After(time.Now()) {
			return ErrorBirth
		}
	}
	if p.IncomeAmount < 0 {
		return ErrorIncome
	}
	if p.PregnancyWeeks < 0 || p.PregnancyWeeks > 42 {
		return ErrorWeeks
	}
	p.Situation = strings.TrimSpace(p.Situation)
	p.ExpectedHelp = strings.TrimSpace(p.ExpectedHelp)
	return nil
}

func (p *Profile) Complete() error {
	if p.BirthDate == "" || p.Situation == "" {
		return ErrorNotComplete
	}
	return nil
}