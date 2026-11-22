package httpapi

import (
	"errors"
	"log"
	"net/http"
	"slices"
	"time"

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
 return &handler{createCmd: create, registerCmd: register, markCmd: mark, reads: reads}
}

type handler struct {
	createCmd   command.CreateEventHandler
	registerCmd command.RegisterHandler
	markCmd     command.MarkHandler
	reads       query.ReadModel
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
		StartsAt: v.StartsAt.Format("2006-01-02T15:04"),
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
		fail(w, "list event types", err)
		return
	}
	result := make([]eventTypeResponse, 0, len(items))
	for _, t := range items {
		result = append(result, eventTypeResponse{ID: t.ID, Name: t.Name})
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ReadJSON[createBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	startsAt, err := time.Parse("2006-01-02T15:04", body.StartsAt)
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
		fail(w, "create event", err)
		return
	}
	h.respondWith(w, r, http.StatusCreated, id)
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.reads.ListEvents(r.Context())
	if err != nil {
		fail(w, "list events", err)
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
	items, err := h.reads.EventsByMother(r.Context(), motherID)
	if err != nil {
		fail(w, "events by mother", err)
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
	body, err := httpx.ReadJSON[registerBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if err := h.registerCmd.Handle(r.Context(), command.Register{EventID: eventID, MotherID: body.MotherID}); err != nil {
		fail(w, "register participant", err)
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
		fail(w, "list participants", err)
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
	body, err := httpx.ReadJSON[markBody](r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if err := h.markCmd.Handle(r.Context(), command.Mark{EventID: eventID, MotherID: motherID, Status: body.Status}); err != nil {
		fail(w, "mark participant", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handler) respondWith(w http.ResponseWriter, r *http.Request, status, id int) {
	v, err := h.reads.GetEvent(r.Context(), id)
	if err != nil {
		fail(w, "get event", err)
		return
	}
	httpx.WriteJSON(w, status, toResponse(v))
}

var validationErrors = []error{
	domain.ErrTitle, domain.ErrLimit, domain.ErrMotherNotFound,
	domain.ErrAlreadyJoined, domain.ErrFull, domain.ErrNotRegistered,
	domain.ErrInvalidStatus,
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
