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
			expectedCode:      api.ReturnCodes().Failed,
			expectedKey:       api.MessageKeys().ValidationError,
			expectedMessage:   api.Messages().EmptyMeetingName,
			expectedSessionID: "",
		},
		{
			name:              "Create with no meeting id should returns a `no meeting id error`",
			url:               "/bigbluebutton/api/create?name=doe&checksum=291411768c7aeb243819983459755b32b96aae34",
			expectedCode:      api.ReturnCodes().Failed,
			expectedKey:       api.MessageKeys().ValidationError,
			expectedMessage:   api.Messages().EmptyMeetingID,
			expectedSessionID: "",
		},
		{
			name:              "Valid create call should create a meeting",
			url:               "/bigbluebutton/api/create?name=doe&meetingID=id&moderatorPW=pwd&checksum=f4db98b7cab8ebc1df423e547ed3fa995d13ad72",
			expectedCode:      api.ReturnCodes().Success,
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
		assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
		assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
		assert.Equal(t, api.Messages().NotFound, response.Message)
	})

	t.Run("Joining a session should redirect", func(t *testing.T) {
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/join?meetingID=id&fullName=Simon&password=pwd&checksum=561af1e2093bba548d1c66b66a4484a93f3d5c80", nil)

		assert.Equal(t, http.StatusFound, w.Code)
	})
}

func TestEnd(t *testing.T) {
	type test struct {
		name            string
		url             string
		expectedKey     string
		expectedMessage string
		expectedCode    string
	}

	router := launchRouter(defaultConfig())

	tests := []test{
		{
			name:            "End with an invalid meeting id should returns a `not found error`",
			url:             "/bigbluebutton/api/end?meetingID=123&password=pwd&checksum=65f76052d9ad86af739dd3fa4abc10540b8ab487",
			expectedKey:     api.MessageKeys().NotFound,
			expectedMessage: api.Messages().NotFound,
			expectedCode:    api.ReturnCodes().Failed,
		},
		{
			name:            "End with an invalid password should returns a `invalid password error`",
			url:             "/bigbluebutton/api/end?meetingID=id&password=pwd2&checksum=19e0cafcffe893b246913327c847625c013efe0c",
			expectedKey:     api.MessageKeys().ValidationError,
			expectedMessage: api.Messages().InvalidModeratorPW,
			expectedCode:    api.ReturnCodes().Failed,
		},
		{
			name:            "End with a valid meeting id and a valid moderator password should returns a `success`",
			url:             "/bigbluebutton/api/end?meetingID=meeting_id&password=pwd&checksum=c748c8cc6cc38a47380c64699920b0e7d6affa80",
			expectedKey:     api.MessageKeys().SendEndMeetingRequest,
			expectedCode:    api.ReturnCodes().Success,
			expectedMessage: api.Messages().EndMeeting,
		},
	}

	TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/create?name=doe&meetingID=meeting_id&moderatorPW=pwd&checksum=b4a2dcd6ecab0c4697d46f3713f18ebc65c7e827", nil)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := TestUtil.ExecuteRequest(router, "GET", test.url, nil)
			var response api.EndResponse
			if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
				panic(err)
			}

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, test.expectedCode, response.ReturnCode)
			assert.Equal(t, test.expectedKey, response.MessageKey)
			assert.Equal(t, test.expectedMessage, response.Message)
		})
	}
}

func TestIsMeetingRunning(t *testing.T) {
	router := launchRouter(defaultConfig())

	t.Run("Checking if an non existant meeting is running should returns a not found error", func(t *testing.T) {
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/isMeetingRunning?meetingID=meeting_id&checksum=d791c7a640acb66cd30ae06d9cde28b306497308", nil)
		var response api.Response
		if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
		assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
		assert.Equal(t, api.Messages().NotFound, response.Message)
	})

	t.Run("Checking if a meeting is running should returns success response", func(t *testing.T) {
		TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/create?name=doe&meetingID=meeting_id&moderatorPW=pwd&checksum=b4a2dcd6ecab0c4697d46f3713f18ebc65c7e827", nil)
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/isMeetingRunning?meetingID=meeting_id&checksum=d791c7a640acb66cd30ae06d9cde28b306497308", nil)
		var response api.IsMeetingsRunningResponse
		if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
		assert.False(t, response.Running)
	})
}

func TestGetMeetingInfo(t *testing.T) {
	router := launchRouter(defaultConfig())
	t.Run("Checking a meeting that does not exists should returns a not found error", func(t *testing.T) {
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/getMeetingInfo?meetingID=meeting_id&checksum=b49265e3baa6b4ecb5e6caf443e3511f62c434e5", nil)
		var response api.Response
		if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
		assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
		assert.Equal(t, api.Messages().NotFound, response.Message)
	})

	t.Run("Checking info from existing meeting should returns a success response and a valid meeting identifier", func(t *testing.T) {
		TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/create?name=doe&meetingID=meeting_id&moderatorPW=pwd&checksum=b4a2dcd6ecab0c4697d46f3713f18ebc65c7e827", nil)
		w := TestUtil.ExecuteRequest(router, "GET", "/bigbluebutton/api/getMeetingInfo?meetingID=meeting_id&checksum=b49265e3baa6b4ecb5e6caf443e3511f62c434e5", nil)
		var response api.GetMeetingInfoResponse
		if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
			panic(err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
		assert.Equal(t, "meeting_id", response.MeetingID)
	})
}
