package requestcounter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIncrement(t *testing.T) {
	tests := []struct {
		name                          string
		prepareLocalMockCounterFunc   func(counter *MockCounter)
		prepareClusterMockCounterFunc func(counter *MockCounter)
		assertFunc                    func(t *testing.T, actual IncrementResult, err error)
	}{
		{
			name: "no error",
			prepareLocalMockCounterFunc: func(counter *MockCounter) {
				counter.EXPECT().Increment(mock.Anything, int64(1)).Return(2, nil)
			},
			prepareClusterMockCounterFunc: func(counter *MockCounter) {
				counter.EXPECT().Increment(mock.Anything, int64(1)).Return(3, nil)
			},
			assertFunc: func(t *testing.T, actual IncrementResult, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, IncrementResult{
					LocalCount:   2,
					ClusterCount: 3,
				}, actual)
			},
		},
		{
			name: "cluster counter error",
			prepareLocalMockCounterFunc: func(counter *MockCounter) {
			},
			prepareClusterMockCounterFunc: func(counter *MockCounter) {
				counter.EXPECT().Increment(mock.Anything, int64(1)).Return(0, errors.New("something happened for cluster counter"))
			},
			assertFunc: func(t *testing.T, actual IncrementResult, err error) {
				assert.EqualError(t, err, "something happened for cluster counter")
				assert.Empty(t, actual)
			},
		},
		{
			name: "local counter error",
			prepareLocalMockCounterFunc: func(counter *MockCounter) {
				counter.EXPECT().Increment(mock.Anything, int64(1)).Return(0, errors.New("something happened for local counter"))
			},
			prepareClusterMockCounterFunc: func(counter *MockCounter) {
				counter.EXPECT().Increment(mock.Anything, int64(1)).Return(3, nil)
			},
			assertFunc: func(t *testing.T, actual IncrementResult, err error) {
				assert.EqualError(t, err, "something happened for local counter")
				assert.EqualValues(t, IncrementResult{
					ClusterCount: 3,
				}, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			localMockCounter := NewMockCounter(t)
			clusterMockCounter := NewMockCounter(t)

			test.prepareLocalMockCounterFunc(localMockCounter)
			test.prepareClusterMockCounterFunc(clusterMockCounter)

			service := NewService(localMockCounter, clusterMockCounter)

			result, err := service.Increment(context.TODO())
			test.assertFunc(t, result, err)
		})
	}
}
