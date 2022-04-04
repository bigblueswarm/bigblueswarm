package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SLedunois/b3lb/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
)

const tenantPrefix = "tenant:%s"
const hostSpec = "host"

// TenantManager is a struct manager b3lb tenants
type TenantManager interface {
	// AddTenant add a tenant in the manager
	AddTenant(tenant *Tenant) error
	// ListTenants list all tenants in the system
	ListTenants() ([]string, error)
	// DeleteTenant delete a specific tenant based on tenant hostname
	DeleteTenant(hostname string) error
	// GetTenant retrieve a tenant from a hostname
	GetTenant(hostname string) (*Tenant, error)
}

// RedisTenantManager is the redis implementation of TenantManager
type RedisTenantManager struct {
	RDB *redis.Client
}

// NewTenantManager initialize a new tenant manager
func NewTenantManager(redis redis.Client) TenantManager {
	return &RedisTenantManager{
		RDB: &redis,
	}
}

func tenantKey(key string) string {
	return fmt.Sprintf(tenantPrefix, key)
}

// AddTenant store tenant in redis
func (r *RedisTenantManager) AddTenant(tenant *Tenant) error {
	if host, ok := tenant.Spec[hostSpec]; !ok || host == "" {
		return errors.New("tenant host chould not be nil or empty string")
	}

	value, err := yaml.Marshal(tenant)
	if err != nil {
		return err
	}

	_, rErr := r.RDB.Set(context.Background(), tenantKey(tenant.Spec[hostSpec]), string(value), 0).Result()
	return utils.ComputeErr(rErr)
}

// ListTenants list all tenants in the system
func (r *RedisTenantManager) ListTenants() ([]string, error) {
	tenants, err := r.RDB.Keys(context.Background(), tenantKey("*")).Result()
	list := []string{}
	for _, tenant := range tenants {
		list = append(list, strings.ReplaceAll(tenant, "tenant:", ""))
	}

	return list, utils.ComputeErr(err)
}

// DeleteTenant delete a specific tenant based on tenant hostname
func (r *RedisTenantManager) DeleteTenant(hostname string) error {
	_, err := r.RDB.Del(context.Background(), tenantKey(hostname)).Result()
	return utils.ComputeErr(err)
}

// GetTenant retrieve a tenant from a hostname
func (r *RedisTenantManager) GetTenant(hostname string) (*Tenant, error) {
	res, err := r.RDB.Get(context.Background(), tenantKey(hostname)).Result()
	if utils.ComputeErr(err) != nil {
		return nil, err
	}

	if res == "" {
		return nil, nil
	}

	var tenant Tenant
	if err := yaml.Unmarshal([]byte(res), &tenant); err != nil {
		return nil, err
	}

	return &tenant, nil
}
