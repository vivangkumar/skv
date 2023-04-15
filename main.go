package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/vivangkumar/skv/pkg/backend"
	"github.com/vivangkumar/skv/pkg/config"
	"github.com/vivangkumar/skv/pkg/node"
	"github.com/vivangkumar/skv/pkg/store"
)

func main() {
	ctx := signalCtx(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	err := run(ctx)
	if err != nil {
		log.Errorf("stopping skv")
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	configureLog(cfg.Log.Level)

	st := store.NewStore()
	be := backend.NewBackend(st)

	n, err := node.NewNode(be, cfg.Node)
	if err != nil {
		return fmt.Errorf("node: %w", err)
	}

	log.WithField("node_id", n.ID()).Info("started node")

	go n.Listen(ctx)

	<-ctx.Done()
	n.Stop()

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
