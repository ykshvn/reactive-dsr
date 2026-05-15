// Package app required for preparing and starting logic
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ykshvn/reactive-dsr/simulation-service/internal/handler"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

func Run(ctx context.Context, l *zap.Logger) error {
	wsHub := ws.NewHub(l)
	graphHandler := handler.NewGraphHandler(wsHub)
	wsHandler := handler.NewWSHandler(wsHub)

	srv := &http.Server{
		Addr:         ":6970",
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/graph/generate", graphHandler.Generate)
	http.Handle("/ws/simulation", wsHandler)

	go func() {
		<-ctx.Done()
		l.Info("Server is shutting down")
		srv.Shutdown(ctx)
	}()

	l.Info("simulation service started on port 6970")

	go wsHub.Run()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
