package app

import (
	"context"
	"net/http"

	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/router"
	"github.com/ykshvn/reactive-dsr/api-gateway/internal/server"
	"go.uber.org/zap"
)

func Run(ctx context.Context, l *zap.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	r := router.NewRouter(cfg, l)
	srv := server.NewServer(cfg, r)

	go func() {
		<-ctx.Done()
		l.Info("Server is shutting down")
		srv.HTTPServer.Shutdown(ctx)
	}()

	l.Info("Starting server<З")
	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
