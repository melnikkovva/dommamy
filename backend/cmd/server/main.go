package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := app.LoadConfig()
	if err != nil {
		return err
	}

	conns, err := app.NewConnectionContainer(ctx, cfg)
	if err != nil {
		return err
	}
	defer conns.Close()

	deps, err := app.NewDependencyContainer(cfg, conns)
	if err != nil {
		return err
	}

	go deps.LearnSyncWorker().Run(ctx)

	return app.NewServer(cfg, deps).Run(ctx)
}