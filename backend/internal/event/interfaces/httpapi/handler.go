package httpapi

import (
	"net/http"
	"time"

	"backend/internal/auth"
	"backend/internal/event/application/command"
	"backend/internal/event/application/query"
	"backend/internal/event/domain"
	"backend/internal/shared/httpx"
)

func NewHandler(
	create command.CreateEventHandler,
	register command.RegisterHandler,
	mark command.MarkHandler,
	reads query.ReadModel,
) httpx.Router {
	return &handler{
		createCmd: create, registerCmd: register, markCmd: mark, reads: reads,
		mapper: httpx.NewErrorMapper(
			httpx.WithNotFound(domain.ErrNotFound),
			httpx.WithValidation(
				domain.ErrTitle, domain.ErrLimit, domain.ErrMotherNotFound,
				domain.ErrAlreadyJoined, domain.ErrFull, domain.ErrNotRegistered,
				domain.ErrInvalidStatus, domain.ErrInvalidTransition,
			),
		),
	}
}

type handler struct {
	createCmd   command.CreateEventHandler
	registerCmd command.RegisterHandler
	markCmd     command.MarkHandler
	reads       query.ReadModel
	mapper      *httpx.ErrorMapper
}

func (h *handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/event-types", h.types)
	mux.HandleFunc("POST /api/events", h.create)
	mux.HandleFunc("GET /api/events", h.list)
	mux.HandleFunc("POST /api/events/{id}/participants", h.register)
	mux.HandleFunc("GET /api/events/{id}/participants", h.participants)
	mux.HandleFunc("PUT /api/events/{id}/participants/{motherId}", h.mark)
	mux.HandleFunc("GET /api/mother/{id}/events", h.byMother)
}

type eventResponse struct {
	ID        int    `json:"id"`
	TypeID    int    `json:"type_id"`
	Title     string `json:"title"`
	StartsAt  string `json:"starts_at"`
	Place     string `json:"place"`
	MaxPeople int    `json:"max_people"`
}

func toResponse(v query.EventView) eventResponse {
	return eventResponse{
		ID: v.ID, TypeID: v.TypeID, Title: v.Title,
		StartsAt: v.StartsAt.Format(httpx.APITime),
		Place: v.Place, MaxPeople: v.MaxPeople,
	}
}

func toResponses(items []query.EventView) []eventResponse {
	result := make([]eventResponse, 0, len(items))
	for _, v := range items {
		result = append(result, toResponse(v))
	}
	return result
}

type eventTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type participationResponse struct {
	MotherID int                       `json:"mother_id"`
	Status   domain.ParticipationStatus `json:"status"`
}

type createBody struct {
	TypeID    int    `json:"type_id"`
	Title     string `json:"title"`
	StartsAt  string `json:"starts_at"`
	Place     string `json:"place"`
	MaxPeople int    `json:"max_people"`
}

type registerBody struct {
	MotherID int `json:"mother_id"`
}

type markBody struct {
	Status string `json:"status"`
}

func (h *handler) types(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.ListEventTypes(r.Context())
	if err != nil {
		h.mapper.Handle(w, "list event types", err)
		return
	}
	result := make([]eventTypeResponse, 0, len(items))
	for _, t := range items {
		result = append(result, eventTypeResponse{ID: t.ID, Name: t.Name})
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ReadJSON[createBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	startsAt, err := time.Parse(httpx.APITime, body.StartsAt)
	if err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, "specify the date and time")
		return
	}
	id, err := h.createCmd.Handle(r.Context(), command.CreateEvent{
		TypeID:    body.TypeID,
		Title:     body.Title,
		StartsAt:  startsAt,
		Place:     body.Place,
		MaxPeople: body.MaxPeople,
	})
	if err != nil {
		h.mapper.Handle(w, "create event", err)
		return
	}
	h.respondWith(w, r, http.StatusCreated, id)
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.ListEvents(r.Context())
	if err != nil {
		h.mapper.Handle(w, "list events", err)
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
	if !auth.RequireOwnershipOrAdmin(w, r, motherID) {
		return
	}
	items, err := h.reads.EventsByMother(r.Context(), motherID)
	if err != nil {
		h.mapper.Handle(w, "events by mother", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toResponses(items))
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	eventID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[registerBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if err := h.registerCmd.Handle(r.Context(), command.Register{EventID: eventID, MotherID: body.MotherID}); err != nil {
		h.mapper.Handle(w, "register participant", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (h *handler) participants(w http.ResponseWriter, r *http.Request) {
	eventID, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	items, err := h.reads.Participants(r.Context(), eventID)
	if err != nil {
		h.mapper.Handle(w, "list participants", err)
		return
	}
	result := make([]participationResponse, 0, len(items))
	for _, p := range items {
		result = append(result, participationResponse{MotherID: p.MotherID, Status: p.Status})
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *handler) mark(w http.ResponseWriter, r *http.Request) {
	eventID, err1 := httpx.PathID(r, "id")
	motherID, err2 := httpx.PathID(r, "motherId")
	if err1 != nil || err2 != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	body, err := httpx.ReadJSON[markBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if err := h.markCmd.Handle(r.Context(), command.Mark{EventID: eventID, MotherID: motherID, Status: body.Status}); err != nil {
		h.mapper.Handle(w, "mark participant", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handler) respondWith(w http.ResponseWriter, r *http.Request, status, id int) {
	v, err := h.reads.GetEvent(r.Context(), id)
	if err != nil {
		h.mapper.Handle(w, "get event", err)
		return
	}
	httpx.WriteJSON(w, status, toResponse(v))
}
