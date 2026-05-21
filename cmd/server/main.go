package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/filipbabicdev/url-shortener/internal/config"
	"github.com/filipbabicdev/url-shortener/internal/database"
	"github.com/filipbabicdev/url-shortener/internal/handler"
	"github.com/filipbabicdev/url-shortener/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPool, err := database.NewPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer dbPool.Close()

	urlRepo := repository.NewURLRepository(dbPool)
	urlHandler := handler.NewURLHandler(urlRepo)

	r := chi.NewRouter()
	handler.SetupRoutes(r, urlHandler)

	log.Printf("Starting server on port %s...", cfg.ServerPort)
	server := http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// TODO: graceful shutdown (handle http.ErrServerClosed, signal.NotifyContext)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
