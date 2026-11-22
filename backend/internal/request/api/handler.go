package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/request"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s *request.Service
}

func New(s *request.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/requests", h.create)
	mux.HandleFunc("GET /api/requests", h.list)
	mux.HandleFunc("PUT /api/requests/{id}/status", h.setStatus)
	mux.HandleFunc("POST /api/requests/{id}/issue", h.issue)
	mux.HandleFunc("GET /api/mother/{id}/requests", h.byMother)
}

type response struct {
	Id int `json:"id"`
	MotherId int `json:"mother_id"`
	Kind string `json:"kind"`
	Description string `json:"description"`
	Status string `json:"status"`
	Comment string `json:"comment"`
	Source string `json:"source"`
	HelpId int `json:"help_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toResponse(r *request.Request) response {
	return response{
		Id: r.Id, MotherId: r.MotherId, Kind: r.Kind, Description: r.Description,
		Status: r.Status, Comment: r.Comment, Source: r.Source, HelpId: r.HelpId,
		CreatedAt: r.CreatedAt.Format(time.RFC3339), UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	}
}

func toResponses(items []*request.Request) []response {
	result := []response{}
	for _, r := range items {
		result = append(result, toResponse(r))
	}
	return result
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MotherId int `json:"mother_id"`
		Kind string `json:"kind"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	created, err := h.s.Create(req.MotherId, req.Kind, req.Description, request.SourceStaff)
	if err != nil {
		fail(w, "create request", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toResponse(created))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List(r.URL.Query().Get("status"))
	if err != nil {
		fail(w, "list requests", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *Handler) byMother(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	items, err := h.s.ByMother(motherId)
	if err != nil {
		fail(w, "requests by mother", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *Handler) setStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Status string `json:"status"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	updated, err := h.s.SetStatus(id, req.Status, req.Comment)
	if err != nil {
		fail(w, "set request status", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(updated))
}

func (h *Handler) issue(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		HelpId int `json:"help_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	updated, err := h.s.Fulfill(id, req.HelpId)
	if err != nil {
		fail(w, "fulfill request", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponse(updated))
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, request.ErrorNotFound):
		httpx.WriteJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, request.ErrorMother), errors.Is(err, request.ErrorNoMother),
		errors.Is(err, request.ErrorKind), errors.Is(err, request.ErrorSource),
		errors.Is(err, request.ErrorStatus), errors.Is(err, request.ErrorTransition),
		errors.Is(err, request.ErrorComment), errors.Is(err, request.ErrorUseIssue),
		errors.Is(err, request.ErrorNoHelp), errors.Is(err, request.ErrorHelpMismatch):
		httpx.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		log.Println("error", where+":", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
	}
}