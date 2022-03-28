package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/SLedunois/b3lb/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
)

const hostSpec = "host"

// TenantManager is a struct manager b3lb tenants
type TenantManager interface {
	// AddTenant add a tenant in the manager
	AddTenant(tenant *Tenant) error
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
	return fmt.Sprintf("tenant:%s", key)
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
