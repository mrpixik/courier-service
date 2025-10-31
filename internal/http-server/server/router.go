package router

import (
	"github.com/go-chi/chi/v5"
	"service-order-avito/internal/config"
)

func MustInit(cfd config.Config) chi.Router {
	router := chi.NewRouter()

	router.Get("/ping", pingGetHandler)

}
