package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/vivangkumar/skv/internal/backend"
	"github.com/vivangkumar/skv/internal/config"
	"github.com/vivangkumar/skv/internal/server"
	"github.com/vivangkumar/skv/internal/store"
)

func main() {
	ctx := signalCtx(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	err := run(ctx)
	if err != nil {
		log.WithError(err).Errorf("stopping skv")
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	configureLog(cfg.Log.Level)

	st := store.New()
	be := backend.New(st)

	s, err := server.New(be, cfg.Server)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	log.Info("starting skv server")

	go s.Listen(ctx)

	<-ctx.Done()
	s.Stop()

	return nil
}

func signalCtx(ctx context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal, len(sig))
		signal.Notify(c, sig...)

		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			cancel()
		}
	}()

	return ctx
}

func configureLog(l log.Level) {
	log.SetLevel(l)
	log.SetFormatter(&log.JSONFormatter{})
}
