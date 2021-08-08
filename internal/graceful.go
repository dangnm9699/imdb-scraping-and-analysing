package internal

import (
	"adlq/logger"
	"context"
	"os"
	"os/signal"
)

func GracefulShutdown(ctx context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		logger.Debug.Println("shutting down...")
		<-ctx.Done()
	}()
}
