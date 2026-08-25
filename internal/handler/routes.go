package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux, handler *URLHandler, db pinger, env string) http.Handler {
	r.Post("/shorten", handler.Shorten)
	r.Get("/health", HealthHandler(db))
	r.Get("/{code}", handler.Redirect)
	r.Get("/", RootHandler(r, env))
	r.NotFound(NoRouteHandler())
	return r
}
