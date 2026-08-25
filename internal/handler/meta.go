package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	serviceName        = "url-shortener"
	serviceVersion     = "1.0.0"
	serviceDescription = "URL shortener MVP"
	serviceDocs        = "https://github.com/filipbabicdev/url-shortener"

	healthCheckTimeout = 2 * time.Second
)

// pinger is satisfied by *pgxpool.Pool (Ping(ctx) error). Kept narrow so
// HealthHandler doesn't need to import pgx just to type its dependency.
type pinger interface {
	Ping(ctx context.Context) error
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Error    string `json:"error,omitempty"`
}

func HealthHandler(db pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), healthCheckTimeout)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(healthResponse{
				Status:   "error",
				Database: "down",
				Error:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(healthResponse{
			Status:   "ok",
			Database: "up",
		})
	}
}

type RouteInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Type   string `json:"type"` // "api" | "redirect"
}

type rootResponse struct {
	Service     string      `json:"service"`
	Version     string      `json:"version"`
	Status      string      `json:"status"`
	Environment string      `json:"environment"`
	Description string      `json:"description"`
	Docs        string      `json:"docs"`
	Timestamp   string      `json:"timestamp"`
	Routes      []RouteInfo `json:"routes"`
}

// routeType tells apart the JSON API surface from the redirect route so a
// client of GET / doesn't try to treat /{code} as a REST resource.
func routeType(path string) string {
	if path == "/{code}" {
		return "redirect"
	}
	return "api"
}

// RootHandler closes over r so the route list is walked per-request, not at
// registration time -- routes added after this handler is wired up still
// show up in the response.
func RootHandler(r chi.Router, env string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var routes []RouteInfo
		chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
			routes = append(routes, RouteInfo{
				Method: method,
				Path:   route,
				Type:   routeType(route),
			})
			return nil
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rootResponse{
			Service:     serviceName,
			Version:     serviceVersion,
			Status:      "ok",
			Environment: env,
			Description: serviceDescription,
			Docs:        serviceDocs,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Routes:      routes,
		})
	}
}
