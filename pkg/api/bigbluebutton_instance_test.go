package api

import (
	"b3lb/pkg/restclient/mock"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const meetingID = "id"

type test struct {
	Name            string
	Params          string
	ExpectedError   bool
	MockFunction    func(req *http.Request) (*http.Response, error)
	CustomValidator func(response interface{}) bool
}

func getTests(action string, params string, validResponse interface{}, customValidator func(interface{}) bool) []test {
	return []test{
		{
			Name:            fmt.Sprintf("%s should return a valid response", action),
			Params:          params,
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
			ExpectedError: true,
			MockFunction: func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("Unexpected error")
			},
		},
		{
			Name:          fmt.Sprintf("%s should return an error if remote instance response is not a valid XML", action),
			Params:        params,
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

			values := method.Call([]reflect.Value{reflect.ValueOf(test.Params)})
			response := values[0].Interface()
			err := values[1].Interface()

			assert.True(t, test.ExpectedError == (err != nil))

			if test.CustomValidator != nil {
				assert.True(t, test.CustomValidator(response))
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

	customValidator := func(response interface{}) bool {
		creation, ok := response.(*CreateResponse)
		if !ok {
			return false
		}

		return creation.ReturnCode == ReturnCodes().Success && creation.MeetingID == meetingID
	}
	tests := getTests("Create", fmt.Sprintf("name=doe&meetingID=%s&moderatorPW=pwd", meetingID), validResponse, customValidator)

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

	customValidator := func(response interface{}) bool {
		end, ok := response.(*EndResponse)
		if !ok {
			return false
		}

		return end.ReturnCode == ReturnCodes().Success
	}
	tests := getTests("End", fmt.Sprintf("meetingID=%s&password=pwd", meetingID), validResponse, customValidator)

	executeTests(t, "End", tests)
}

func TestIsMeetingRunning(t *testing.T) {
	validResponse := &IsMeetingsRunningResponse{
		ReturnCode: ReturnCodes().Success,
		Running:    true,
	}

	customValidator := func(response interface{}) bool {
		running, ok := response.(*IsMeetingsRunningResponse)
		if !ok {
			return false
		}

		return running.ReturnCode == ReturnCodes().Success && running.Running
	}

	tests := getTests("IsMeetingRunning", fmt.Sprintf("meetingID=%s", meetingID), validResponse, customValidator)

	executeTests(t, "IsMeetingRunning", tests)
}

func TestGetMeetingInfo(t *testing.T) {
	validResponse := &GetMeetingInfoResponse{
		ReturnCode: ReturnCodes().Success,
		MeetingID:  meetingID,
	}

	customValidator := func(response interface{}) bool {
		info, ok := response.(*GetMeetingInfoResponse)
		if !ok {
			return false
		}

		return info.ReturnCode == ReturnCodes().Success && info.MeetingID == meetingID
	}

	tests := getTests("GetMeetingInfo", fmt.Sprintf("meetingID=%s", meetingID), validResponse, customValidator)

	executeTests(t, "GetMeetingInfo", tests)
}
