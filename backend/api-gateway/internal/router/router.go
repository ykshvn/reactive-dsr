package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/middleware"
	"go.uber.org/zap"
)

func NewRouter(cfg *config.Config, l *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	// r.Use(middleware.Recoverer)
	// r.Use(middleware.RealIP)
	// r.Use(middleware.RequestID)

	// Middleware Chain
	r.Use(middleware.Chain(
		middleware.Recovery(l),
		middleware.RequestID,
		middleware.Logging(l),
		middleware.CORS(cfg.Cors.AllowedOrigins),
		// TODO: rate limiter
	))

	// Routes

	// CheckHeath
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// Api...
	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`{"message": "api-gateway is alive"}`))
		})
	})

	// WebSocket
	r.Route("/ws/simulation", func(r chi.Router) {
		// TODO: WS Endpoint Implementation
	})

	return r
}
