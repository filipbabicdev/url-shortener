package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(handler *URLHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/shorten", handler.Shorten)
	r.Get("/health", handler.HealthCheck)
	r.Get("/{code}", handler.Redirect)
	return r
}