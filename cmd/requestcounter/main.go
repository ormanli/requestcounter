package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/ormanli/requestcounter/internal"
	"github.com/ormanli/requestcounter/internal/app/requestcounter"
)

func main() {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	ctx, cncl := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cncl()

	var c requestcounter.Config

	err := envconfig.Process("app", &c)
	if err != nil {
		slog.Error("", "error", err.Error())
		code = 1
		return
	}

	err = internal.Run(ctx, c)
	if err != nil {
		slog.Error("", "error", err.Error())
		code = 1
		return
	}
}
