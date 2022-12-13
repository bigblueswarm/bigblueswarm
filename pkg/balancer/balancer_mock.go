// Package balancer manage the balancer progress and choose the next server
package balancer

// Mock is a mock implementation of the Balancer interface
type Mock struct{}

var (
	// BalancerMockProcessFunc is the function to be called when Process is called
	BalancerMockProcessFunc func(instances []string) (string, error)
	// BalancerMockClusterStatusFunc is the function to be called when ClusterStatus is called
	BalancerMockClusterStatusFunc func(instances []string) ([]InstanceStatus, error)
	//BalancerGetCurrentStateFunc is the function to be called when GetCurrentState is called
	BalancerGetCurrentStateFunc func(measurement string, field string) (int64, error)
)

// Process is a mock implementation of the Process method
func (b *Mock) Process(instances []string) (string, error) {
	return BalancerMockProcessFunc(instances)
}

// ClusterStatus is a mock implementation of the ClusterStatus method
func (b *Mock) ClusterStatus(instances []string) ([]InstanceStatus, error) {
	return BalancerMockClusterStatusFunc(instances)
}

// GetCurrentState is a mock implementation of the GetCurrentState method
func (b *Mock) GetCurrentState(measurement string, field string) (int64, error) {
	return BalancerGetCurrentStateFunc(measurement, field)
}
