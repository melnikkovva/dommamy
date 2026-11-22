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

func (p *Profile) SetID(id int) {
}

func (p *Profile) Fill(birthDate, birthPlace, citizenship, registration, living, email,
	religion, education, profession, incomeSource, maritalStatus, fatherInfo, source string,
	incomeAmount, pregnancyWeeks int, situation, expectedHelp string) error {

	p.birthDate = strings.TrimSpace(birthDate)
	if p.birthDate != "" {
		d, err := time.Parse("2006-01-02", p.birthDate)
		if err != nil || d.After(time.Now()) {
			return ErrBirth
		}
	}

	p.birthPlace = strings.TrimSpace(birthPlace)
	p.citizenship = strings.TrimSpace(citizenship)
	p.registration = strings.TrimSpace(registration)
	p.living = strings.TrimSpace(living)
	p.email = strings.TrimSpace(email)
	p.religion = strings.TrimSpace(religion)
	p.education = strings.TrimSpace(education)
	p.profession = strings.TrimSpace(profession)
	p.incomeSource = strings.TrimSpace(incomeSource)
	p.incomeAmount = incomeAmount
	if p.incomeAmount < 0 {
		return ErrIncome
	}
	p.maritalStatus = strings.TrimSpace(maritalStatus)
	p.fatherInfo = strings.TrimSpace(fatherInfo)
	p.pregnancyWeeks = pregnancyWeeks
	if p.pregnancyWeeks < 0 || p.pregnancyWeeks > 42 {
		return ErrWeeks
	}
	p.situation = strings.TrimSpace(situation)
	p.expectedHelp = strings.TrimSpace(expectedHelp)
	p.source = strings.TrimSpace(source)

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
func (p Profile) UpdatedAt() time.Time     { return p.updatedAt }

func (p *Profile) SetUpdatedAt(t time.Time) {
	p.updatedAt = t
}

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
