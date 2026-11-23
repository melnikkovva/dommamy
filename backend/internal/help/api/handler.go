package api

import (
	"errors"
	"net/http"

	"backend/internal/help"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s *help.Service
}

func New(s *help.Service) *Handler {
	return &Handler{s: s}
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
	Id       int    `json:"id"`
	MotherId int    `json:"mother_id"`
	Kind     string `json:"kind"`
	Sum      int    `json:"sum"`
	CheckUrl string `json:"check_url"`
	CreatedAt string `json:"created_at"`
}

func handleToResponse(i *help.Item) itemResponse {
	return itemResponse{
		Id: i.Id, MotherId: i.MotherId, Kind: string(i.Kind),
		Sum: i.Sum, CheckUrl: i.CheckUrl, CreatedAt: i.CreatedAt.Format(httpx.APITime),
	}
}

func (h *Handler) item(w http.ResponseWriter, r *http.Request) {
	req, err := httpx.ReadJSON[itemRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	item, err := h.s.Item(r.Context(), req.MotherId, help.Kind(req.Kind), req.Sum, req.CheckUrl)
	if err != nil {
		h.writeDomainError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, handleToResponse(item))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal database error")
		return
	}

	result := make([]itemResponse, 0, len(items))
	for _, i := range items {
		result = append(result, handleToResponse(i))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) writeDomainError(w http.ResponseWriter, err error) {
	if errors.Is(err, help.ErrorCheck) || errors.Is(err, help.ErrorSum) ||
		errors.Is(err, help.ErrorMum) || errors.Is(err, help.ErrorKind) {

		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
}


