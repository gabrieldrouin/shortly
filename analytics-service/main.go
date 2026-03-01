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

	"github.com/gabrieldrouin/shortly/analytics-service/internal/config"
	"github.com/gabrieldrouin/shortly/analytics-service/internal/handler"
	"github.com/gabrieldrouin/shortly/analytics-service/internal/kafka"
	"github.com/gabrieldrouin/shortly/analytics-service/internal/repository"
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

	// Dependencies
	repo := repository.NewClickRepository(pool)
	analyticsHandler := handler.NewAnalyticsHandler(repo)

	// Kafka consumer
	consumer := kafka.NewConsumer(cfg.KafkaBroker, repo)
	defer consumer.Close()
	go consumer.Run(ctx)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handler.Health)
	r.Get("/api/analytics/{code}", analyticsHandler.ServeHTTP)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("analytics service listening", "port", cfg.Port)
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
