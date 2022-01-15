package mock

import "github.com/SLedunois/b3lb/pkg/api"

// InstanceManager is a mock implementation of the InstanceManager interface.
type InstanceManager struct{}

var (
	// ExistsFunc is the function that will be called when the mock instance manager is used.
	ExistsFunc func(instance api.BigBlueButtonInstance) (bool, error)
	// ListFunc is the function that will be called when the mock instance manager is used.
	ListFunc func() ([]string, error)
	// AddFunc is the function that will be called when the mock instance manager is used.
	AddFunc func(instance api.BigBlueButtonInstance) error
	// RemoveFunc is the function that will be called when the mock instance manager is used.
	RemoveFunc func(URL string) error
	// GetFunc is the function that will be called when the mock instance manager is used.
	GetFunc func(URL string) (api.BigBlueButtonInstance, error)
	// ListInstancesFunc is the function that will be called when the mock instance manager is used.
	ListInstancesFunc func() ([]api.BigBlueButtonInstance, error)
)

// Exists is a mock implementation that returns true if the instance exists.
func (m *InstanceManager) Exists(instance api.BigBlueButtonInstance) (bool, error) {
	return ExistsFunc(instance)
}

// List is a mock implementation that returns a list of all instances.
func (m *InstanceManager) List() ([]string, error) {
	return ListFunc()
}

// Add is a mock implementation that adds a new instance.
func (m *InstanceManager) Add(instance api.BigBlueButtonInstance) error {
	return AddFunc(instance)
}

// Remove is a mock implementation that removes an instance.
func (m *InstanceManager) Remove(URL string) error {
	return RemoveFunc(URL)
}

// Get is a mock implementation that returns an instance.
func (m *InstanceManager) Get(URL string) (api.BigBlueButtonInstance, error) {
	return GetFunc(URL)
}

// ListInstances is a mock implementation that returns all instances.
func (m *InstanceManager) ListInstances() ([]api.BigBlueButtonInstance, error) {
	return ListInstancesFunc()
}
