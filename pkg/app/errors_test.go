package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoInstanceFoundErr(t *testing.T) {
	err := noInstanceFoundError()
	assert.Equal(t, "FAILED", err.ReturnCode)
	assert.Equal(t, "noInstanceFound", err.MessageKey)
	assert.Equal(t, "BigBlueSwarm do not found a valid BigBlueButton instance for your request", err.Message)
}

func TestServerError(t *testing.T) {
	message := "dummy message"
	err := serverError(message)
	assert.Equal(t, "FAILED", err.ReturnCode)
	assert.Equal(t, "internalError", err.MessageKey)
	assert.Equal(t, message, err.Message)
}

func TestGetTenantError(t *testing.T) {
	assert.Equal(t, "BigBlueSwarm failed to retrieve the request tenant", getTenantError().Message)
}

func TestMeetingPoolReacherError(t *testing.T) {
	err := meetingPoolReachedError()
	assert.Equal(t, "meetingPoolReached", err.MessageKey)
	assert.Equal(t, "Your tenant reached the meeting pool limit and can't create a new one.", err.Message)
}

func TestUserPoolReachedError(t *testing.T) {
	err := userPoolReachedError()
	assert.Equal(t, "userPoolReached", err.MessageKey)
	assert.Equal(t, "Your tenant reached the user pool limit.", err.Message)
}
