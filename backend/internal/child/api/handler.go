package api

import (
	"errors"
	"net/http"
	"time"

	"backend/internal/auth"
	"backend/internal/child"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s *child.Service
}

func New(s *child.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother/{id}/children", h.add)
	mux.HandleFunc("GET /api/mother/{id}/children", h.list)
}

type childRequest struct {
	FirstName string `json:"first_name"`
	BirthDate string `json:"birth_date"`
}

type childResponse struct {
	Id int `json:"id"`
	MotherId int `json:"mother_id"`
	FirstName string `json:"first_name"`
	BirthDate string `json:"birth_date"`
}

func handleToResponse(c *child.Child) childResponse {
	return childResponse{
		Id: c.Id, MotherId: c.MotherId, FirstName: c.FirstName,
		BirthDate: c.BirthDate.Format(httpx.APITime),
	}
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	motherId, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherId) {
		return
	}

	req, err := httpx.ReadJSON[childRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	birth, err := time.Parse(httpx.LocalDate, req.BirthDate)
	if err != nil {
		birth, err = time.Parse(httpx.APITime, req.BirthDate)
	}
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "specify the date of birth")
		return
	}

	c, err := h.s.Add(r.Context(), motherId, req.FirstName, birth)
	if err != nil {
		switch {
		case errors.Is(err, child.ErrorName), errors.Is(err, child.ErrorBirth),
			errors.Is(err, child.ErrorMother), errors.Is(err, child.ErrorNoMother):
			httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, handleToResponse(c))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	motherId, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherId) {
		return
	}
	items, err := h.s.SearchByMother(motherId)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	result := []childResponse{}
	for _, c := range items {
		result = append(result, handleToResponse(c))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}