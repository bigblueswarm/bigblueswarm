package admin

import (
	"b3lb/pkg/api"
	"b3lb/pkg/utils"
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// B3LBInstances is the key for the list of instances
const B3LBInstances = "b3lb_instances"

// InstanceManager manager Bigbluebutton instances
type InstanceManager interface {
	Exists(instance api.BigBlueButtonInstance) (bool, error)
	List() ([]string, error)
	ListInstances() ([]api.BigBlueButtonInstance, error)
	Add(instance api.BigBlueButtonInstance) error
	Remove(URL string) error
	Get(URL string) (api.BigBlueButtonInstance, error)
}

// RedisInstanceManager is the redis implementation of InstanceManager
type RedisInstanceManager struct {
	RDB *redis.Client
}

// NewInstanceManager creates a new instance manager
func NewInstanceManager(rdb redis.Client) InstanceManager {
	return &RedisInstanceManager{
		RDB: &rdb,
	}
}

// Exists checks if an instance exists
func (m *RedisInstanceManager) Exists(instance api.BigBlueButtonInstance) (bool, error) {
	exists, err := m.RDB.HExists(ctx, B3LBInstances, instance.URL).Result()
	return exists, utils.ComputeErr(err)
}

// List returns the list of instances
func (m *RedisInstanceManager) List() ([]string, error) {
	instances, err := m.RDB.HKeys(ctx, B3LBInstances).Result()
	return instances, utils.ComputeErr(err)
}

// Add adds an instance to the manager
func (m *RedisInstanceManager) Add(instance api.BigBlueButtonInstance) error {
	_, err := m.RDB.HSet(ctx, B3LBInstances, instance.URL, instance.Secret).Result()
	return utils.ComputeErr(err)
}

// Remove and instance from the manager
func (m *RedisInstanceManager) Remove(URL string) error {
	_, err := m.RDB.HDel(ctx, B3LBInstances, URL).Result()
	return utils.ComputeErr(err)
}

// Get retrieve a BigBlueButton instance based on its url
func (m *RedisInstanceManager) Get(URL string) (api.BigBlueButtonInstance, error) {
	secret, err := m.RDB.HGet(ctx, B3LBInstances, URL).Result()

	if secret == "" {
		return api.BigBlueButtonInstance{}, errors.New("Instance not found")
	}

	return api.BigBlueButtonInstance{
		URL:    URL,
		Secret: secret,
	}, utils.ComputeErr(err)
}

// ListInstances retrieve all instance as a BigBlueButtonInstance array
func (m *RedisInstanceManager) ListInstances() ([]api.BigBlueButtonInstance, error) {
	instanceMap, err := m.RDB.HGetAll(ctx, B3LBInstances).Result()

	instances := make([]api.BigBlueButtonInstance, 0)
	for k, v := range instanceMap {
		instances = append(instances, api.BigBlueButtonInstance{
			URL:    k,
			Secret: v,
		})
	}

	return instances, utils.ComputeErr(err)
}
