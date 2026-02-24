package main

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/gabrieldrouin/shortly/redirect-service/internal/cache"
	"github.com/gabrieldrouin/shortly/redirect-service/internal/config"
	"github.com/gabrieldrouin/shortly/redirect-service/internal/handler"
	"github.com/gabrieldrouin/shortly/redirect-service/internal/repository"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Postgres
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to create postgres pool", "error", err)
		return
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping postgres", "error", err)
		return
	}
	slog.Info("connected to postgres")

	// Redis
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to parse redis url", "error", err)
		return
	}
	rdb := redis.NewClient(opts)
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed to ping redis", "error", err)
		return
	}
	slog.Info("connected to redis")

	// Dependencies
	repo := repository.NewURLRepository(pool)
	redisCache := cache.NewRedisCache(rdb)
	redirectHandler := handler.NewRedirectHandler(repo, redisCache)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)
	r.Get("/{code}", redirectHandler.ServeHTTP)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("redirect service listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
