package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"fmt"
	"backend/internal/help"
)

type Handler struct {
	s *help.Service
}

func New(s *help.Service) *Handler {
	return &Handler{ s : s }
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/help", h.item)
	mux.HandleFunc("GET /api/help", h.list)
}

type itemRequest struct {
	MotherId int `json:"mother_id"`
	Kind string `json:"kind"`
	Sum int `json:"sum"`
	CheckUrl string `json:"check_url"`
}

type itemResponse struct {
	Id int `json:"id"`
	MotherId int `json:"mother_id"`
	Kind string `json:"kind"`
	Sum int `json:"sum"`
	CheckUrl string `json:"check_url"`
	CreatedAt time.Time `json:"created_at"`
}

func handleToResponse(i *help.Item) itemResponse {
	return itemResponse{
		Id: i.Id, MotherId: i.MotherId, Kind: string(i.Kind),
		Sum: i.Sum, CheckUrl: i.CheckUrl, CreatedAt: i.CreatedAt,
	}
}

func (h *Handler) item(w http.ResponseWriter, r *http.Request) {
	var req itemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	item, err := h.s.Item(r.Context(), req.MotherId, help.Kind(req.Kind), req.Sum, req.CheckUrl)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, handleToResponse(item))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "internal database error")
		return
	}

	result := make([]itemResponse, 0, len(items))
	for _, i := range items {
		result = append(result, handleToResponse(i))
	}
	h.writeJSON(w, http.StatusOK, result)
}

func (h *Handler) writeDomainError(w http.ResponseWriter, err error) {
	if errors.Is(err, help.ErrorCheck) || errors.Is(err, help.ErrorSum) ||
		errors.Is(err, help.ErrorMum) || errors.Is(err, help.ErrorKind) {

		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	h.writeError(w, http.StatusInternalServerError, "something went wrong")
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": "%s"}`, msg)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


