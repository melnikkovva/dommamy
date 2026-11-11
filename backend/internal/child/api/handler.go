package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/child"
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

func handleToResponse (c *child.Child) childResponse {
	return childResponse{
		Id: c.Id, MotherId: c.MotherId, FirstName: c.FirstName,
		BirthDate: c.BirthDate.Format("2006-01-02"),
	}
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id mum"})
		return
	}

	var req childRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	birth, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "specify the date of birth"})
		return
	}

	c, err := h.s.Add(motherId, req.FirstName, birth)
	if err != nil {
		if errors.Is(err, child.ErrorName) || errors.Is(err, child.ErrorBirth) ||
			errors.Is(err, child.ErrorMother) || errors.Is(err, child.ErrorNoMother) {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
			return
		}
		log.Println("error of adding child:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	writeJSON(w, http.StatusCreated, handleToResponse(c))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect mum id"})
		return
	}
	items, err := h.s.SearchByMother(motherId)
	if err != nil {
		log.Println("error list of child:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	result := []childResponse{}
	for _, c := range items {
		result = append(result, handleToResponse(c))
	}
	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}