package requestcounter

// IncrementResult represents the result of an increment operation.
type IncrementResult struct {
	// LocalCount is the number of increments performed on the local machine.
	LocalCount int64
	// ClusterCount is the number of increments performed across all machines in the cluster.
	ClusterCount int64
}
