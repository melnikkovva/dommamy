package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/lifecycle"
)

type Handler struct {
	s *lifecycle.Service
}

func New(s *lifecycle.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mother/{id}/status", h.get)
	mux.HandleFunc("PUT /api/mother/{id}/status", h.change)
	mux.HandleFunc("GET /api/mother-statuses", h.all)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	current, err := h.s.Status(motherId)
	if err != nil {
		fail(w, "status", err)
		return
	}
	items, err := h.s.History(motherId)
	if err != nil {
		fail(w, "history", err)
		return
	}
	history := []map[string]string{}
	for _, c := range items {
		history = append(history, map[string]string{
			"status": c.Status,
			"reason": c.Reason,
			"created_at": c.CreatedAt.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": current, "history": history})
}

func (h *Handler) change(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	if err := h.s.Change(motherId, req.Status, req.Reason); err != nil {
		fail(w, "change status", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) all(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.Currents()
	if err != nil {
		fail(w, "currents", err)
		return
	}
	result := []map[string]any{}
	for _, c := range items {
		result = append(result, map[string]any{"mother_id": c.MotherId, "status": c.Status})
	}
	writeJSON(w, http.StatusOK, result)
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, lifecycle.ErrorMother), errors.Is(err, lifecycle.ErrorNoMother),
		errors.Is(err, lifecycle.ErrorStatus), errors.Is(err, lifecycle.ErrorTransition),
		errors.Is(err, lifecycle.ErrorReason):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		log.Println("error", where+":", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}