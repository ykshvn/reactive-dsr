package main

import (
	"fmt"
	"net/http"
)

// NOTE: this is for reverse proxy testing
func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service": "simulation-service", "message": "Hello from simulation!"}`)
	})

	http.HandleFunc("/graph/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "nodes": 15, "edges": 42}`)
	})

	fmt.Println("Simulation service started on localhost:6970")

	http.ListenAndServe(":6970", nil)
}
