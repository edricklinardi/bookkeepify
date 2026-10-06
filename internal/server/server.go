// Package server is the HTTP API. It reads only from Postgres and never calls
// Spotify or Last.fm on the request path, except for the OAuth exchange.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/edricklinardi/bookkeepify/internal/config"
)

func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", cfg.Addr)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})
	return mux
}
