// Package admin manages the bigblueswarm admin part
package admin

// TenantManagerMock is a mock implementation of the TenantManager interface.
type TenantManagerMock struct{}

var (
	// AddTenantTenantManagerMockFunc is the function that will be called when the mock tenant manager is used.
	AddTenantTenantManagerMockFunc func(tenant *Tenant) error
	// ListTenantsTenantManagerMockFunc is the function that will be called when the mock tenant manager is used
	ListTenantsTenantManagerMockFunc func() ([]TenantListObject, error)
	// DeleteTenantTenantManagerMockFunc is the function that will be called when the mock tenant manager is used
	DeleteTenantTenantManagerMockFunc func(hostname string) error
	// GetTenantTenantManagerMockFunc is the function that will be called when the mock tenant manager is used
	GetTenantTenantManagerMockFunc func(hostname string) (*Tenant, error)
)

// AddTenant is a mock implementation that add a tenant
func (t *TenantManagerMock) AddTenant(tenant *Tenant) error {
	return AddTenantTenantManagerMockFunc(tenant)
}

// ListTenants is a mock implementation that list all tenants
func (t *TenantManagerMock) ListTenants() ([]TenantListObject, error) {
	return ListTenantsTenantManagerMockFunc()
}

// DeleteTenant is a mock implementation that will delete a given tenant
func (t *TenantManagerMock) DeleteTenant(hostname string) error {
	return DeleteTenantTenantManagerMockFunc(hostname)
}

// GetTenant is a mock implementation that retrieve a tenant from a hostname
func (t *TenantManagerMock) GetTenant(hostname string) (*Tenant, error) {
	return GetTenantTenantManagerMockFunc(hostname)
}
