package server

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"service-order-avito/internal/config"
	"service-order-avito/internal/http/middleware"
	"service-order-avito/internal/http/server/handlers"
)

type CourierHandler interface {
	Post(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	GetAll(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
}

func InitRouter(cfg config.HTTPServer, log *slog.Logger, courierHandler CourierHandler) chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.WithLogger(log),
	)

	router.Get("/ping", handlers.PingGetHandler)
	router.Head("/healthcheck", handlers.HealthcheckHeadHandler)

	router.Get("/couriers", courierHandler.GetAll)

	router.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", courierHandler.Get)
		r.Post("/", courierHandler.Post)
		r.Put("/", courierHandler.Put)

	})
	return router
}
