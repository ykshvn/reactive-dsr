package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ykshvn/reactive-dsr/simulation-service/internal/handler"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

// TODO: Refactor this to different package
func realMain() error {
	graphHandler := handler.NewGraphHandler()

	// TODO: Process query params
	http.HandleFunc("/graph/generate", graphHandler.Generate)

	fmt.Println("Simulation Service started on http://localhost:6970")

	if err := http.ListenAndServe(":6970", nil); err != nil {
		return err
	}
	return nil
}
