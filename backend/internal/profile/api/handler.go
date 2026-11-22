package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/profile"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s *profile.Service
}

func New(s *profile.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mother/{id}/profile", h.get)
	mux.HandleFunc("PUT /api/mother/{id}/profile", h.save)
	mux.HandleFunc("PUT /api/mother/{id}/profile/status", h.setStatus)
	mux.HandleFunc("GET /api/profiles", h.list)
}

type profileBody struct {
	BirthDate string `json:"birth_date"`
	BirthPlace string `json:"birth_place"`
	Citizenship string `json:"citizenship"`
	Registration string `json:"registration"`
	Living string `json:"living"`
	Email string `json:"email"`
	Religion string `json:"religion"`
	Education string `json:"education"`
	Profession string `json:"profession"`
	IncomeSource string `json:"income_source"`
	IncomeAmount int `json:"income_amount"`
	MaritalStatus string `json:"marital_status"`
	FatherInfo string `json:"father_info"`
	PregnancyWeeks int `json:"pregnancy_weeks"`
	Situation string `json:"situation"`
	ExpectedHelp string `json:"expected_help"`
	Source string `json:"source"`
}

type profileResponse struct {
	profileBody
	MotherId int `json:"mother_id"`
	Status string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type profileListItem struct {
	MotherId int `json:"mother_id"`
	Status string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

func toResponse(p *profile.Profile) profileResponse {
	updated := ""
	if !p.UpdatedAt.IsZero() {
		updated = p.UpdatedAt.Format(time.RFC3339)
	}
	return profileResponse{
		profileBody: profileBody{
			BirthDate: p.BirthDate, BirthPlace: p.BirthPlace, Citizenship: p.Citizenship,
			Registration: p.Registration, Living: p.Living, Email: p.Email,
			Religion: p.Religion, Education: p.Education, Profession: p.Profession,
			IncomeSource: p.IncomeSource, IncomeAmount: p.IncomeAmount,
			MaritalStatus: p.MaritalStatus, FatherInfo: p.FatherInfo,
			PregnancyWeeks: p.PregnancyWeeks, Situation: p.Situation,
			ExpectedHelp: p.ExpectedHelp, Source: p.Source,
		},
		MotherId: p.MotherId, Status: p.Status, UpdatedAt: updated,
	}
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	p, err := h.s.Get(motherId)
	if err != nil {
		fail(w, "get profile", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(p))
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req profileBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	p := &profile.Profile{
		MotherId: motherId, BirthDate: req.BirthDate, BirthPlace: req.BirthPlace,
		Citizenship: req.Citizenship, Registration: req.Registration, Living: req.Living,
		Email: req.Email, Religion: req.Religion, Education: req.Education,
		Profession: req.Profession, IncomeSource: req.IncomeSource, IncomeAmount: req.IncomeAmount,
		MaritalStatus: req.MaritalStatus, FatherInfo: req.FatherInfo,
		PregnancyWeeks: req.PregnancyWeeks, Situation: req.Situation,
		ExpectedHelp: req.ExpectedHelp, Source: req.Source,
	}
	if err := h.s.Save(p); err != nil {
		fail(w, "save profile", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) setStatus(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	if err := h.s.SetStatus(motherId, req.Status); err != nil {
		fail(w, "set profile status", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List(r.URL.Query().Get("status"))
	if err != nil {
		fail(w, "list profiles", err)
		return
	}
	result := []profileListItem{}
	for _, p := range items {
		result = append(result, profileListItem{
			MotherId: p.MotherId,
			Status: p.Status,
			UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		})
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, profile.ErrorMother), errors.Is(err, profile.ErrorNoMother),
		errors.Is(err, profile.ErrorBirth), errors.Is(err, profile.ErrorIncome),
		errors.Is(err, profile.ErrorWeeks), errors.Is(err, profile.ErrorStatus),
		errors.Is(err, profile.ErrorTransition), errors.Is(err, profile.ErrorNotFilled),
		errors.Is(err, profile.ErrorNotComplete):
		httpx.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		log.Println("error", where+":", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
	}
}