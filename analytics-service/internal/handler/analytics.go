package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gabrieldrouin/shortly/analytics-service/internal/repository"
)

type AnalyticsHandler struct {
	repo *repository.ClickRepository
}

func NewAnalyticsHandler(repo *repository.ClickRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

type analyticsResponse struct {
	ShortCode  string `json:"short_code"`
	ClickCount int64  `json:"click_count"`
}

func (h *AnalyticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	count, err := h.repo.GetClickCount(r.Context(), code)
	if err != nil {
		slog.Error("failed to get click count", "error", err, "code", code)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, analyticsResponse{
		ShortCode:  code,
		ClickCount: count,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
