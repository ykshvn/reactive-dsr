package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ykshvn/reactive-dsr/simulation-service/internal/simulation"
)

type SimulationHandler struct {
	engine *simulation.Engine
}

func NewSimulationHandler(e *simulation.Engine) *SimulationHandler {
	return &SimulationHandler{
		engine: e,
	}
}

func (h *SimulationHandler) StartRouteDiscovery(w http.ResponseWriter, r *http.Request) {
	srcStr := r.URL.Query().Get("src")
	dstStr := r.URL.Query().Get("dst")

	if srcStr == "" || dstStr == "" {
		http.Error(w, `{"error": "source and dest parameters are required"}`, http.StatusBadRequest)
		return
	}

	src, err1 := strconv.Atoi(srcStr)
	dst, err2 := strconv.Atoi(dstStr)

	if err1 != nil || err2 != nil {
		http.Error(w, `{"error": "invalid source or dest"}`, http.StatusBadRequest)
		return
	}

	h.engine.StartStepRouteDiscovery(src, dst)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Route discovery started",
		"source":  src,
		"dest":    dst,
	})
}

func (h *SimulationHandler) Step(w http.ResponseWriter, r *http.Request) {
	event, err := h.engine.Step()
	if err != nil {
		http.Error(w, `{"error": "step failed"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":      "ok",
		"step":        h.engine.GetCurrentStep(),
		"queueLength": h.engine.GetQueueLength(),
	}

	if event.Step != -1 {
		response["event"] = event
	} else {
		response["message"] = "No events to process"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
