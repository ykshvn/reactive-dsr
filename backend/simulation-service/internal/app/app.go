// Package app required for preparing and starting logic
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ykshvn/reactive-dsr/shared/events"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/handler"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

func Run(ctx context.Context, l *zap.Logger) error {
	graphHandler := handler.NewGraphHandler()
	wsHub := ws.NewHub(l)
	wsHandler := handler.NewWSHandler(wsHub)

	srv := &http.Server{
		Addr:         ":6970",
		Handler:      nil, // будем регистрировать вручную
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// TODO: Process query params
	http.HandleFunc("/graph/generate", graphHandler.Generate)
	http.Handle("/ws/simulation", wsHandler)

	go func() {
		<-ctx.Done()
		l.Info("Server is shutting down")
		srv.Shutdown(ctx)
	}()

	l.Info("simulation service started on port 6970")

	go wsHub.Run()

	go func() {
		step := 0
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				step++
				testEvent := events.NewStepEvent(events.EventSimulationStep, step, map[string]interface{}{
					"message": "Test event from simulation service",
					"step":    step,
				})
				wsHub.BroadcastEvent(testEvent)
				l.Info("Broadcasted test event", zap.Int("step", step))
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
