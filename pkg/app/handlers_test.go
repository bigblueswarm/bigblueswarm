package app

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoute(t *testing.T) {
	router := launchRouter(&config.Config{})
	w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api", nil)

	response := "<response><returncode>SUCCESS</returncode><version>2.0</version></response>"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, response, w.Body.String())
}

func TestCreate(t *testing.T) {
	type test struct {
		name              string
		url               string
		expectedKey       string
		expectedMessage   string
		expectedCode      string
		expectedSessionID string
	}

	router := launchRouter(defaultConfig())

	tests := []test{
		{
			name:              "Create with no name should returns a `no name error`",
			url:               "/bigbluebutton/api/create?checksum=025401fa251bdbcbba29c347a9cf811f29aa15a1",
			expectedCode:      api.FailedReturnCode,
			expectedKey:       api.ValidationErrorMessageKey,
			expectedMessage:   api.EmptyMeetingNameMessage,
			expectedSessionID: "",
		},
		{
			name:              "Create with no meeting id should returns a `no meeting id error`",
			url:               "/bigbluebutton/api/create?name=doe&checksum=291411768c7aeb243819983459755b32b96aae34",
			expectedCode:      api.FailedReturnCode,
			expectedKey:       api.ValidationErrorMessageKey,
			expectedMessage:   api.EmptyMeetingIDMessage,
			expectedSessionID: "",
		},
		{
			name:              "Valid create call should create a meeting",
			url:               "/bigbluebutton/api/create?name=doe&meetingID=id&moderatorPW=pwd&checksum=f4db98b7cab8ebc1df423e547ed3fa995d13ad72",
			expectedCode:      api.SuccessReturnCode,
			expectedKey:       "",
			expectedMessage:   "",
			expectedSessionID: "id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := TestUtil.ExecuteRequest(router, "GET", test.url, nil)
			var response api.CreateResponse
			if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
				panic(err)
			}

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, test.expectedCode, response.ReturnCode)
			assert.Equal(t, test.expectedKey, response.MessageKey)
			assert.Equal(t, test.expectedMessage, response.Message)
			assert.Equal(t, test.expectedSessionID, response.MeetingID)
		})
	}
}

func TestJoin(t *testing.T) {
	router := launchRouter(defaultConfig())

	t.Run("Joining a non existing session should returns a notFound error", func(t *testing.T) {
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/join?meetingID=123&fullName=Simon&password=pwd&checksum=9215d3a80656cf2aa7e2a772e27fbfe7e3e4ddc8", nil)
		var response api.Error
		if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.FailedReturnCode, response.ReturnCode)
		assert.Equal(t, api.NotFoundMessageKey, response.MessageKey)
		assert.Equal(t, api.NotFoundMeetingIDMessage, response.Message)
	})

	t.Run("Joining a session should redirect", func(t *testing.T) {
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/join?meetingID=id&fullName=Simon&password=pwd&checksum=561af1e2093bba548d1c66b66a4484a93f3d5c80", nil)

		assert.Equal(t, http.StatusFound, w.Code)
	})
}
