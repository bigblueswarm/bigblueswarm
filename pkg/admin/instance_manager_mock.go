package admin

import "github.com/bigblueswarm/bigblueswarm/v2/pkg/api"

// InstanceManagerMock is a mock implementation of the InstanceManager interface.
type InstanceManagerMock struct{}

var (
	// ListInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	ListInstanceManagerMockFunc func() ([]string, error)
	// AddInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	AddInstanceManagerMockFunc func(instance api.BigBlueButtonInstance) error
	// GetInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	GetInstanceManagerMockFunc func(URL string) (api.BigBlueButtonInstance, error)
	// ListInstancesInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	ListInstancesInstanceManagerMockFunc func() ([]api.BigBlueButtonInstance, error)
	// SetInstancesInstanceManagerMockFunc is the function that will be called when the mock instance manager is used.
	SetInstancesInstanceManagerMockFunc func(instances map[string]string) error
)

// List is a mock implementation that returns a list of all instances.
func (m *InstanceManagerMock) List() ([]string, error) {
	return ListInstanceManagerMockFunc()
}

// Add is a mock implementation that adds a new instance.
func (m *InstanceManagerMock) Add(instance api.BigBlueButtonInstance) error {
	return AddInstanceManagerMockFunc(instance)
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
