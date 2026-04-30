// Package server used for creating and starting http server
package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(cfg *config.Config, r *chi.Mux) *Server {
	return &Server{
		HTTPServer: &http.Server{
			Addr:    ":" + fmt.Sprintf("%d", cfg.Server.Port),
			Handler: r,

			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	return s.HTTPServer.ListenAndServe()
}
