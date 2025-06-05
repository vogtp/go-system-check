package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/vogtp/go-system-check/cmd/root"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	root.Command(ctx)
}
