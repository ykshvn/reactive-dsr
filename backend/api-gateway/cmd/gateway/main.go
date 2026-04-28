package main

import (
	"fmt"
	"log"

	"github.com/ykshvn/reactive-dsr/api-gateway/internal/config"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatalf("[ERROR]: %v", err)
	}
}

func realMain() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println(cfg)

	return nil
}
