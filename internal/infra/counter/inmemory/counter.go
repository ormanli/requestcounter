package inmemory

import (
	"context"
	"sync/atomic"
)

// Counter represents a counter using atomic.Int64.
type Counter struct {
	counter atomic.Int64
}

// Increment increments the counter by the specified amount and returns the new value.
func (c *Counter) Increment(_ context.Context, delta int64) (int64, error) {
	return c.counter.Add(delta), nil
}
