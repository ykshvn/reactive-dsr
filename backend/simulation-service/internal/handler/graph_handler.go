// Package handler used for handling requests
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ykshvn/reactive-dsr/shared/events"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/graph"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/ws"
)

type GraphHandler struct {
	generator *graph.Generator
	hub       *ws.Hub
}

func NewGraphHandler(hub *ws.Hub) *GraphHandler {
	return &GraphHandler{
		generator: graph.NewGenerator(),
		hub:       hub,
	}
}

func (h *GraphHandler) Generate(w http.ResponseWriter, r *http.Request) {
	nodes := 10
	nodesStr := r.URL.Query().Get("nodes")
	if nodesStr != "" {
		if n, err := strconv.Atoi(nodesStr); err == nil {
			if n >= 3 && n <= 50 {
				nodes = n
			} else if n > 50 {
				nodes = 50
			} else if n < 3 {
				nodes = 3
			}
		}
	}

	resp, err := h.generator.GenerateGraph(nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Test PREQ Message
	if false {
		h.hub.BroadcastEvent(events.NewEvent(events.EventGraphGenerated, resp))
		testRREQ := events.RREQPayload{
			From:       0,
			To:         nodes - 1,
			RouteSoFar: []int{0, 1, 3},
			RequestID:  42,
		}

		h.hub.BroadcastEvent(events.NewEvent(events.EventRREQPropagated, testRREQ))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
