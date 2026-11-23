package domain

import (
	"slices"
	"strings"
	"time"
)

type Status string

const (
	StatusDraft   Status = "draft"
	StatusSent    Status = "sent"
	StatusChecked Status = "checked"
)

var validStatuses = []Status{StatusDraft, StatusSent, StatusChecked}

func (s Status) Valid() bool {
	return slices.Contains(validStatuses, s)
}

var transitions = map[Status][]Status{
	StatusDraft:   {StatusSent},
	StatusSent:    {StatusChecked, StatusDraft},
	StatusChecked: {StatusSent},
}

func CanTransition(from, to Status) bool {
	for _, next := range transitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

type Profile struct {
	motherID       int
	birthDate      string
	birthPlace     string
	citizenship    string
	registration   string
	living         string
	email          string
	religion       string
	education      string
	profession     string
	incomeSource   string
	incomeAmount   int
	maritalStatus  string
	fatherInfo     string
	pregnancyWeeks int
	situation      string
	expectedHelp   string
	source         string
	status         Status
	updatedAt      time.Time
}

func NewProfile(motherID int) (Profile, error) {
	if motherID <= 0 {
		return Profile{}, ErrMother
	}
	return Profile{
		motherID: motherID,
		status:   StatusDraft,
	}, nil
}

type ProfileInput struct {
	BirthDate      string
	BirthPlace     string
	Citizenship    string
	Registration   string
	Living         string
	Email          string
	Religion       string
	Education      string
	Profession     string
	IncomeSource   string
	IncomeAmount   int
	MaritalStatus  string
	FatherInfo     string
	PregnancyWeeks int
	Situation      string
	ExpectedHelp   string
	Source         string
}

func (p *Profile) Fill(input ProfileInput, now time.Time) error {
	p.birthDate = strings.TrimSpace(input.BirthDate)
	if p.birthDate != "" {
		d, err := time.Parse("2006-01-02", p.birthDate)
		if err != nil || d.After(now) {
			return ErrBirth
		}
	}

	p.birthPlace = strings.TrimSpace(input.BirthPlace)
	p.citizenship = strings.TrimSpace(input.Citizenship)
	p.registration = strings.TrimSpace(input.Registration)
	p.living = strings.TrimSpace(input.Living)
	p.email = strings.TrimSpace(input.Email)
	p.religion = strings.TrimSpace(input.Religion)
	p.education = strings.TrimSpace(input.Education)
	p.profession = strings.TrimSpace(input.Profession)
	p.incomeSource = strings.TrimSpace(input.IncomeSource)
	p.incomeAmount = input.IncomeAmount
	if p.incomeAmount < 0 {
		return ErrIncome
	}
	p.maritalStatus = strings.TrimSpace(input.MaritalStatus)
	p.fatherInfo = strings.TrimSpace(input.FatherInfo)
	p.pregnancyWeeks = input.PregnancyWeeks
	if p.pregnancyWeeks < 0 || p.pregnancyWeeks > 42 {
		return ErrWeeks
	}
	p.situation = strings.TrimSpace(input.Situation)
	p.expectedHelp = strings.TrimSpace(input.ExpectedHelp)
	p.source = strings.TrimSpace(input.Source)
	p.updatedAt = now

	return nil
}

func (p *Profile) Send() error {
	if p.birthDate == "" || p.situation == "" {
		return ErrNotComplete
	}
	if !CanTransition(p.status, StatusSent) {
		return ErrTransition
	}
	p.status = StatusSent
	return nil
}

func (p *Profile) Transition(to Status) error {
	if !to.Valid() {
		return ErrStatus
	}
	if !CanTransition(p.status, to) {
		return ErrTransition
	}
	p.status = to
	return nil
}

func (p Profile) MotherID() int            { return p.motherID }
func (p Profile) BirthDate() string        { return p.birthDate }
func (p Profile) BirthPlace() string       { return p.birthPlace }
func (p Profile) Citizenship() string      { return p.citizenship }
func (p Profile) Registration() string     { return p.registration }
func (p Profile) Living() string           { return p.living }
func (p Profile) Email() string            { return p.email }
func (p Profile) Religion() string         { return p.religion }
func (p Profile) Education() string        { return p.education }
func (p Profile) Profession() string       { return p.profession }
func (p Profile) IncomeSource() string     { return p.incomeSource }
func (p Profile) IncomeAmount() int        { return p.incomeAmount }
func (p Profile) MaritalStatus() string    { return p.maritalStatus }
func (p Profile) FatherInfo() string       { return p.fatherInfo }
func (p Profile) PregnancyWeeks() int      { return p.pregnancyWeeks }
func (p Profile) Situation() string        { return p.situation }
func (p Profile) ExpectedHelp() string     { return p.expectedHelp }
func (p Profile) Source() string           { return p.source }
func (p Profile) Status() Status           { return p.status }
func (p Profile) UpdatedAt() time.Time { return p.updatedAt }

type Snapshot struct {
	MotherID       int
	BirthDate      string
	BirthPlace     string
	Citizenship    string
	Registration   string
	Living         string
	Email          string
	Religion       string
	Education      string
	Profession     string
	IncomeSource   string
	IncomeAmount   int
	MaritalStatus  string
	FatherInfo     string
	PregnancyWeeks int
	Situation      string
	ExpectedHelp   string
	Source         string
	Status         Status
	UpdatedAt      time.Time
}

func Restore(s Snapshot) Profile {
	return Profile{
		motherID:       s.MotherID,
		birthDate:      s.BirthDate,
		birthPlace:     s.BirthPlace,
		citizenship:    s.Citizenship,
		registration:   s.Registration,
		living:         s.Living,
		email:          s.Email,
		religion:       s.Religion,
		education:      s.Education,
		profession:     s.Profession,
		incomeSource:   s.IncomeSource,
		incomeAmount:   s.IncomeAmount,
		maritalStatus:  s.MaritalStatus,
		fatherInfo:     s.FatherInfo,
		pregnancyWeeks: s.PregnancyWeeks,
		situation:      s.Situation,
		expectedHelp:   s.ExpectedHelp,
		source:         s.Source,
		status:         s.Status,
		updatedAt:      s.UpdatedAt,
	}
}

func (p Profile) Snapshot() Snapshot {
	return Snapshot{
		MotherID:       p.motherID,
		BirthDate:      p.birthDate,
		BirthPlace:     p.birthPlace,
		Citizenship:    p.citizenship,
		Registration:   p.registration,
		Living:         p.living,
		Email:          p.email,
		Religion:       p.religion,
		Education:      p.education,
		Profession:     p.profession,
		IncomeSource:   p.incomeSource,
		IncomeAmount:   p.incomeAmount,
		MaritalStatus:  p.maritalStatus,
		FatherInfo:     p.fatherInfo,
		PregnancyWeeks: p.pregnancyWeeks,
		Situation:      p.situation,
		ExpectedHelp:   p.expectedHelp,
		Source:         p.source,
		Status:         p.status,
		UpdatedAt:      p.updatedAt,
	}
}
