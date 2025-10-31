package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"service-order-avito/internal/config"
	"service-order-avito/internal/http-server/server"
	"service-order-avito/pkg/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.MustInit(cfg.Env)
	log.Info("Logger initialized")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := server.MustInitRouter(ctx, cfg, log)

	srv := &http.Server{Addr: "localhost:" + cfg.Port, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server start up failed: %s", err.Error())
		}
	}()
	log.Info("Listening on " + "localhost:" + cfg.Port)

	gracefulShutdownServer(ctx, cfg, log, srv)
}

func gracefulShutdownServer(ctx context.Context, cfg config.Config, log *slog.Logger, srv *http.Server) {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.IdleTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Info("Server shutdown error: %s\n", err.Error())
	} else {
		log.Info("Server gracefully stopped")
	}
}
