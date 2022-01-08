package mock

// Balancer is a mock implementation of the Balancer interface
type Balancer struct{}

var (
	// BalancerProcessFunc is the function to be called when Process is called
	BalancerProcessFunc func(instances []string) (string, error)
)

// Process is a mock implementation of the Process method
func (b *Balancer) Process(instances []string) (string, error) {
	return BalancerProcessFunc(instances)
}
