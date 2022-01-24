package api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/SLedunois/b3lb/pkg/restclient/mock"

	"github.com/stretchr/testify/assert"
)

const meetingID = "id"

type test struct {
	Name            string
	HasParams       bool
	Params          string
	ExpectedError   bool
	MockFunction    func(req *http.Request) (*http.Response, error)
	CustomValidator func(t *testing.T, response interface{})
}

func getTests(action string, hasParams bool, params string, validResponse interface{}, customValidator func(*testing.T, interface{})) []test {
	return []test{
		{
			Name:            fmt.Sprintf("%s should return a valid response", action),
			Params:          params,
			HasParams:       hasParams,
			ExpectedError:   false,
			CustomValidator: customValidator,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				response, err := xml.Marshal(validResponse)

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
		{
			Name:          fmt.Sprintf("%s should return an error if remote instance response is not a valid XML", action),
			Params:        params,
			HasParams:     hasParams,
			ExpectedError: true,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"id": "123"}`))),
				}, nil
			},
		},
	}
}

func executeTests(t *testing.T, action string, tests []test) {
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mock.DoFunc = test.MockFunction
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
