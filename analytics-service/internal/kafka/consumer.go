package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"

	"github.com/gabrieldrouin/shortly/analytics-service/internal/repository"
)

type ClickEvent struct {
	ShortCode string `json:"short_code"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}

type Consumer struct {
	reader *kafka.Reader
	repo   *repository.ClickRepository
}

func NewConsumer(broker string, repo *repository.ClickRepository) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "url-clicks",
		GroupID: "analytics-service",
	})
	return &Consumer{reader: r, repo: repo}
}

func (c *Consumer) Run(ctx context.Context) {
	slog.Info("kafka consumer started")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("kafka consumer stopped")
				return
			}
			slog.Error("failed to read kafka message", "error", err)
			continue
		}

		var event ClickEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			slog.Error("failed to unmarshal click event", "error", err)
			continue
		}

		if err := c.repo.Insert(ctx, event.ShortCode, event.UserAgent, event.Referer); err != nil {
			slog.Error("failed to insert click event", "error", err, "short_code", event.ShortCode)
			continue
		}

		slog.Debug("click event persisted", "short_code", event.ShortCode)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
