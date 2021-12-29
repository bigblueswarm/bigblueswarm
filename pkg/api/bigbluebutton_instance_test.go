package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Create with valid parameters should return a valid create response", func(t *testing.T) {
		meetingID := "id"
		params := fmt.Sprintf("name=doe&meetingID=%s&moderatorPW=pwd", meetingID)
		response := instance.Create(params)
		assert.Equal(t, response.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, response.MeetingID, meetingID)
	})
}

func TestGetJoinRedirectURL(t *testing.T) {
	t.Run("Valid join call should return a valid join redirect url", func(t *testing.T) {
		params := "meetingID=id&fullName=Simon&password=pwd"
		url, err := instance.GetJoinRedirectURL(params)

		if err != nil {
			panic(err)
		}

		expectedURL := fmt.Sprintf("%s/api/join?%s&%s", instance.URL, params, "checksum=fd830ac8255d26170825bb676a746754dab86731")
		assert.Equal(t, expectedURL, url)
	})
}
