package redis

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff/v4"
	"github.com/redis/go-redis/v9"

	"github.com/ormanli/requestcounter/internal/app/requestcounter"
)

// Counter represents a counter that is using Redis.
type Counter struct {
	client       *redis.Client
	incrementKey string
}

// Increment increments the counter by the specified amount and returns the new value.
func (c *Counter) Increment(ctx context.Context, delta int64) (int64, error) {
	return c.client.IncrBy(ctx, c.incrementKey, delta).Result()
}

// NewCounter creates a new instance of the Counter type with the specified Redis client and configuration.
func NewCounter(ctx context.Context, client *redis.Client, cfg requestcounter.Config) (*Counter, error) {
	err := backoff.Retry(func() error {
		return client.Ping(ctx).Err()
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(cfg.RedisRetryDuration), uint64(cfg.RedisMaxRetry)))
	if err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &Counter{
		client:       client,
		incrementKey: cfg.RedisCounterKey,
	}, nil
}
