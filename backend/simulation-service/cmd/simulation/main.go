package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ykshvn/reactive-dsr/shared/logger"
	"github.com/ykshvn/reactive-dsr/simulation-service/internal/app"
	"go.uber.org/zap"
)

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	if err := realMain(l); err != nil {
		l.Error(err.Error())
	}
}

// TODO: Refactor this to different package
func realMain(l *zap.Logger) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	return app.Run(ctx, l)
}
