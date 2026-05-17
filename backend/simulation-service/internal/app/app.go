// Package app required for preparing and starting logic
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/ykshvn/reactive-dsr/simulation-service/internal/handler"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/simulation"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
	"go.uber.org/zap"
)

type App struct {
	log          *zap.Logger
	engine       *simulation.Engine
	wsHub        *ws.Hub
	graphHandler *handler.GraphHandler
	wsHandler    *handler.WSHandler
	simHandler   *handler.SimulationHandler
}

func NewApp(l *zap.Logger) *App {
	hub := ws.NewHub(l)
	engine := simulation.NewEngine(hub, l)
	wsHandler := handler.NewWSHandler(hub)
	simHandler := handler.NewSimulationHandler(engine)

	return &App{
		log:          l,
		engine:       engine,
		wsHub:        hub,
		graphHandler: handler.NewGraphHandler(hub, engine),
		wsHandler:    wsHandler,
		simHandler:   simHandler,
	}
}

func (a *App) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         ":6970",
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/graph/generate", a.graphHandler.Generate)
	http.HandleFunc("/simulation/start", a.simHandler.StartRouteDiscovery)
	http.Handle("/ws/simulation", a.wsHandler)

	go func() {
		<-ctx.Done()
		a.log.Info("Server is shutting down")
		srv.Shutdown(ctx)
	}()

	a.log.Info("simulation service started on port 6970")

	go a.wsHub.Run()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
