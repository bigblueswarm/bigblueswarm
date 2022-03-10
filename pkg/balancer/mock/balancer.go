package mock

import "github.com/SLedunois/b3lb/pkg/balancer"

// Balancer is a mock implementation of the Balancer interface
type Balancer struct{}

var (
	// BalancerProcessFunc is the function to be called when Process is called
	BalancerProcessFunc func(instances []string) (string, error)
	// BalancerClusterStatusFunc is the function to be called when ClusterStatus is called
	BalancerClusterStatusFunc func(instances []string) ([]balancer.InstanceStatus, error)
)

// Process is a mock implementation of the Process method
func (b *Balancer) Process(instances []string) (string, error) {
	return BalancerProcessFunc(instances)
}

// ClusterStatus is a mock implementation of the ClusterStatus method
func (b *Balancer) ClusterStatus(instances []string) ([]balancer.InstanceStatus, error) {
	return BalancerClusterStatusFunc(instances)
}
