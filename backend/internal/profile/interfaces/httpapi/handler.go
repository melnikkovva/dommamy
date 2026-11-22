package httpapi

import (
	"errors"
	"log"
	"net/http"

	"backend/internal/profile/application/command"
	"backend/internal/profile/application/query"
	"backend/internal/profile/domain"
	"backend/internal/shared/httpx"
)

type Handler struct {
	saveHandler    command.SaveHandler
	setStatusHandler command.SetStatusHandler
	reads          query.ReadModel
}

func NewHandler(
	save command.SaveHandler,
	setStatus command.SetStatusHandler,
	reads query.ReadModel,
) *Handler {
	return &Handler{
		saveHandler:    save,
		setStatusHandler: setStatus,
		reads:          reads,
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
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	v, err := h.reads.GetByMother(r.Context(), motherID)
	if err != nil {
		fail(w, "get profile", err)
		return
	}

	if v.MotherID == 0 {
		v.MotherID = motherID
		v.Status = domain.StatusDraft
	}

	httpx.WriteJSON(w, http.StatusOK, v)
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	req, err := httpx.ReadJSON[saveRequest](r)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
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
		fail(w, "save profile", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) setStatus(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	req, err := httpx.ReadJSON[struct {
		Status string `json:"status"`
	}](r)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}

	if err := h.setStatusHandler.Handle(r.Context(), command.SetStatus{
		MotherID: motherID,
		Status:   req.Status,
	}); err != nil {
		fail(w, "set profile status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	items, err := h.reads.ListByStatus(r.Context(), status)
	if err != nil {
		fail(w, "list profiles", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, domain.ErrMother), errors.Is(err, domain.ErrNoMother),
		errors.Is(err, domain.ErrBirth), errors.Is(err, domain.ErrIncome),
		errors.Is(err, domain.ErrWeeks), errors.Is(err, domain.ErrStatus),
		errors.Is(err, domain.ErrTransition), errors.Is(err, domain.ErrNotFilled),
		errors.Is(err, domain.ErrNotComplete):
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Println("error", where+":", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
	}
}
