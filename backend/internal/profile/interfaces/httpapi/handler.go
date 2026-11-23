package httpapi

import (
	"errors"
	"net/http"

	"backend/internal/auth"
	"backend/internal/profile/application/command"
	"backend/internal/profile/application/query"
	"backend/internal/profile/domain"
	"backend/internal/shared/httpx"
)

type Handler struct {
	saveHandler      command.SaveHandler
	setStatusHandler command.SetStatusHandler
	reads            query.ReadModel
	mapper           *httpx.ErrorMapper
}

func NewHandler(
	save command.SaveHandler,
	setStatus command.SetStatusHandler,
	reads query.ReadModel,
) *Handler {
	return &Handler{
		saveHandler:      save,
		setStatusHandler: setStatus,
		reads:            reads,
		mapper: httpx.NewErrorMapper(
			httpx.WithValidation(
				domain.ErrMother, domain.ErrNoMother,
				domain.ErrBirth, domain.ErrIncome,
				domain.ErrWeeks, domain.ErrStatus,
				domain.ErrTransition, domain.ErrNotFilled,
				domain.ErrNotComplete,
			),
		),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mother/{id}/profile", h.get)
	mux.HandleFunc("PUT /api/mother/{id}/profile", h.save)
	mux.HandleFunc("PUT /api/mother/{id}/profile/status", h.setStatus)
	mux.HandleFunc("GET /api/profiles", h.list)
}

type saveRequest struct {
	BirthDate      string `json:"birth_date"`
	BirthPlace     string `json:"birth_place"`
	Citizenship    string `json:"citizenship"`
	Registration   string `json:"registration"`
	Living         string `json:"living"`
	Email          string `json:"email"`
	Religion       string `json:"religion"`
	Education      string `json:"education"`
	Profession     string `json:"profession"`
	IncomeSource   string `json:"income_source"`
	IncomeAmount   int    `json:"income_amount"`
	MaritalStatus  string `json:"marital_status"`
	FatherInfo     string `json:"father_info"`
	PregnancyWeeks int    `json:"pregnancy_weeks"`
	Situation      string `json:"situation"`
	ExpectedHelp   string `json:"expected_help"`
	Source         string `json:"source"`
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	v, err := h.reads.GetByMother(r.Context(), motherID)
	if errors.Is(err, domain.ErrNotFilled) {
		v.MotherID = motherID
		v.Status = domain.StatusDraft
	} else if err != nil {
		h.mapper.Handle(w, "get profile", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	req, err := httpx.ReadJSON[saveRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	cmd := command.Save{
		MotherID:       motherID,
		BirthDate:      req.BirthDate,
		BirthPlace:     req.BirthPlace,
		Citizenship:    req.Citizenship,
		Registration:   req.Registration,
		Living:         req.Living,
		Email:          req.Email,
		Religion:       req.Religion,
		Education:      req.Education,
		Profession:     req.Profession,
		IncomeSource:   req.IncomeSource,
		IncomeAmount:   req.IncomeAmount,
		MaritalStatus:  req.MaritalStatus,
		FatherInfo:     req.FatherInfo,
		PregnancyWeeks: req.PregnancyWeeks,
		Situation:      req.Situation,
		ExpectedHelp:   req.ExpectedHelp,
		Source:         req.Source,
	}

	if err := h.saveHandler.Handle(r.Context(), cmd); err != nil {
		h.mapper.Handle(w, "save profile", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) setStatus(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	req, err := httpx.ReadJSON[struct {
		Status string `json:"status"`
	}](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	if err := h.setStatusHandler.Handle(r.Context(), command.SetStatus{
		MotherID: motherID,
		Status:   req.Status,
	}); err != nil {
		h.mapper.Handle(w, "set profile status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	items, err := h.reads.ListByStatus(r.Context(), status)
	if err != nil {
		h.mapper.Handle(w, "list profiles", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}
