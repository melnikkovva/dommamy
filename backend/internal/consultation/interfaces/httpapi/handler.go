package httpapi

import (
	"errors"
	"net/http"
	"time"

	"backend/internal/consultation/domain"
	consultationmysql "backend/internal/consultation/infrastructure/mysql"
	"backend/internal/shared/httpx"
)

type Handler struct {
	repo *consultationmysql.Repo
}

func NewHandler(repo *consultationmysql.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/consultations", h.create)
	mux.HandleFunc("GET /api/consultations", h.list)
	mux.HandleFunc("GET /api/consultations/{id}", h.get)
	mux.HandleFunc("PUT /api/consultations/{id}", h.update)
	mux.HandleFunc("DELETE /api/consultations/{id}", h.delete)
	mux.HandleFunc("GET /api/consultations/mother/{motherId}", h.byMother)
}

type createRequest struct {
	MotherId         int    `json:"mother_id"`
	SpecialistType   string `json:"specialist_type"`
	ConsultationDate string `json:"consultation_date"`
	Notes            string `json:"notes"`
	Status           string `json:"status"`
}

type consultationResponse struct {
	Id               int    `json:"id"`
	MotherId         int    `json:"mother_id"`
	SpecialistType   string `json:"specialist_type"`
	ConsultationDate string `json:"consultation_date"`
	Notes            string `json:"notes"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
}

type rowResponse struct {
	Id               int    `json:"id"`
	MotherId         int    `json:"mother_id"`
	MotherName       string `json:"mother_name"`
	SpecialistType   string `json:"specialist_type"`
	ConsultationDate string `json:"consultation_date"`
	Notes            string `json:"notes"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
}

func toResponse(c *domain.Consultation) consultationResponse {
	return consultationResponse{
		Id:               c.Id,
		MotherId:         c.MotherId,
		SpecialistType:   string(c.SpecialistType),
		ConsultationDate: c.ConsultationDate.Format("2006-01-02"),
		Notes:            c.Notes,
		Status:           string(c.Status),
		CreatedAt:        c.CreatedAt.Format(httpx.APITime),
	}
}

func rowToResponse(r domain.Row) rowResponse {
	return rowResponse{
		Id:               r.Id,
		MotherId:         r.MotherId,
		MotherName:       r.MotherName,
		SpecialistType:   string(r.SpecialistType),
		ConsultationDate: r.ConsultationDate.Format("2006-01-02"),
		Notes:            r.Notes,
		Status:           string(r.Status),
		CreatedAt:        r.CreatedAt.Format(httpx.APITime),
	}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	req, err := httpx.ReadJSON[createRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	consultationDate, err := time.Parse("2006-01-02", req.ConsultationDate)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "неверный формат даты")
		return
	}

	status := domain.Status(req.Status)
	if status == "" {
		status = domain.Completed
	}

	c, err := domain.NewConsultation(domain.ConsultationInput{
		MotherId:         req.MotherId,
		SpecialistType:   domain.SpecialistType(req.SpecialistType),
		ConsultationDate: consultationDate,
		Notes:            req.Notes,
		Status:           status,
	})
	if err != nil {
		h.writeDomainError(w, err)
		return
	}

	if err := h.repo.Save(r.Context(), c); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка сохранения")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toResponse(c))
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	rows, err := h.repo.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка загрузки")
		return
	}
	result := make([]rowResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToResponse(row))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "неверный id")
		return
	}
	rows, err := h.repo.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка загрузки")
		return
	}
	for _, row := range rows {
		if int(row.Id) == id {
			httpx.WriteJSON(w, http.StatusOK, rowToResponse(row))
			return
		}
	}
	httpx.WriteError(w, http.StatusNotFound, "консультация не найдена")
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "неверный id")
		return
	}

	req, err := httpx.ReadJSON[struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	if req.Status != "" {
		if err := h.repo.UpdateStatus(r.Context(), id, domain.Status(req.Status)); err != nil {
			httpx.WriteError(w, http.StatusInternalServerError, "ошибка обновления")
			return
		}
	}

	rows, err := h.repo.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка загрузки")
		return
	}
	for _, row := range rows {
		if int(row.Id) == id {
			httpx.WriteJSON(w, http.StatusOK, rowToResponse(row))
			return
		}
	}
	httpx.WriteError(w, http.StatusNotFound, "консультация не найдена")
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "неверный id")
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка удаления")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) byMother(w http.ResponseWriter, r *http.Request) {
	motherId, err := httpx.PathID(r, "motherId")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "неверный id мамы")
		return
	}
	rows, err := h.repo.ByMother(r.Context(), motherId)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "ошибка загрузки")
		return
	}
	result := make([]rowResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, rowToResponse(row))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) writeDomainError(w http.ResponseWriter, err error) {
	if errors.Is(err, domain.ErrorDate) || errors.Is(err, domain.ErrorType) ||
		errors.Is(err, domain.ErrorUnknown) || errors.Is(err, domain.ErrorMum) {
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	httpx.WriteError(w, http.StatusInternalServerError, "что-то пошло не так")
}
