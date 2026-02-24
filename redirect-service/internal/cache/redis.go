package cache

import (
	"context"
	"errors"
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

func (c *RedisCache) GetURL(ctx context.Context, shortCode string) (string, error) {
	val, err := c.client.Get(ctx, fmt.Sprintf("short:%s", shortCode)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *RedisCache) SetURL(ctx context.Context, shortCode, originalURL string) error {
	return c.client.Set(ctx, fmt.Sprintf("short:%s", shortCode), originalURL, urlTTL).Err()
}

func (c *RedisCache) DeleteURL(ctx context.Context, shortCode string) error {
	return c.client.Del(ctx, fmt.Sprintf("short:%s", shortCode)).Err()
}
