package admin

import "github.com/SLedunois/b3lb/pkg/api"

// InstanceManagerMock is a mock implementation of the InstanceManager interface.
type InstanceManagerMock struct{}

var (
	// ExistsInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	ExistsInstanceManagerMockFunc func(instance api.BigBlueButtonInstance) (bool, error)
	// ListInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	ListInstanceManagerMockFunc func() ([]string, error)
	// AddInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	AddInstanceManagerMockFunc func(instance api.BigBlueButtonInstance) error
	// RemoveInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	RemoveInstanceManagerMockFunc func(URL string) error
	// GetInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	GetInstanceManagerMockFunc func(URL string) (api.BigBlueButtonInstance, error)
	// ListInstancesInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	ListInstancesInstanceManagerMockFunc func() ([]api.BigBlueButtonInstance, error)
	// SetInstancesInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	SetInstancesInstanceManagerMockFunc func(instances map[string]string) error
)

// Exists is a mock implementation that returns true if the instance exists.
func (m *InstanceManagerMock) Exists(instance api.BigBlueButtonInstance) (bool, error) {
	return ExistsInstanceManagerMockFunc(instance)
}

// List is a mock implementation that returns a list of all instances.
func (m *InstanceManagerMock) List() ([]string, error) {
	return ListInstanceManagerMockFunc()
}

// Add is a mock implementation that adds a new instance.
func (m *InstanceManagerMock) Add(instance api.BigBlueButtonInstance) error {
	return AddInstanceManagerMockFunc(instance)
}

// Remove is a mock implementation that removes an instance.
func (m *InstanceManagerMock) Remove(URL string) error {
	return RemoveInstanceManagerMockFunc(URL)
}

// Get is a mock implementation that returns an instance.
func (m *InstanceManagerMock) Get(URL string) (api.BigBlueButtonInstance, error) {
	return GetInstanceManagerMockFunc(URL)
}

// ListInstances is a mock implementation that returns all instances.
func (m *InstanceManagerMock) ListInstances() ([]api.BigBlueButtonInstance, error) {
	return ListInstancesInstanceManagerMockFunc()
}

// SetInstances is a mock implementation that set all instances.
func (m *InstanceManagerMock) SetInstances(instances map[string]string) error {
	return SetInstancesInstanceManagerMockFunc(instances)
}
