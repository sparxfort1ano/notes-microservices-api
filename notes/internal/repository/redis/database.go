package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedDatabase struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCachedDatabase(cfg config) (*CachedDatabase, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &CachedDatabase{
		client: client,
		ttl:    cfg.TTL,
	}, nil
}

func (c *CachedDatabase) TTL() time.Duration {
	return c.ttl
}

func (c *CachedDatabase) Close() error {
	return c.client.Close()
}
