package app

import (
	"b3lb/pkg/utils"
	"context"

	"github.com/go-redis/redis/v8"
)

// SessionManager internally manage remote bigbluebutton session
type SessionManager struct {
	RDB *redis.Client
}

// NewSessionManager creates a new SessionManager
func NewSessionManager(rdb redis.Client) *SessionManager {
	return &SessionManager{
		RDB: &rdb,
	}
}

// Add persist the session in the redis database
func (m *SessionManager) Add(sessionID string, host string) error {
	_, err := m.RDB.Set(context.Background(), sessionID, host, 0).Result()

	return utils.ComputeErr(err)
}

// Get retrieve the session from the redis database
func (m *SessionManager) Get(sessionID string) (string, error) {
	host, err := m.RDB.Get(context.Background(), sessionID).Result()

	return host, utils.ComputeErr(err)
}
