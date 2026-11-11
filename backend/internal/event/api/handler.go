package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/internal/event"
)

type Handler struct {
	s *event.Service
}

func New(s *event.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/event-types", h.types)
	mux.HandleFunc("POST /api/events", h.create)
	mux.HandleFunc("GET /api/events", h.list)
	mux.HandleFunc("POST /api/events/{id}/participants", h.register)
	mux.HandleFunc("GET /api/events/{id}/participants", h.participants)
	mux.HandleFunc("PUT /api/events/{id}/participants/{motherId}", h.mark)
	mux.HandleFunc("GET /api/mother/{id}/events", h.byMother)
}

type eventRequest struct {
	TypeId int `json:"type_id"`
	Title string `json:"title"`
	StartsAt string `json:"starts_at"`
	Place string `json:"place"`
	MaxPeople int `json:"max_people"`
}

type eventResponse struct {
	Id int `json:"id"`
	TypeId int `json:"type_id"`
	Title string `json:"title"`
	StartsAt string `json:"starts_at"`
	Place string `json:"place"`
	MaxPeople int `json:"max_people"`
}

func toResponse(e *event.Event) eventResponse {
	return eventResponse{
		Id: e.Id, TypeId: e.TypeId, Title: e.Title,
		StartsAt: e.StartsAt.Format("2006-01-02T15:04"),
		Place: e.Place, MaxPeople: e.MaxPeople,
	}
}

func (h *Handler) types(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.Types()
	if err != nil {
		return
	}
	result := []map[string]any{}
	for _, t := range items {
		result = append(result, map[string]any{"id": t.Id, "name": t.Name})
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	startsAt, err := time.Parse("2006-01-02T15:04", req.StartsAt)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "specify the date and time"})
		return
	}
	e, err := h.s.Create(req.TypeId, req.Title, startsAt, req.Place, req.MaxPeople)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(e))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List()
	if err != nil {
		return
	}
	writeJSON(w, http.StatusOK, toResponses(items))
}

func (h *Handler) byMother(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	items, err := h.s.ByMother(motherId)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusOK, toResponses(items))
}

func toResponses(items []*event.Event) []eventResponse {
	result := []eventResponse{}
	for _, e := range items {
		result = append(result, toResponse(e))
	}
	return result
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		MotherId int `json:"mother_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	if err := h.s.Register(eventId, req.MotherId); err != nil {
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *Handler) participants(w http.ResponseWriter, r *http.Request) {
	eventId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	items, err := h.s.Participants(eventId)
	if err != nil {
		return
	}
	result := []map[string]any{}
	for _, p := range items {
		result = append(result, map[string]any{"mother_id": p.MotherId, "status": p.Status})
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) mark(w http.ResponseWriter, r *http.Request) {
	eventId, err1 := strconv.Atoi(r.PathValue("id"))
	motherId, err2 := strconv.Atoi(r.PathValue("motherId"))
	if err1 != nil || err2 != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	if err := h.s.Mark(eventId, motherId, req.Status); err != nil {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}