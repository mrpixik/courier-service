package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"service-order-avito/internal/config"
	"service-order-avito/internal/http-server/middleware"
)

func MustInitRouter(ctx context.Context, cfg config.Config, log *slog.Logger) chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.WithGracefulShutdown(ctx),
		middleware.WithLogger(log),
	)

	router.Get("/ping", pingGetHandler)
	router.Head("/healthcheck", healthcheckHeadHandler)

	return router
}
