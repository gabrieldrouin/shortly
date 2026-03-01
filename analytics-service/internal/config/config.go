package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	KafkaBroker string
}

func Load() Config {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	return Config{
		Port:        port,
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", user, pass, db),
		KafkaBroker: kafkaBroker,
	}
}
