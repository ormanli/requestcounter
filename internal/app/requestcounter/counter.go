package requestcounter

import "context"

// Counter defines the interface for a counter.
type Counter interface {
	// Increment the counter by the specified amount and returns the new value.
	Increment(ctx context.Context, delta int64) (int64, error)
}
