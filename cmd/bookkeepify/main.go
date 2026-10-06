// Command bookkeepify runs one of the Bookkeepify processes. All three share
// one binary and one Postgres database; they split into separate deployments
// only in M6.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/edricklinardi/bookkeepify/internal/config"
	"github.com/edricklinardi/bookkeepify/internal/server"
)

const usage = `usage: bookkeepify <command>

commands:
  api        serve the HTTP API
  worker     run background jobs (sync, enrich, route, apply)
  scheduler  enqueue periodic jobs`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cmd := os.Args[1]
	logger = logger.With("cmd", cmd)

	switch cmd {
	case "api":
		err = server.Run(ctx, cfg, logger)
	case "worker", "scheduler":
		// Implemented in M1 (worker: SyncLikes) and M5 (scheduler).
		logger.Info("not implemented yet; waiting for shutdown signal")
		<-ctx.Done()
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		logger.Error("exit", "err", err)
		os.Exit(1)
	}
}
