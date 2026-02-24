package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
}

func Load() Config {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	return Config{
		Port:        port,
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", user, pass, db),
		RedisURL:    redisURL,
	}
}
