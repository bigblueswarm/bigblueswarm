package app

import (
	"context"

	"github.com/SLedunois/b3lb/pkg/utils"

	"github.com/go-redis/redis/v8"
)

// Mapper manages BigBlueButton sessions
type Mapper interface {
	Add(key string, host string) error
	Get(key string) (string, error)
	Remove(key string) error
}

// RedisMapper internally manage remote bigbluebutton session
type RedisMapper struct {
	RDB *redis.Client
}

// NewMapper creates a new SessionManager
func NewMapper(rdb redis.Client) Mapper {
	return &RedisMapper{
		RDB: &rdb,
	}
}

// MeetingMapKey format meetingID as a valid meeting map key
func MeetingMapKey(id string) string {
	return "meeting:" + id
}

//RecordingMapKey format recordingID as a valid recording map key
func RecordingMapKey(id string) string {
	return "recording:" + id
}

// Add persist the session in the redis database
func (m *RedisMapper) Add(key string, host string) error {
	_, err := m.RDB.Set(context.Background(), key, host, 0).Result()

	return utils.ComputeErr(err)
}

// Get retrieve the session from the redis database
func (m *RedisMapper) Get(key string) (string, error) {
	host, err := m.RDB.Get(context.Background(), key).Result()

	return host, utils.ComputeErr(err)
}

// Remove remove the session from the redis database
func (m *RedisMapper) Remove(key string) error {
	_, err := m.RDB.Del(context.Background(), key).Result()

	return utils.ComputeErr(err)
}
