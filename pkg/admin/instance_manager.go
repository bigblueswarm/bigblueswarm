package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/api"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// BBSInstances is the key for the list of instances
const BBSInstances = "instances:list"

// InstanceManager manager Bigbluebutton instances
type InstanceManager interface {
	List() ([]string, error)
	ListInstances() ([]api.BigBlueButtonInstance, error)
	Add(instance api.BigBlueButtonInstance) error
	Get(URL string) (api.BigBlueButtonInstance, error)
	SetInstances(instances map[string]string) error
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

// List returns the list of instances
func (m *RedisInstanceManager) List() ([]string, error) {
	instances, err := m.RDB.HKeys(ctx, BBSInstances).Result()
	return instances, utils.ComputeErr(err)
}

// Add adds an instance to the manager
func (m *RedisInstanceManager) Add(instance api.BigBlueButtonInstance) error {
	_, err := m.RDB.HSet(ctx, BBSInstances, instance.URL, instance.Secret).Result()
	return utils.ComputeErr(err)
}

// Get retrieve a BigBlueButton instance based on its url
func (m *RedisInstanceManager) Get(URL string) (api.BigBlueButtonInstance, error) {
	secret, err := m.RDB.HGet(ctx, BBSInstances, URL).Result()

	if secret == "" {
		return api.BigBlueButtonInstance{}, errors.New("instance not found")
	}

	return api.BigBlueButtonInstance{
		URL:    URL,
		Secret: secret,
	}, utils.ComputeErr(err)
}

// ListInstances retrieve all instance as a BigBlueButtonInstance array
func (m *RedisInstanceManager) ListInstances() ([]api.BigBlueButtonInstance, error) {
	instanceMap, err := m.RDB.HGetAll(ctx, BBSInstances).Result()

	instances := make([]api.BigBlueButtonInstance, 0)
	for k, v := range instanceMap {
		instances = append(instances, api.BigBlueButtonInstance{
			URL:    k,
			Secret: v,
		})
	}

	return instances, utils.ComputeErr(err)
}

// SetInstances set instances map
func (m *RedisInstanceManager) SetInstances(instances map[string]string) error {
	_, err := m.RDB.Del(ctx, BBSInstances).Result()
	if utils.ComputeErr(err) != nil {
		return fmt.Errorf("failed to clear instances: %s", err)
	}

	for url, secret := range instances {
		err := m.Add(api.BigBlueButtonInstance{
			URL:    url,
			Secret: secret,
		})

		if err != nil {
			return fmt.Errorf("failed to add instance %s (secret %s). Process stopped: %s", url, secret, err)
		}
	}

	return nil
}
