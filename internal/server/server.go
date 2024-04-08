package server

import (
	"avito-test2024-spring/internal/config"
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.HTTPConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
