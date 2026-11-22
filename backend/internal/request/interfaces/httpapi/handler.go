package httpapi

import (
	"errors"
	"log"
	"net/http"
	"slices"
	"time"

	"backend/internal/request/application/command"
	"backend/internal/request/application/query"
	"backend/internal/request/domain"
	"backend/internal/shared/httpx"
)

func NewHandler(
	create command.CreateRequestHandler,
	changeStatus command.ChangeStatusHandler,
	fulfill command.FulfillRequestHandler,
	reads query.ReadModel,
) httpx.Router {
	return &handler{create: create, changeStatus: changeStatus, fulfill: fulfill, reads: reads}
}

type handler struct {
	create       command.CreateRequestHandler
	changeStatus command.ChangeStatusHandler
	fulfill      command.FulfillRequestHandler
	reads        query.ReadModel
}

func (h *handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/requests", h.createRequest)
	mux.HandleFunc("GET /api/requests", h.list)
	mux.HandleFunc("PUT /api/requests/{id}/status", h.setStatus)
	mux.HandleFunc("POST /api/requests/{id}/issue", h.issue)
	mux.HandleFunc("GET /api/mother/{id}/requests", h.byMother)
}

type response struct {
	ID          int    `json:"id"`
	MotherID    int    `json:"mother_id"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
	Source      string `json:"source"`
	HelpID      int    `json:"help_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toResponse(v query.RequestView) response {
	return response{
		ID: v.ID, MotherID: v.MotherID, Kind: v.Kind, Description: v.Description,
		Status: v.Status, Comment: v.Comment, Source: v.Source, HelpID: v.HelpID,
		CreatedAt: v.CreatedAt.Format(time.RFC3339), UpdatedAt: v.UpdatedAt.Format(time.RFC3339),
	}
}

func toResponses(items []query.RequestView) []response {
	result := make([]response, 0, len(items))
	for _, v := range items {
		result = append(result, toResponse(v))
	}
	return result
}

type createBody struct {
	MotherID    int    `json:"mother_id"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
}

type statusBody struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type issueBody struct {
	HelpID int `json:"help_id"`
}

func (h *handler) createRequest(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ReadJSON[createBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	id, err := h.create.Handle(r.Context(), command.CreateRequest{
		MotherID:    body.MotherID,
		Kind:        body.Kind,
		Description: body.Description,
		Source:      string(domain.SourceStaff),
	})
	if err != nil {
		fail(w, "create request", err)
		return
	}
	h.respondWith(w, r, http.StatusCreated, id)
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		fail(w, "list requests", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *handler) byMother(w http.ResponseWriter, r *http.Request) {
	motherID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	items, err := h.reads.ByMother(r.Context(), motherID)
	if err != nil {
		fail(w, "requests by mother", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *handler) setStatus(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[statusBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	err = h.changeStatus.Handle(r.Context(), command.ChangeStatus{RequestID: id, Status: body.Status, Comment: body.Comment})
	if err != nil {
		fail(w, "set request status", err)
		return
	}
	h.respondWith(w, r, http.StatusOK, id)
}

func (h *handler) issue(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[issueBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if err := h.fulfill.Handle(r.Context(), command.FulfillRequest{RequestID: id, HelpID: body.HelpID}); err != nil {
		fail(w, "fulfill request", err)
		return
	}
	h.respondWith(w, r, http.StatusOK, id)
}

func (h *handler) respondWith(w http.ResponseWriter, r *http.Request, status, id int) {
	v, err := h.reads.Get(r.Context(), id)
	if err != nil {
		fail(w, "get request", err)
		return
	}
	httpx.WriteJSON(w, status, toResponse(v))
}

var validationErrors = []error{
	domain.ErrMotherNotSpecified, domain.ErrMotherNotFound, domain.ErrInvalidKind,
	domain.ErrInvalidSource, domain.ErrInvalidStatus, domain.ErrTransition,
	domain.ErrCommentRequired, domain.ErrUseFulfill, domain.ErrHelpNotFound,
	domain.ErrHelpMismatch,
}

func fail(w http.ResponseWriter, where string, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	case slices.ContainsFunc(validationErrors, func(e error) bool { return errors.Is(err, e) }):
		httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Println("error", where+":", err)
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
	}
}