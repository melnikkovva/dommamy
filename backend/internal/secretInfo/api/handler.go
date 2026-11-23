package api

import (
	"errors"
	"net/http"

	"backend/internal/auth"
	"backend/internal/secretInfo"
	"backend/internal/shared/httpx"
)

type Handler struct {
	s *secretInfo.Service
}

func New(s *secretInfo.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mother/{id}/sensitive", h.reveal)
	mux.HandleFunc("PUT /api/mother/{id}/sensitive", h.save)
}

type dataBody struct {
	CriminalRecord string `json:"criminal_record"`
	NdPnd string `json:"nd_pnd"`
	Drugs string `json:"drugs"`
	ChronicDiseases string `json:"chronic_diseases"`
}

func (h *Handler) reveal(w http.ResponseWriter, r *http.Request) {
	motherId, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherId) {
		return
	}

	d, err := h.s.Reveal(motherId)
	if errors.Is(err, secretInfo.ErrorNotFound) {
		httpx.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dataBody{
		CriminalRecord: d.CriminalRecord,
		NdPnd: d.NdPnd,
		Drugs: d.Drugs,
		ChronicDiseases: d.ChronicDiseases,
	})
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	motherId, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherId) {
		return
	}
	req, err := httpx.ReadJSON[dataBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	err = h.s.Save(r.Context(), motherId, req.CriminalRecord, req.NdPnd, req.Drugs, req.ChronicDiseases)
	if errors.Is(err, secretInfo.ErrorNoMother) || errors.Is(err, secretInfo.ErrorMother) {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}