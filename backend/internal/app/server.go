package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg Config, deps *DependencyContainer) *Server {
	mux := http.NewServeMux()
	for _, r := range deps.Routers() {
		r.Register(mux)
	}
	// TODO: перенести в контекст document
	mux.Handle("GET /api/photos/{path...}", photosHandler(deps.FilesDir()))

	return &Server{srv: &http.Server{Addr: cfg.HTTPAddr, Handler: mux}}
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