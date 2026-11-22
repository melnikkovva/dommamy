package query

import (
	"context"
	"time"

	"backend/internal/profile/domain"
)

type ProfileView struct {
	MotherID       int       `json:"mother_id"`
	BirthDate      string    `json:"birth_date"`
	BirthPlace     string    `json:"birth_place"`
	Citizenship    string    `json:"citizenship"`
	Registration   string    `json:"registration"`
	Living         string    `json:"living"`
	Email          string    `json:"email"`
	Religion       string    `json:"religion"`
	Education      string    `json:"education"`
	Profession     string    `json:"profession"`
	IncomeSource   string    `json:"income_source"`
	IncomeAmount   int       `json:"income_amount"`
	MaritalStatus  string    `json:"marital_status"`
	FatherInfo     string    `json:"father_info"`
	PregnancyWeeks int       `json:"pregnancy_weeks"`
	Situation      string    `json:"situation"`
	ExpectedHelp   string    `json:"expected_help"`
	Source         string    `json:"source"`
	Status         domain.Status `json:"status"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProfileListItem struct {
	MotherID  int             `json:"mother_id"`
	Status    domain.Status   `json:"status"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ReadModel interface {
	GetByMother(ctx context.Context, motherID int) (ProfileView, error)
	ListByStatus(ctx context.Context, status string) ([]ProfileListItem, error)
}
