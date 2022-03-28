package app

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/stretchr/testify/assert"

	"github.com/SLedunois/b3lb/pkg/admin"
	"github.com/SLedunois/b3lb/pkg/api"
	RestClientMock "github.com/SLedunois/b3lb/pkg/restclient/mock"

	log "github.com/sirupsen/logrus"
	LogTest "github.com/sirupsen/logrus/hooks/test"
)

func TestToDuration(t *testing.T) {
	t.Run("A valid value should return a valid duration", func(t *testing.T) {
		assert.Equal(t, time.Duration(1*time.Minute), toDuration("1m"))
	})

	t.Run("An invalid value should panic toDuration", func(t *testing.T) {
		assert.Panics(t, func() {
			toDuration("invalid")
		})
	})
}

func TestPollRecordings(t *testing.T) {
	logHook := LogTest.NewGlobal()
	log.AddHook(logHook)
	tests := []test.Test{
		{
			Name: "An error returned by the clear recordings method should be logged",
			Mock: func() {
				mock := redisMock.ExpectKeys(RecodingPattern())
				mock.SetVal([]string{})
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, "Failed to clear recordings. redis error", logHook.LastEntry().Message)
			},
		},
		{
			Name: "An error returned by the list instances method should be logged",
			Mock: func() {
				redisMock.ExpectKeys(RecodingPattern()).SetVal([]string{})
				admin.ListInstancesInstanceManagerMockFunc = func() ([]api.BigBlueButtonInstance, error) {
					return nil, errors.New("admin error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, "Failed to retrieve instances. admin error", logHook.LastEntry().Message)
			},
		},
		{
			Name: "An error returned by the instance get recordings method should be logged",
			Mock: func() {
				redisMock.ExpectKeys(RecodingPattern()).SetVal([]string{})
				admin.ListInstancesInstanceManagerMockFunc = func() ([]api.BigBlueButtonInstance, error) {
					return []api.BigBlueButtonInstance{
						{
							URL:    "http://localhost:8080/bigbluebutton",
							Secret: test.DefaultSecret(),
						},
					}, nil
				}
				RestClientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("rest client error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, "Failed to retrieve recordings for instance http://localhost:8080/bigbluebutton. rest client error", logHook.LastEntry().Message)
			},
		},
		{
			Name: "An error returned by the mapper add method should be logged",
			Mock: func() {
				redisMock.ExpectKeys(RecodingPattern()).SetVal([]string{})
				admin.ListInstancesInstanceManagerMockFunc = func() ([]api.BigBlueButtonInstance, error) {
					return []api.BigBlueButtonInstance{
						{
							URL:    "http://localhost:8080/bigbluebutton",
							Secret: test.DefaultSecret(),
						},
					}, nil
				}
				RestClientMock.DoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := api.GetRecordingsResponse{
						Recordings: []api.Recording{
							{
								RecordID: "recording-id",
							},
						},
					}

					value, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(value)),
					}, nil
				}
				mock := redisMock.ExpectSet(RecordingMapKey("recording-id"), "http://localhost:8080/bigbluebutton", 0)
				mock.SetVal("")
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, "Failed to store recording recording-id. redis error", logHook.LastEntry().Message)
			},
		},
	}
	server := doGenericInitialization()
	server.InstanceManager = &admin.InstanceManagerMock{}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			server.pollRecordings()
			test.Validator(t, nil, nil)
		})
	}
}
