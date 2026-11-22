package httpapi

import (
	"errors"
	"log"
	"net/http"

	lifecyclecommand "backend/internal/lifecycle/application/command"
	lifecyclequery "backend/internal/lifecycle/application/query"
	"backend/internal/lifecycle/domain"
	"backend/internal/shared/httpx"
)

type Handler struct {
	changeHandler lifecyclecommand.ChangeHandler
	reads         lifecyclequery.ReadModel
}

func NewHandler(
	change lifecyclecommand.ChangeHandler,
	reads lifecyclequery.ReadModel,
) *Handler {
	return &Handler{
		changeHandler: change,
		reads:         reads,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/mother/{id}/status", h.get)
	mux.HandleFunc("PUT /api/mother/{id}/status", h.change)
	mux.HandleFunc("GET /api/mother-statuses", h.all)
}

type changeRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	view, err := h.reads.GetStatus(r.Context(), motherID)
	if err != nil {
		fail(w, "get status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, view)
}

func (h *Handler) change(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}

	req, err := httpx.ReadJSON[changeRequest](r)
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}

	if err := h.changeHandler.Handle(r.Context(), lifecyclecommand.Change{
		MotherID: motherID,
		Status:   req.Status,
		Reason:   req.Reason,
	}); err != nil {
		fail(w, "change status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) all(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.Currents(r.Context())
	if err != nil {
		fail(w, "currents", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, domain.ErrMother), errors.Is(err, domain.ErrNoMother),
		errors.Is(err, domain.ErrStatus), errors.Is(err, domain.ErrTransition),
		errors.Is(err, domain.ErrReason):
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Println("error", where+":", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
	}
}
