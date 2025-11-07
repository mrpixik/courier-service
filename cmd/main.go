package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"service-order-avito/internal/config"
	"service-order-avito/internal/http/server"
	"service-order-avito/internal/http/server/handlers"
	"service-order-avito/internal/repository/postgres"
	"service-order-avito/internal/service"
	"service-order-avito/pkg/logger"
	"syscall"
)

func main() {
	// CONFIG
	cfg := config.MustLoad()

	// LOGGER
	log := logger.MustInit(cfg.Env)
	log.Info("Logger initialized")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Repository's Lay
	courierRepository, err := postgres.NewCourierRepositoryPostgres(ctx, cfg)
	if err != nil {
		log.Error("Failed to init repository: " + err.Error())
		os.Exit(1)
	}
	log.Info("Courier repository postgres initialized")

	// Service lay
	courierService := service.NewCourierService(courierRepository)
	log.Info("Courier service initialized")

	// Controller's lay
	courierHandler := handlers.NewCourierHandler(courierService)
	log.Info("Courier handler initialized")

	// ROUTER & SERVER
	r := server.InitRouter(ctx, cfg, log, courierHandler)

	srv := &http.Server{Addr: "localhost:" + cfg.ServerPort, Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server start up failed: %s", err.Error())
		}
	}()
	log.Info("Listening on " + "localhost:" + cfg.ServerPort)

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
