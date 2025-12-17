package server

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"service-order-avito/internal/config"
	"service-order-avito/internal/handler/http/middleware"
	"service-order-avito/internal/handler/http/server/handler"
)

type courierHandler interface {
	Post(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
	GetAll(http.ResponseWriter, *http.Request)
	Put(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type deliveryHandler interface {
	PostAssign(http.ResponseWriter, *http.Request)
	PostUnassign(http.ResponseWriter, *http.Request)
}

func InitRouter(cfg config.HTTPServer, log *slog.Logger, courierHandler courierHandler, deliveryHandler deliveryHandler) chi.Router {
	router := chi.NewRouter()

	router.Use(
		middleware.WithLogger(log),
	)

	router.Get("/ping", handler.PingGetHandler)
	router.Head("/healthcheck", handler.HealthcheckHeadHandler)

	router.Get("/couriers", courierHandler.GetAll)

	router.Route("/courier", func(r chi.Router) {
		r.Get("/{id}", courierHandler.Get)
		r.Post("/", courierHandler.Post)
		r.Put("/", courierHandler.Put)
		r.Delete("/{id}", courierHandler.Delete)
	})

	router.Route("/delivery", func(r chi.Router) {
		r.Post("/assign", deliveryHandler.PostAssign)
		r.Post("/unassign", deliveryHandler.PostUnassign)
	})
	return router
}
