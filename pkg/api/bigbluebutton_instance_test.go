package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v3/pkg/restclient"
	"github.com/bigblueswarm/test_utils/pkg/test"
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

const meetingID = "id"

type bbbTest struct {
	Name            string
	HasParams       bool
	Params          string
	ExpectedError   bool
	MockFunction    func(req *http.Request) (*http.Response, error)
	CustomValidator func(t *testing.T, response interface{})
}

func marshall(action string, value interface{}) ([]byte, error) {
	switch action {
	case "GetRecordingTextTracks":
		return json.Marshal(value)
	default:
		return xml.Marshal(value)
	}
}

func getTests(action string, hasParams bool, params string, validResponse interface{}, customValidator func(*testing.T, interface{})) []bbbTest {
	return []bbbTest{
		{
			Name:            fmt.Sprintf("%s should return a valid response", action),
			Params:          params,
			HasParams:       hasParams,
			ExpectedError:   false,
			CustomValidator: customValidator,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				response, err := marshall(action, validResponse)

				if err != nil {
					panic(err)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(response)),
				}, nil
			},
		},
		{
			Name:          fmt.Sprintf("%s should return an error if remote instance returns a 500 status code", action),
			Params:        "",
			HasParams:     hasParams,
			ExpectedError: true,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, nil
			},
		},
		{
			Name:          fmt.Sprintf("%s should return an error if remote instance call returns and error", action),
			Params:        "",
			HasParams:     hasParams,
			ExpectedError: true,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("Unexpected error")
			},
		},
	}
}

func executeTests(t *testing.T, action string, tests []bbbTest) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			restclient.RestClientMockDoFunc = test.MockFunction
			method := reflect.ValueOf(instance).MethodByName(action)
			if method.IsNil() {
				panic(fmt.Sprintf("Method %s not found", action))
			}

			var values []reflect.Value
			if test.HasParams {
				values = method.Call([]reflect.Value{reflect.ValueOf(test.Params)})
			} else {
				values = method.Call([]reflect.Value{})
			}
			response := values[0].Interface()
			err := values[1].Interface()

			assert.True(t, test.ExpectedError == (err != nil))

			if test.CustomValidator != nil {
				test.CustomValidator(t, response)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	validResponse := &CreateResponse{
		Response: Response{
			ReturnCode: ReturnCodes().Success,
		},
		MeetingID: meetingID,
	}

	customValidator := func(t *testing.T, response interface{}) {
		creation, ok := response.(*CreateResponse)
		if !ok {
			t.Error("Response is not a CreateResponse")
			return
		}

		assert.Equal(t, creation.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, creation.MeetingID, meetingID)
	}
	tests := getTests("Create", true, fmt.Sprintf("name=doe&meetingID=%s&moderatorPW=pwd", meetingID), validResponse, customValidator)

	executeTests(t, "Create", tests)
}

func TestGetJoinRedirectURL(t *testing.T) {
	t.Run("Valid join call should return a valid join redirect url", func(t *testing.T) {
		params := fmt.Sprintf("meetingID=%s&fullName=Simon&password=pwd", meetingID)
		url, err := instance.GetJoinRedirectURL(params)

		if err != nil {
			panic(err)
		}

		expectedURL := fmt.Sprintf("%s/api/join?%s&%s", instance.URL, params, "checksum=ca7b6a04636c6fba1dd6158a1f2b72ab4811472a")
		assert.Equal(t, expectedURL, url)
	})
}

func TestEnd(t *testing.T) {
	validResponse := &EndResponse{
		Response: Response{
			ReturnCode: ReturnCodes().Success,
		},
	}

	customValidator := func(t *testing.T, response interface{}) {
		end, ok := response.(*EndResponse)
		if !ok {
			t.Error("Response is not a EndResponse")
			return
		}

		assert.Equal(t, end.ReturnCode, ReturnCodes().Success)
	}
	tests := getTests("End", true, fmt.Sprintf("meetingID=%s&password=pwd", meetingID), validResponse, customValidator)

	executeTests(t, "End", tests)
}

func TestIsMeetingRunning(t *testing.T) {
	validResponse := &IsMeetingsRunningResponse{
		ReturnCode: ReturnCodes().Success,
		Running:    true,
	}

	customValidator := func(t *testing.T, response interface{}) {
		running, ok := response.(*IsMeetingsRunningResponse)
		if !ok {
			t.Error("Response is not a IsMeetingsRunningResponse")
			return
		}

		assert.Equal(t, running.ReturnCode, ReturnCodes().Success)
		assert.True(t, running.Running)
	}

	tests := getTests("IsMeetingRunning", true, fmt.Sprintf("meetingID=%s", meetingID), validResponse, customValidator)

	executeTests(t, "IsMeetingRunning", tests)
}

func TestGetMeetingInfo(t *testing.T) {
	validResponse := &GetMeetingInfoResponse{
		ReturnCode: ReturnCodes().Success,
		MeetingInfo: MeetingInfo{
			MeetingID: meetingID,
		},
	}

	customValidator := func(t *testing.T, response interface{}) {
		info, ok := response.(*GetMeetingInfoResponse)
		if !ok {
			t.Error("Response is not a GetMeetingInfoResponse")
			return
		}

		assert.Equal(t, info.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, info.MeetingInfo.MeetingID, meetingID)
	}

	tests := getTests("GetMeetingInfo", true, fmt.Sprintf("meetingID=%s", meetingID), validResponse, customValidator)

	executeTests(t, "GetMeetingInfo", tests)
}

func TestGetMeetings(t *testing.T) {
	validResponse := &GetMeetingsResponse{
		ReturnCode: ReturnCodes().Success,
		Meetings: []MeetingInfo{
			{
				MeetingID: meetingID,
			},
		},
	}

	customValidator := func(t *testing.T, response interface{}) {
		meetings, ok := response.(*GetMeetingsResponse)
		if !ok {
			t.Error("Response is not a GetMeetingInfoResponse")
			return
		}

		assert.Equal(t, meetings.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, 1, len(meetings.Meetings))
		assert.Equal(t, meetingID, meetings.Meetings[0].MeetingID)
	}

	tests := getTests("GetMeetings", false, "", validResponse, customValidator)

	executeTests(t, "GetMeetings", tests)
}

func TestJoin(t *testing.T) {
	joinURL := "http://localhost/client/BigBlueButton.html?sessionToken=ai1wqj8wb6s7rnk0"
	validResponse := &JoinRedirectResponse{
		Response: Response{
			ReturnCode: ReturnCodes().Success,
		},
		URL: joinURL,
	}

	customValidator := func(t *testing.T, response interface{}) {
		join, ok := response.(*JoinRedirectResponse)
		if !ok {
			t.Error("Response is not a JoinRedirectResponse")
			return
		}

		assert.Equal(t, ReturnCodes().Success, join.ReturnCode)
		assert.Equal(t, joinURL, join.URL)
	}

	tests := getTests("Join", true, fmt.Sprintf("meetingID=%s&fullName=Simon&password=pwd", meetingID), validResponse, customValidator)

	executeTests(t, "Join", tests)
}

func TestGetRecordings(t *testing.T) {
	validResponse := &GetRecordingsResponse{
		Response: Response{
			ReturnCode: ReturnCodes().Success,
		},
		Recordings: []Recording{
			{
				RecordID:  "recordID",
				MeetingID: meetingID,
			},
		},
	}

	customValidator := func(t *testing.T, response interface{}) {
		recordings, ok := response.(*GetRecordingsResponse)
		if !ok {
			t.Error("Response is not a GetRecordingsResponse")
			return
		}

		assert.Equal(t, recordings.Response.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, 1, len(recordings.Recordings))
		assert.Equal(t, meetingID, recordings.Recordings[0].MeetingID)
		assert.Equal(t, "recordID", recordings.Recordings[0].RecordID)
	}

	tests := getTests("GetRecordings", true, "", validResponse, customValidator)

	executeTests(t, "GetRecordings", tests)
}

func TestUpdateRecordings(t *testing.T) {
	validResponse := &UpdateRecordingsResponse{
		ReturnCode: ReturnCodes().Success,
		Updated:    true,
	}

	customValidator := func(t *testing.T, response interface{}) {
		recordings, ok := response.(*UpdateRecordingsResponse)
		if !ok {
			t.Error("Response is not a UpdateRecordingsResponse")
			return
		}

		assert.Equal(t, recordings.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, recordings.Updated, true)
	}

	tests := getTests("UpdateRecordings", true, "recordID=recording-id", validResponse, customValidator)

	executeTests(t, "UpdateRecordings", tests)
}

func TestDeleteRecordings(t *testing.T) {
	validResponse := &DeleteRecordingsResponse{
		ReturnCode: ReturnCodes().Success,
		Deleted:    true,
	}

	customValidator := func(t *testing.T, response interface{}) {
		recordings, ok := response.(*DeleteRecordingsResponse)
		if !ok {
			t.Error("Response is not a DeleteRecordingsRespomse")
			return
		}

		assert.Equal(t, recordings.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, recordings.Deleted, true)
	}

	tests := getTests("DeleteRecordings", true, "recordID=recording-id", validResponse, customValidator)

	executeTests(t, "DeleteRecordings", tests)
}

func TestPublishRecordings(t *testing.T) {
	validResponse := &PublishRecordingsResponse{
		ReturnCode: ReturnCodes().Success,
		Published:  true,
	}

	customValidator := func(t *testing.T, response interface{}) {
		recordings, ok := response.(*PublishRecordingsResponse)
		if !ok {
			t.Error("Response is not a PublishRecordingsResponse")
			return
		}

		assert.Equal(t, recordings.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, recordings.Published, true)
	}

	tests := getTests("PublishRecordings", true, "recordID=recording-id&published=true", validResponse, customValidator)

	executeTests(t, "PublishRecordings", tests)
}

func TestGetRecordingsTextTracks(t *testing.T) {
	validResponse := &GetRecordingsTextTracksResponse{
		Response: RecordingsTextTrackResponseType{
			ReturnCode: ReturnCodes().Success,
			Tracks: []Track{
				{
					Href:   "http://localhost/client/api/v1/recordings/recording-id/texttracks/track-id",
					Kind:   "subtitles",
					Lang:   "en",
					Label:  "English",
					Source: "upload",
				},
			},
		},
	}

	customValidator := func(t *testing.T, response interface{}) {
		tracks, ok := response.(*GetRecordingsTextTracksResponse)
		if !ok {
			t.Error("Response is not a GetRecordingsTextTracksResponse")
			return
		}

		assert.Equal(t, tracks.Response.ReturnCode, ReturnCodes().Success)
		assert.Equal(t, len(tracks.Response.Tracks), 1)
		assert.Equal(t, tracks.Response.Tracks[0].Href, "http://localhost/client/api/v1/recordings/recording-id/texttracks/track-id")
		assert.Equal(t, tracks.Response.Tracks[0].Kind, "subtitles")
		assert.Equal(t, tracks.Response.Tracks[0].Lang, "en")
		assert.Equal(t, tracks.Response.Tracks[0].Label, "English")
		assert.Equal(t, tracks.Response.Tracks[0].Source, "upload")
	}

	tests := getTests("GetRecordingTextTracks", true, "recordID=recording-id", validResponse, customValidator)

	executeTests(t, "GetRecordingTextTracks", tests)
}

func TestRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	baseURL := "http://localhost/bigbluebutton"
	parameters := "meetingID=meeting-id"

	tests := []test.Test{
		{
			Name: "Redirect should end with a 302 status code",
			Mock: func() {},
			Validator: func(t *testing.T, response interface{}, err error) {
				assert.Equal(t, http.StatusFound, w.Code)
				assert.Equal(t, baseURL+"/api/"+GetMeetingInfo+"?"+parameters+"&checksum=606d824c6e7faeb58108561bbb1df8a3153be6e4", w.Header().Get("Location"))
			},
		},
	}

	instance := &BigBlueButtonInstance{
		URL:    baseURL,
		Secret: test.DefaultSecret(),
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			instance.Redirect(c, GetMeetingInfo, parameters)
			test.Validator(t, nil, nil)
		})
	}
}
