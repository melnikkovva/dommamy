package command

import (
	"context"
	"errors"
	"time"

	"backend/internal/profile/domain"
)

type Save struct {
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
}

type SaveHandler interface {
	Handle(ctx context.Context, cmd Save) error
}

func NewSaveHandler(repo domain.Repository, mothers MotherDirectory, now func() time.Time) SaveHandler {
	return &saveHandler{repo: repo, mothers: mothers, now: now}
}

type saveHandler struct {
	repo    domain.Repository
	mothers MotherDirectory
	now     func() time.Time
}

func (h *saveHandler) Handle(ctx context.Context, cmd Save) error {
	ok, err := h.mothers.Exists(ctx, cmd.MotherID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNoMother
	}

	p, err := h.repo.ByMother(ctx, cmd.MotherID)
	if errors.Is(err, domain.ErrNotFilled) {
		p, err = domain.NewProfile(cmd.MotherID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := p.Fill(domain.ProfileInput{
		BirthDate:      cmd.BirthDate,
		BirthPlace:     cmd.BirthPlace,
		Citizenship:    cmd.Citizenship,
		Registration:   cmd.Registration,
		Living:         cmd.Living,
		Email:          cmd.Email,
		Religion:       cmd.Religion,
		Education:      cmd.Education,
		Profession:     cmd.Profession,
		IncomeSource:   cmd.IncomeSource,
		IncomeAmount:   cmd.IncomeAmount,
		MaritalStatus:  cmd.MaritalStatus,
		FatherInfo:     cmd.FatherInfo,
		PregnancyWeeks: cmd.PregnancyWeeks,
		Situation:      cmd.Situation,
		ExpectedHelp:   cmd.ExpectedHelp,
		Source:         cmd.Source,
	}, h.now()); err != nil {
		return err
	}

	return h.repo.Save(ctx, p)
}
