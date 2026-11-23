package httpx

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

const APITime = time.RFC3339

const MaxBodySize = 1 << 20 // 1 MB

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

func ReadJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func PathID(r *http.Request, name string) (int, error) {
	return strconv.Atoi(r.PathValue(name))
}

type ErrorMapper struct {
	notFound     []error
	validation   []error
	unauthorized []error
	conflict     []error
	forbidden    []error
}

type ErrorMapperOption func(*ErrorMapper)

func WithNotFound(errs ...error) ErrorMapperOption {
	return func(m *ErrorMapper) { m.notFound = append(m.notFound, errs...) }
}

func WithValidation(errs ...error) ErrorMapperOption {
	return func(m *ErrorMapper) { m.validation = append(m.validation, errs...) }
}

func WithUnauthorized(errs ...error) ErrorMapperOption {
	return func(m *ErrorMapper) { m.unauthorized = append(m.unauthorized, errs...) }
}

func WithConflict(errs ...error) ErrorMapperOption {
	return func(m *ErrorMapper) { m.conflict = append(m.conflict, errs...) }
}

func WithForbidden(errs ...error) ErrorMapperOption {
	return func(m *ErrorMapper) { m.forbidden = append(m.forbidden, errs...) }
}

func NewErrorMapper(opts ...ErrorMapperOption) *ErrorMapper {
	m := &ErrorMapper{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *ErrorMapper) Handle(w http.ResponseWriter, where string, err error) {
	if err == nil {
		return
	}
	status, ok := m.status(err)
	if !ok {
		log.Println("error "+where+":", err)
		WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	WriteError(w, status, err.Error())
}

func (m *ErrorMapper) status(err error) (int, bool) {
	for _, e := range m.notFound {
		if errors.Is(err, e) {
			return http.StatusNotFound, true
		}
	}
	for _, e := range m.validation {
		if errors.Is(err, e) {
			return http.StatusUnprocessableEntity, true
		}
	}
	for _, e := range m.unauthorized {
		if errors.Is(err, e) {
			return http.StatusUnauthorized, true
		}
	}
	for _, e := range m.conflict {
		if errors.Is(err, e) {
			return http.StatusConflict, true
		}
	}
	for _, e := range m.forbidden {
		if errors.Is(err, e) {
			return http.StatusForbidden, true
		}
	}
	return 0, false
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Println("panic recovered:", recovered)
				WriteError(w, http.StatusInternalServerError, "something went wrong")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(lw, r)
		log.Printf("%s %s %d %s", r.Method, r.RequestURI, lw.status, time.Since(start).Round(time.Millisecond))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lw *loggingResponseWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}

func IsMaxBytesExceeded(err error) bool {
	var maxErr *http.MaxBytesError
	return errors.As(err, &maxErr)
}