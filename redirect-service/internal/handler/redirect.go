package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gabrieldrouin/shortly/redirect-service/internal/cache"
	"github.com/gabrieldrouin/shortly/redirect-service/internal/repository"
)

type RedirectHandler struct {
	repo  *repository.URLRepository
	cache *cache.RedisCache
}

func NewRedirectHandler(repo *repository.URLRepository, cache *cache.RedisCache) *RedirectHandler {
	return &RedirectHandler{repo: repo, cache: cache}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	// 1. Check Redis cache
	cached, err := h.cache.GetURL(r.Context(), code)
	if err != nil {
		slog.Error("cache lookup failed", "error", err, "code", code)
	}
	if cached != "" {
		http.Redirect(w, r, cached, http.StatusFound)
		return
	}

	// 2. Cache miss — query DB
	u, err := h.repo.GetByShortCode(r.Context(), code)
	if err != nil {
		slog.Error("db lookup failed", "error", err, "code", code)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	if u == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	// 3. Check expiry
	if u.ExpiresAt != nil && u.ExpiresAt.Before(time.Now()) {
		if err := h.cache.DeleteURL(r.Context(), code); err != nil {
			slog.Error("failed to delete expired cache entry", "error", err, "code", code)
		}
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	// 4. Backfill cache (best-effort)
	if err := h.cache.SetURL(r.Context(), code, u.OriginalURL); err != nil {
		slog.Error("failed to backfill cache", "error", err, "code", code)
	}

	http.Redirect(w, r, u.OriginalURL, http.StatusFound)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
