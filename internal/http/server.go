package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	HTTP *http.Server
}

func (s *Server) Start(addr string, handler http.Handler) {
	s.HTTP = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("server starting", "addr", addr)
	if err := s.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server listen error", "err", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	slog.Info("shutting down server")
	if err := s.HTTP.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "err", err)
	}
	slog.Info("server stopped")
}
