package handler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/filipbabicdev/url-shortener/internal/repository"
)

type URLHandler struct {
	repo *repository.URLRepository
}

func NewURLHandler(repo *repository.URLRepository) *URLHandler {
	return &URLHandler{repo: repo}
}

//[ ] `POST /shorten`
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	type request struct {
		URL string `json:"url"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if parsedURL, err := url.ParseRequestURI(req.URL); err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	parsedURL, err := h.repo.Create(r.Context(), req.URL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	response := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: fmt.Sprintf("http://%s/%s", r.Host, parsedURL.ShortCode),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response) 
}
//[ ] `GET /{code}` — redirect (301/302) na originalni URL
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Implementacija će biti dodana kasnije.
}
//- [ ] `GET /health` — health check
func (h *URLHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}