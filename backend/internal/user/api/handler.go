package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"backend/internal/user"
	"backend/internal/shared/httpx"
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
}

type loginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token       string `json:"token"`
	Id          int    `json:"id"`
	Login       string `json:"login"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
	MotherId    int    `json:"mother_id"`
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
		Role: u.Role, MotherId: u.MotherId, CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	token, u, err := h.s.Login(req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorInvalidCredentials):
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный логин или пароль"})
		default:
			log.Println("error login:", err)
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
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
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "use login endpoint"})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect JSON"})
		return
	}
	if req.Role == "" {
		req.Role = user.RoleUser
	}
	u, err := h.s.Create(req.Login, req.Password, req.FirstName, req.LastName, req.Role, req.MotherId)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorLoginTaken):
			httpx.WriteJSON(w, http.StatusConflict, map[string]string{"error": "Логин уже занят"})
		case errors.Is(err, user.ErrorLogin), errors.Is(err, user.ErrorPassword),
			errors.Is(err, user.ErrorName), errors.Is(err, user.ErrorSurname):
			httpx.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		default:
			log.Println("error create user:", err)
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	items, err := h.s.ListUsers()
	if err != nil {
		log.Println("error list users:", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}
	result := []userResponse{}
	for _, u := range items {
		result = append(result, toUserResponse(u))
	}
	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "incorrect id"})
		return
	}
	u, err := h.s.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrorNotFound):
			httpx.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			log.Println("error get user:", err)
			httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}
	httpx.WriteJSON(w, http.StatusOK, toUserResponse(u))
}
