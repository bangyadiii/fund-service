package shutdown

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func Gracefully(ctx context.Context) {
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-notifyCtx.Done()
	log.Println("Shutdown the server")
}
