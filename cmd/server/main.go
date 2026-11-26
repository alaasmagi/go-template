package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alaasmagi/go-template/internal/config"
	"github.com/alaasmagi/go-template/internal/platform/logging"
	"github.com/alaasmagi/go-template/internal/web/api"
)

func main() {
	logger := logging.New()

	if err := config.LoadDotEnv(); err != nil {
		logger.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":" + strconv.Itoa(cfg.AppPort)
	if err := api.Serve(ctx, addr, api.NewRouter(), logger); err != nil {
		logger.Error("http server failed", "error", err)
		os.Exit(1)
	}
}
