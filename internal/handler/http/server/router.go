package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
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

func InitRouter(log *slog.Logger,
	courierHandler courierHandler,
	deliveryHandler deliveryHandler,
	metricObserver middleware.MetricsObserverHTTP) chi.Router {

	router := chi.NewRouter()

	router.Use(
		middleware.WithMonitoring(log, metricObserver),
	)

	router.Handle("/metrics", promhttp.Handler())

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
