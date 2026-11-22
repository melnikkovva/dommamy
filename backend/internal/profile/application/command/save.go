package command

import (
	"context"
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
	ok, err := h.mothers.Exists(cmd.MotherID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNoMother
	}

	p, err := h.repo.ByMother(ctx, cmd.MotherID)
	if err != nil {
		return err
	}

	if p.MotherID() == 0 {
		p, err = domain.NewProfile(cmd.MotherID)
		if err != nil {
			return err
		}
	}

	if err := p.Fill(
		cmd.BirthDate, cmd.BirthPlace, cmd.Citizenship, cmd.Registration,
		cmd.Living, cmd.Email, cmd.Religion, cmd.Education, cmd.Profession,
		cmd.IncomeSource, cmd.MaritalStatus, cmd.FatherInfo, cmd.Source,
		cmd.IncomeAmount, cmd.PregnancyWeeks, cmd.Situation, cmd.ExpectedHelp,
	); err != nil {
		return err
	}

	p.SetUpdatedAt(h.now())
	return h.repo.Save(ctx, p)
}
