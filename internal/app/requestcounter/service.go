package requestcounter

import "context"

type Service struct {
	localCounter   Counter
	clusterCounter Counter
}

// NewService creates a new instance of the service with the specified counters.
func NewService(localCounter, clusterCounter Counter) *Service {
	return &Service{
		localCounter:   localCounter,
		clusterCounter: clusterCounter,
	}
}

// Increment increments the counter by 1 and returns the result.
func (s *Service) Increment(ctx context.Context) (IncrementResult, error) {
	clusterCount, err := s.clusterCounter.Increment(ctx, 1)
	if err != nil {
		return IncrementResult{}, err
	}

	localCount, err := s.localCounter.Increment(ctx, 1)
	if err != nil {
		return IncrementResult{
			ClusterCount: clusterCount,
		}, err
	}

	return IncrementResult{
		LocalCount:   localCount,
		ClusterCount: clusterCount,
	}, nil
}
