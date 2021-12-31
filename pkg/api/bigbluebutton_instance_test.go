package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const meetingID = "id"

func TestCreate(t *testing.T) {
	t.Run("Create with valid parameters should return a valid create response", func(t *testing.T) {
		params := fmt.Sprintf("name=doe&meetingID=%s&moderatorPW=pwd", meetingID)
		response := instance.Create(params)
		assert.Equal(t, response.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, response.MeetingID, meetingID)
	})
}

func TestGetJoinRedirectURL(t *testing.T) {
	t.Run("Valid join call should return a valid join redirect url", func(t *testing.T) {
		params := fmt.Sprintf("meetingID=%s&fullName=Simon&password=pwd", meetingID)
		url, err := instance.GetJoinRedirectURL(params)

		if err != nil {
			panic(err)
		}

		expectedURL := fmt.Sprintf("%s/api/join?%s&%s", instance.URL, params, "checksum=fd830ac8255d26170825bb676a746754dab86731")
		assert.Equal(t, expectedURL, url)
	})
}

func TestEnd(t *testing.T) {
	t.Run("Ending a meeting should return a valid end response", func(t *testing.T) {
		response := instance.End(fmt.Sprintf("meetingID=%s&password=pwd", meetingID))
		assert.Equal(t, response.ReturnCode, ReturnCodes().Success)
	})
}

func TestIsMeetingRunning(t *testing.T) {
	t.Run("IsMeetingRunning should return a failed response", func(t *testing.T) {
		response := instance.End(fmt.Sprintf("meetingID=%s", meetingID))
		assert.Equal(t, response.ReturnCode, ReturnCodes().Failed)
	})
}

func TestGetMeetingInfo(t *testing.T) {
	t.Run("GetMeetingInfo should return a failed response", func(t *testing.T) {
		response := instance.GetMeetingInfo(fmt.Sprintf("meetingID=%s", meetingID))
		assert.Equal(t, response.ReturnCode, ReturnCodes().Failed)
	})
}
