package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/filipbabicdev/url-shortener/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type URLHandler struct {
	repo *repository.URLRepository
}

func NewURLHandler(repo *repository.URLRepository) *URLHandler {
	return &URLHandler{repo: repo}
}

// POST /shorten
func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	type request struct {
		URL string `json:"url"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if parsedURL, err := url.ParseRequestURI(req.URL); err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		writeJSONError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	parsedURL, err := h.repo.Create(r.Context(), req.URL)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to shorten URL")
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

// GET /{code}
func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.repo.GetByShortCode(r.Context(), code)
	if err == pgx.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "URL not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to retrieve URL")
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
