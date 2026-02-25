package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type ClickEvent struct {
	ShortCode string `json:"short_code"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        "url-clicks",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
	}
	return &Producer{writer: w}
}

func (p *Producer) PublishClick(ctx context.Context, shortCode, userAgent, referer string) error {
	event := ClickEvent{
		ShortCode: shortCode,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		UserAgent: userAgent,
		Referer:   referer,
	}

	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(shortCode),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
