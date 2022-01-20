package app

import (
	"context"

	"github.com/SLedunois/b3lb/pkg/utils"

	"github.com/go-redis/redis/v8"
)

// SessionManager manages BigBlueButton sessions
type SessionManager interface {
	Add(sessionID string, host string) error
	Get(sessionID string) (string, error)
	Remove(sessionID string) error
}

// RedisSessionManager internally manage remote bigbluebutton session
type RedisSessionManager struct {
	RDB *redis.Client
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(rdb redis.Client) SessionManager {
	return &RedisSessionManager{
		RDB: &rdb,
	}
}

// cacheKey format meetingID as a valid meeting cache key
func cacheKey(id string) string {
	return "meeting:" + id
}

// Add persist the session in the redis database
func (m *RedisSessionManager) Add(sessionID string, host string) error {
	_, err := m.RDB.Set(context.Background(), cacheKey(sessionID), host, 0).Result()

	return utils.ComputeErr(err)
}

// Get retrieve the session from the redis database
func (m *RedisSessionManager) Get(sessionID string) (string, error) {
	host, err := m.RDB.Get(context.Background(), cacheKey(sessionID)).Result()

	return host, utils.ComputeErr(err)
}

// Remove remove the session from the redis database
func (m *RedisSessionManager) Remove(sessionID string) error {
	_, err := m.RDB.Del(context.Background(), cacheKey(sessionID)).Result()

	return utils.ComputeErr(err)
}
