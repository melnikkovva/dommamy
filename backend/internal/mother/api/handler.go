package api

import (
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/mother"
)

type MotherHandler struct {
	s *mother.Service
}

func NewMotherHandler(s *mother.Service) *MotherHandler{
	return &MotherHandler{s: s}
}

func (h *MotherHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/mother", h.create)
	mux.HandleFunc("GET /api/mother", h.search)
}

type motherRequest struct {
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Phone string `json:"phone"`
}

type motherResponse struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Phone string `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

func handleToMotherResponse(m *mother.Mother) motherResponse {
	return motherResponse{
		Id: m.Id, FirstName: m.FirstName, LastName: m.LastName,
		Phone: m.Phone, CreatedAt: m.CreatedAt,
	}
}

func (h *MotherHandler) create(w http.ResponseWriter, r *http.Request) {
	var req motherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}

	m, err := h.s.Create(req.FirstName, req.LastName, req.Phone)
	if err != nil {
		sendJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	sendJSON(w, http.StatusCreated, handleToMotherResponse(m))
}

func (h *MotherHandler) search(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.Search(r.URL.Query().Get("q"))
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "server error"})
		return
	}
	result := []motherResponse{}
	for _, m := range items {
		result = append(result, handleToMotherResponse(m))
	}
	sendJSON(w, http.StatusOK, result)
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

