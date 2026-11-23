package httpapi

import (
	"net/http"

	"backend/internal/auth"
	lifecyclecommand "backend/internal/lifecycle/application/command"
	lifecyclequery "backend/internal/lifecycle/application/query"
	"backend/internal/lifecycle/domain"
	"backend/internal/shared/httpx"
)

type Handler struct {
	changeHandler lifecyclecommand.ChangeHandler
	reads         lifecyclequery.ReadModel
	mapper        *httpx.ErrorMapper
}

func NewHandler(
	change lifecyclecommand.ChangeHandler,
	reads lifecyclequery.ReadModel,
) *Handler {
	return &Handler{
		changeHandler: change,
		reads:         reads,
		mapper: httpx.NewErrorMapper(
			httpx.WithValidation(
				domain.ErrMother, domain.ErrNoMother,
				domain.ErrStatus, domain.ErrTransition,
				domain.ErrReason,
			),
		),
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
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	view, err := h.reads.GetStatus(r.Context(), motherID)
	if err != nil {
		h.mapper.Handle(w, "get status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, view)
}

func (h *Handler) change(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}

	req, err := httpx.ReadJSON[changeRequest](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	if err := h.changeHandler.Handle(r.Context(), lifecyclecommand.Change{
		MotherID: motherID,
		Status:   req.Status,
		Reason:   req.Reason,
	}); err != nil {
		h.mapper.Handle(w, "change status", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) all(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.Currents(r.Context())
	if err != nil {
		h.mapper.Handle(w, "currents", err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, items)
}
