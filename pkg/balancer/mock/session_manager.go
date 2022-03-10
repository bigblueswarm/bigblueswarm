package mock

// SessionManager is a mock implementation of the SessionManager interface
type SessionManager struct{}

var (
	// AddFunc is the function to be called when Add is called
	AddFunc func(sessionID string, host string) error
	// GetFunc is the function to be called when Get is called
	GetFunc func(sessionID string) (string, error)
	// RemoveFunc is the function to be called when Remove is called
	RemoveFunc func(sessionID string) error
)

// Add is a mock implementation that persist the session in the redis database
func (m *SessionManager) Add(sessionID string, host string) error {
	return AddFunc(sessionID, host)
}

// Get is a mock implementation that retrieve the session from the redis database
func (m *SessionManager) Get(sessionID string) (string, error) {
	return GetFunc(sessionID)
}

// Remove is a mock implementation that remove the session from the redis database
func (m *SessionManager) Remove(sessionID string) error {
	return RemoveFunc(sessionID)
}
