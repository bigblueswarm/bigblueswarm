package admin

import (
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

func computeErr(err error) error {
	if err == redis.Nil {
		return nil
	}

	return err
}

// Exists checks if an instance exists
func (m *InstanceManager) Exists(instance BigBlueButtonInstance) (bool, error) {
	return m.RDB.HExists(ctx, B3LBInstances, instance.URL).Result()
}

// List returns the list of instances
func (m *InstanceManager) List() ([]string, error) {
	instances, err := m.RDB.HKeys(ctx, B3LBInstances).Result()
	return instances, computeErr(err)
}

// Add adds an instance to the manager
func (m *InstanceManager) Add(instance BigBlueButtonInstance) error {
	_, err := m.RDB.HSet(ctx, B3LBInstances, instance.URL, instance.Secret).Result()
	return computeErr(err)
}

// Remove and instance from the manager
func (m *InstanceManager) Remove(URL string) error {
	_, err := m.RDB.HDel(ctx, B3LBInstances, URL).Result()
	return computeErr(err)
}
