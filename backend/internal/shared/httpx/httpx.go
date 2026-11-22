package httpx

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Router interface {
	Register(mux *http.ServeMux)
}

type errorBody struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("write json:", err)
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, errorBody{Error: msg})
}

func ReadJSON[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func PathID(r *http.Request, name string) (int, error) {
	return strconv.Atoi(r.PathValue(name))
}