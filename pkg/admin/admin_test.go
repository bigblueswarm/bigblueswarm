package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/pkg/api"
	"github.com/SLedunois/b3lb/pkg/balancer"
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/SLedunois/b3lb/pkg/admin/mock"
	bmock "github.com/SLedunois/b3lb/pkg/balancer/mock"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func toBBBInstanceArray(body []byte) []api.BigBlueButtonInstance {
	instances := []api.BigBlueButtonInstance{}
	json.Unmarshal(body, &instances)
	return instances
}

func toInstanceStatusArray(body []byte) []balancer.InstanceStatus {
	status := []balancer.InstanceStatus{}
	json.Unmarshal(body, &status)
	return status
}

func TestListInstances(t *testing.T) {
	url := "http://localhost/bigbluebutton"
	var w *httptest.ResponseRecorder
	admin := CreateAdmin(&mock.InstanceManager{}, &bmock.Balancer{}, &config.AdminConfig{})

	tests := []test.Test{
		{
			Name: "List should returns a list containg a single bigbluebutton instance",
			Mock: func() {
				mock.ListInstancesFunc = func() ([]api.BigBlueButtonInstance, error) {
					return []api.BigBlueButtonInstance{
						{
							URL:    url,
							Secret: test.DefaultSecret(),
						},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, toBBBInstanceArray(w.Body.Bytes())[0].URL, url)
				assert.Equal(t, toBBBInstanceArray(w.Body.Bytes())[0].Secret, test.DefaultSecret())
			},
		},
		{
			Name: "List should return an internal server error if instance manager returns an error",
			Mock: func() {
				mock.ListInstancesFunc = func() ([]api.BigBlueButtonInstance, error) {
					return nil, errors.New("unexpected error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			test.Mock()
			admin.ListInstances(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestClusterStatus(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&mock.InstanceManager{}, &bmock.Balancer{}, &config.AdminConfig{})

	host := "http://localhost/bigbluebutton"
	cpu := 20.01
	mem := 35.45
	apiStatus := "Up"
	activeMeetings := int64(3)
	activeParticipants := int64(22)

	expectedStatus := []balancer.InstanceStatus{
		{
			Host:               host,
			CPU:                cpu,
			Mem:                mem,
			APIStatus:          apiStatus,
			ActiveMeeting:      int64(activeMeetings),
			ActiveParticipants: int64(activeParticipants),
		},
	}

	tests := []test.Test{
		{
			Name: "an error returned by InstanceManager should return a 500 Internal Server Error status code",
			Mock: func() {
				mock.ListFunc = func() ([]string, error) {
					return nil, errors.New("manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "an error returned by Balancer should return a 500 Internal Server Error status code",
			Mock: func() {
				mock.ListFunc = func() ([]string, error) {
					return []string{}, nil
				}
				bmock.BalancerClusterStatusFunc = func(instances []string) ([]balancer.InstanceStatus, error) {
					return nil, errors.New("balancer error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "a valid request should return a 200 Status OK and a list of status",
			Mock: func() {
				mock.ListFunc = func() ([]string, error) {
					return []string{}, nil
				}
				bmock.BalancerClusterStatusFunc = func(instances []string) ([]balancer.InstanceStatus, error) {
					return expectedStatus, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				status := toInstanceStatusArray(w.Body.Bytes())
				assert.Equal(t, host, status[0].Host)
				assert.Equal(t, cpu, status[0].CPU)
				assert.Equal(t, mem, status[0].Mem)
				assert.Equal(t, apiStatus, status[0].APIStatus)
				assert.Equal(t, activeMeetings, int64(status[0].ActiveMeeting))
				assert.Equal(t, activeParticipants, int64(status[0].ActiveParticipants))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.ClusterStatus(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestSetInstances(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&mock.InstanceManager{}, &bmock.Balancer{}, &config.AdminConfig{})

	tests := []test.Test{
		{
			Name: "an invalid body should return a bad request status and an error",
			Mock: func() {
				test.AddRequestBody(c, "")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "Body does not bind InstanceList object: EOF", w.Body.String())
			},
		},
		{
			Name: "an error returned by InstanceManager should return an internal server error and an error",
			Mock: func() {
				test.AddRequestBody(c, `kind: InstanceList
instances:
  http://bigbluebutton1: secret1`)
				mock.SetInstancesFunc = func(instances map[string]string) error {
					return errors.New("instance manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "instance manager error", w.Body.String())
			},
		},
		{
			Name: "a valid request should return a http 200 ok",
			Mock: func() {
				test.AddRequestBody(c, `kind: InstanceList
instances:
  http://bigbluebutton1: secret1`)
				mock.SetInstancesFunc = func(instances map[string]string) error {
					return nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusCreated, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.SetInstances(c)
			test.Validator(t, nil, nil)
		})
	}
}
