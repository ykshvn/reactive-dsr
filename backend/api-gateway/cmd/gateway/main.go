package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ykshvn/reactive-dsr/api-gateway/internal/app"
	"github.com/ykshvn/reactive-dsr/shared/logger"
	"go.uber.org/zap"
)

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("[ERROR]: %v", err)
	}
	defer l.Sync()

	if err := realMain(l); err != nil {
		l.Error(err.Error())
	}
}

func realMain(l *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	return app.Run(ctx, l)
}
