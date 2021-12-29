package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnsCode(t *testing.T) {
	assert.Equal(t, "SUCCESS", ReturnCodes().Success)
	assert.Equal(t, "FAILED", ReturnCodes().Failed)
}

func TestMessageKeys(t *testing.T) {
	assert.Equal(t, "validationError", MessageKeys().ValidationError)
	assert.Equal(t, "duplicationWarning", MessageKeys().DuplicationWarning)
	assert.Equal(t, "notFound", MessageKeys().NotFound)
	assert.Equal(t, "sentEndMeetingRequest", MessageKeys().SendEndMeetingRequest)
}

func TestMessages(t *testing.T) {
	assert.Equal(t, "You must provide a meeting ID", Messages().EmptyMeetingID)
	assert.Equal(t, "You must provide a meeting name", Messages().EmptyMeetingName)
	assert.Equal(t, "This conference was already in existence and may currently be in progress.", Messages().DuplicationWarning)
	assert.Equal(t, "A meeting with that ID does not exist", Messages().NotFound)
	assert.Equal(t, "A request to end the meeting was sent. Please wait a few seconds, and then use the getMeetingInfo or isMeetingRunning API calls to verify that it was ended.", Messages().EndMeeting)
	assert.Equal(t, "Provided moderator password is incorrect", Messages().InvalidModeratorPW)
}
