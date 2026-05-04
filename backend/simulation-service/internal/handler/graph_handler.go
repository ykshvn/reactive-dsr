// Package handler used for handling requests
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ykshvn/reactive-dsr/simulation-service/internal/graph"
)

type GraphHandler struct {
	generator *graph.Generator
}

func NewGraphHandler() *GraphHandler {
	return &GraphHandler{
		generator: graph.NewGenerator(),
	}
}

func (h *GraphHandler) Generate(w http.ResponseWriter, r *http.Request) {
	n := 10

	resp, err := h.generator.GenerateGraph(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
