package admin

// TenantManagerMock is a mock implementation of the TenantManager interface.
type TenantManagerMock struct{}

var (
	// AddTenantTenantManagerMockFunc is the function that will be called when the mock tenant manager is used.
	AddTenantTenantManagerMockFunc func(tenant *Tenant) error
)

// AddTenant is a mock implementation that add a tenant
func (t *TenantManagerMock) AddTenant(tenant *Tenant) error {
	return AddTenantTenantManagerMockFunc(tenant)
}
