package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gabrieldrouin/shortly/shortener-service/internal/cache"
	"github.com/gabrieldrouin/shortly/shortener-service/internal/repository"
	"github.com/gabrieldrouin/shortly/shortener-service/internal/shortcode"
	"github.com/gabrieldrouin/shortly/shortener-service/internal/validator"
)

const maxRetries = 5

type ShortenHandler struct {
	repo    *repository.URLRepository
	cache   *cache.RedisCache
	baseURL string
}

func NewShortenHandler(repo *repository.URLRepository, cache *cache.RedisCache, baseURL string) *ShortenHandler {
	return &ShortenHandler{repo: repo, cache: cache, baseURL: baseURL}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *ShortenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	if err := validator.ValidateURL(req.URL); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		code, err := shortcode.Generate()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate short code"})
			return
		}

		_, err = h.repo.Insert(r.Context(), code, req.URL)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateShortCode) {
				lastErr = err
				continue
			}
			slog.Error("failed to insert url", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}

		if err := h.cache.SetURL(r.Context(), code, req.URL); err != nil {
			slog.Error("failed to cache url", "error", err, "code", code)
		}

		writeJSON(w, http.StatusCreated, shortenResponse{
			ShortURL: fmt.Sprintf("%s/%s", h.baseURL, code),
		})
		return
	}

	slog.Error("exhausted retries for short code generation", "error", lastErr)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate unique short code"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
