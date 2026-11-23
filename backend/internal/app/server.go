package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"backend/internal/auth"
	"backend/internal/shared/httpx"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg Config, deps *DependencyContainer) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.Handle("GET /api/photos/{path...}", photosHandler(deps.FilesDir()))

	apiMux := http.NewServeMux()
	for _, r := range deps.Routers() {
		r.Register(apiMux)
	}

	apiHandler := auth.Middleware([]byte(cfg.AuthSecret))(apiMux)
	apiHandler = httpx.LoggingMiddleware(apiHandler)
	apiHandler = httpx.RecoverMiddleware(apiHandler)

	mux.Handle("/", apiHandler)

	return &Server{srv: &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Println("сервер запущен на", s.srv.Addr)
		errCh <- s.srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}

func photosHandler(filesDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("path")
		if fileName == "" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(filesDir, fileName))
	})
}