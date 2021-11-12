package admin

import (
	"b3lb/pkg/api"
	"b3lb/pkg/utils"
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// B3LBInstances is the key for the list of instances
const B3LBInstances = "b3lb_instances"

// InstanceManager manager Bigbluebutton instances
type InstanceManager struct {
	RDB *redis.Client
}

// NewInstanceManager creates a nes instance manager
func NewInstanceManager(rdb *redis.Client) *InstanceManager {
	return &InstanceManager{
		RDB: rdb,
	}
}

// Exists checks if an instance exists
func (m *InstanceManager) Exists(instance api.BigBlueButtonInstance) (bool, error) {
	return m.RDB.HExists(ctx, B3LBInstances, instance.URL).Result()
}

// List returns the list of instances
func (m *InstanceManager) List() ([]string, error) {
	instances, err := m.RDB.HKeys(ctx, B3LBInstances).Result()
	return instances, utils.ComputeErr(err)
}

// Add adds an instance to the manager
func (m *InstanceManager) Add(instance api.BigBlueButtonInstance) error {
	_, err := m.RDB.HSet(ctx, B3LBInstances, instance.URL, instance.Secret).Result()
	return utils.ComputeErr(err)
}

// Remove and instance from the manager
func (m *InstanceManager) Remove(URL string) error {
	_, err := m.RDB.HDel(ctx, B3LBInstances, URL).Result()
	return utils.ComputeErr(err)
}

// Get retrieve a BigBlueButton instance based on its url
func (m *InstanceManager) Get(URL string) (api.BigBlueButtonInstance, error) {
	secret, err := m.RDB.HGet(ctx, B3LBInstances, URL).Result()
	return api.BigBlueButtonInstance{
		URL:    URL,
		Secret: secret,
	}, utils.ComputeErr(err)
}
