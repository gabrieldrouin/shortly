package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const urlTTL = 24 * time.Hour

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) SetURL(ctx context.Context, shortCode, originalURL string) error {
	return c.client.Set(ctx, fmt.Sprintf("short:%s", shortCode), originalURL, urlTTL).Err()
}
