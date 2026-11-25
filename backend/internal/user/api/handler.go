package api

import (
	"errors"
	"net/http"

	"backend/internal/auth"
	"backend/internal/shared/httpx"
	"backend/internal/user"
)

type Handler struct {
	s *user.Service
}

func New(s *user.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", h.login)
	mux.HandleFunc("POST /api/auth/me", h.me)
	mux.HandleFunc("POST /api/users", h.createUser)
	mux.HandleFunc("GET /api/users", h.listUsers)
	mux.HandleFunc("GET /api/users/{id}", h.getUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.deleteUser)
}

type loginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token     string `json:"token"`
	Id        int    `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	MotherId  int    `json:"mother_id"`
}

type userResponse struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	MotherId  int    `json:"mother_id"`
	CreatedAt string `json:"created_at"`
}

type createUserBody struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	MotherId  int    `json:"mother_id"`
}

func toUserResponse(u *user.User) userResponse {
	return userResponse{
		Id: u.Id, Login: u.Login, FirstName: u.FirstName, LastName: u.LastName,
		Role: u.Role, MotherId: u.MotherId, CreatedAt: u.CreatedAt.Format(httpx.APITime),
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	req, err := httpx.ReadJSON[loginBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	token, u, err := h.s.Login(req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorInvalidCredentials):
			httpx.WriteError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusOK, loginResponse{
		Token:     token,
		Id:        u.Id,
		Login:     u.Login,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		MotherId:  u.MotherId,
	})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	httpx.WriteError(w, http.StatusOK, "use login endpoint")
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.FromContext(r.Context())
	if err != nil || claims.Role != user.RoleAdmin {
		httpx.WriteError(w, http.StatusForbidden, "только администратор может создавать пользователей")
		return
	}
	req, err := httpx.ReadJSON[createUserBody](w, r)
	if err != nil {
		if httpx.IsMaxBytesExceeded(err) {
			httpx.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		httpx.WriteError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}
	if req.Role == "" {
		req.Role = user.RoleUser
	}
	u, err := h.s.Create(req.Login, req.Password, req.FirstName, req.LastName, req.Role, req.MotherId)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorLoginTaken):
			httpx.WriteError(w, http.StatusConflict, "Логин уже занят")
		case errors.Is(err, user.ErrorLogin), errors.Is(err, user.ErrorPassword),
			errors.Is(err, user.ErrorName), errors.Is(err, user.ErrorSurname):
			httpx.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.FromContext(r.Context())
	if err != nil || claims.Role != user.RoleAdmin {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}
	items, err := h.s.ListUsers()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	result := []userResponse{}
	for _, u := range items {
		result = append(result, toUserResponse(u))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	claims, ctxErr := auth.FromContext(r.Context())
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if ctxErr != nil || (claims.Role != user.RoleAdmin && claims.UserID != id) {
		httpx.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}
	u, err := h.s.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorNotFound):
			httpx.WriteError(w, http.StatusNotFound, "user not found")
		default:
			httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ctxErr := auth.FromContext(r.Context())
	if ctxErr != nil || claims.Role != user.RoleAdmin {
		httpx.WriteError(w, http.StatusForbidden, "только администратор может удалять пользователей")
		return
	}
	id, err := httpx.PathID(r, "id")
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "incorrect id")
		return
	}
	if claims.UserID == id {
		httpx.WriteError(w, http.StatusConflict, "нельзя удалить свою учетную запись")
		return
	}
	if err := h.s.DeleteUser(id); err != nil {
		if errors.Is(err, user.ErrorNotFound) {
			httpx.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		httpx.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
