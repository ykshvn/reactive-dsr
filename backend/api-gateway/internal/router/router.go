package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/middleware"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/proxy"
	"go.uber.org/zap"
)

func NewRouter(cfg *config.Config, l *zap.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()

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

	// Proxy
	httpProxy, err := proxy.NewHTTPProxy(cfg, l)
	if err != nil {
		return nil, err
	}
	r.Route("/api", func(r chi.Router) {
		r.Handle("/*", httpProxy.ProxyHandler())
	})

	// WebSocket
	ws := proxy.NewWSProxy(l)
	r.Get("/ws/simulation", ws.ServeHTTP)

	return r, nil
}
