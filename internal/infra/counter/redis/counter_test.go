package redis

import (
	"context"
	"testing"

	"github.com/orlangure/gnomock"
	gnomockRedis "github.com/orlangure/gnomock/preset/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/ormanli/requestcounter/internal/app/requestcounter"
)

func TestIncrement(t *testing.T) {
	container, err := gnomock.Start(gnomockRedis.Preset(gnomockRedis.WithVersion("7-alpine")))
	require.NoError(t, err)

	defer func() {
		_ = gnomock.Stop(container)
	}()

	client := redis.NewClient(&redis.Options{
		Addr: container.DefaultAddress(),
	})

	counter, err := NewCounter(context.TODO(), client, requestcounter.Config{
		RedisCounterKey: "test",
	})
	require.NoError(t, err)

	err = client.Get(context.TODO(), "test").Err()
	require.ErrorIs(t, err, redis.Nil)

	count, err := counter.Increment(context.TODO(), 1)
	require.NoError(t, err)
	require.EqualValues(t, 1, count)

	result, err := client.Get(context.TODO(), "test").Result()
	require.NoError(t, err)
	require.EqualValues(t, "1", result)

	count, err = counter.Increment(context.TODO(), 5)
	require.NoError(t, err)
	require.EqualValues(t, 6, count)

	result, err = client.Get(context.TODO(), "test").Result()
	require.NoError(t, err)
	require.EqualValues(t, "6", result)
}
