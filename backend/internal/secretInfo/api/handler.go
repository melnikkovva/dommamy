package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"backend/internal/secretInfo"
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
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id mum"})
		return
	}

	d, err := h.s.Reveal(motherId)
	if errors.Is(err, secretInfo.ErrorNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		log.Println("error reveal secret info:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	writeJSON(w, http.StatusOK, dataBody{
		CriminalRecord: d.CriminalRecord,
		NdPnd: d.NdPnd,
		Drugs: d.Drugs,
		ChronicDiseases: d.ChronicDiseases,
	})
}

func (h *Handler) save(w http.ResponseWriter, r *http.Request) {
	motherId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id mum"})
		return
	}
	var req dataBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}

	err = h.s.Save(motherId, req.CriminalRecord, req.NdPnd, req.Drugs, req.ChronicDiseases)
	if errors.Is(err, secretInfo.ErrorNoMother) || errors.Is(err, secretInfo.ErrorMother) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}
	if err != nil {
		log.Println("error save secret info:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}