package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/donor"
)

type Handler struct {
	s *donor.Service
}

func New(s *donor.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/donors", h.create)
	mux.HandleFunc("GET /api/donors", h.list)
	mux.HandleFunc("GET /api/donors/{id}", h.get)
	mux.HandleFunc("PUT /api/donors/{id}", h.update)
	mux.HandleFunc("POST /api/donors/{id}/donations", h.addDonation)
}

type donorBody struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Comment string `json:"comment"`
}

type donorResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Comment string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type summaryResponse struct {
	donorResponse
	Total int `json:"total"`
	Count int `json:"count"`
}

type donationResponse struct {
	Id int `json:"id"`
	DonorId int `json:"donor_id"`
	Amount int `json:"amount"`
	Date string `json:"date"`
	Comment string `json:"comment"`
}

func toDonor(d *donor.Donor) donorResponse {
	return donorResponse{
		Id: d.Id, Name: d.Name, Kind: d.Kind, Phone: d.Phone, Email: d.Email,
		Comment: d.Comment, CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}
}

func toDonation(n *donor.Donation) donationResponse {
	return donationResponse{Id: n.Id, DonorId: n.DonorId, Amount: n.Amount, Date: n.Date, Comment: n.Comment}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req donorBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	d, err := h.s.Create(req.Name, req.Kind, req.Phone, req.Email, req.Comment)
	if err != nil {
		fail(w, "create donor", err)
		return
	}
	writeJSON(w, http.StatusCreated, toDonor(d))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.List(r.URL.Query().Get("q"))
	if err != nil {
		fail(w, "list donors", err)
		return
	}
	result := []summaryResponse{}
	for _, s := range items {
		result = append(result, summaryResponse{donorResponse: toDonor(&s.Donor), Total: s.Total, Count: s.Count})
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	d, donations, err := h.s.Get(id)
	if err != nil {
		fail(w, "get donor", err)
		return
	}
	list := []donationResponse{}
	for _, n := range donations {
		list = append(list, toDonation(n))
	}
	writeJSON(w, http.StatusOK, map[string]any{"donor": toDonor(d), "donations": list})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req donorBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	d, err := h.s.Update(id, req.Name, req.Kind, req.Phone, req.Email, req.Comment)
	if err != nil {
		fail(w, "update donor", err)
		return
	}
	writeJSON(w, http.StatusOK, toDonor(d))
}

func (h *Handler) addDonation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	var req struct {
		Amount int `json:"amount"`
		Date string `json:"date"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	n, err := h.s.AddDonation(id, req.Amount, req.Date, req.Comment)
	if err != nil {
		fail(w, "add donation", err)
		return
	}
	writeJSON(w, http.StatusCreated, toDonation(n))
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, donor.ErrorNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, donor.ErrorName), errors.Is(err, donor.ErrorKind),
		errors.Is(err, donor.ErrorEmail), errors.Is(err, donor.ErrorAmount),
		errors.Is(err, donor.ErrorDate):
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