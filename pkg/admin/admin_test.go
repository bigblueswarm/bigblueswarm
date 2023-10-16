package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v3/pkg/api"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/balancer"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/config"

	"github.com/bigblueswarm/test_utils/pkg/request"
	"github.com/bigblueswarm/test_utils/pkg/test"

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
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	tests := []test.Test{
		{
			Name: "List should returns a list containg a single bigbluebutton instance",
			Mock: func() {
				ListInstancesInstanceManagerMockFunc = func() ([]api.BigBlueButtonInstance, error) {
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
				ListInstancesInstanceManagerMockFunc = func() ([]api.BigBlueButtonInstance, error) {
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
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	host := "http://localhost/bigbluebutton"
	cpu := 20.01
	mem := 35.45
	apiStatus := "Up"
	meetings := int64(3)
	participants := int64(22)

	expectedStatus := []balancer.InstanceStatus{
		{
			Host:         host,
			CPU:          cpu,
			Mem:          mem,
			APIStatus:    apiStatus,
			Meetings:     int64(meetings),
			Participants: int64(participants),
		},
	}

	tests := []test.Test{
		{
			Name: "an error returned by InstanceManager should return a 500 Internal Server Error status code",
			Mock: func() {
				ListInstanceManagerMockFunc = func() ([]string, error) {
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
				ListInstanceManagerMockFunc = func() ([]string, error) {
					return []string{}, nil
				}
				balancer.BalancerMockClusterStatusFunc = func(instances []string) ([]balancer.InstanceStatus, error) {
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
				ListInstanceManagerMockFunc = func() ([]string, error) {
					return []string{}, nil
				}
				balancer.BalancerMockClusterStatusFunc = func(instances []string) ([]balancer.InstanceStatus, error) {
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
				assert.Equal(t, meetings, int64(status[0].Meetings))
				assert.Equal(t, participants, int64(status[0].Participants))
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
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	tests := []test.Test{
		{
			Name: "an invalid body should return a bad request status and an error",
			Mock: func() {
				request.AddRequestBody(c, "")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "body does not bind InstanceList object: EOF", w.Body.String())
			},
		},
		{
			Name: "an error returned by InstanceManager should return an internal server error and an error",
			Mock: func() {
				request.AddRequestBody(c, `{
	"kind": "InstanceList",
	"instances": {
		"http://bigbluebutton1": "secret1"
	}
}`)
				SetInstancesInstanceManagerMockFunc = func(instances map[string]string) error {
					return errors.New("instance manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "failed to set instances in instance manager: instance manager error", w.Body.String())
			},
		},
		{
			Name: "a valid request should return a http 200 ok",
			Mock: func() {
				request.AddRequestBody(c, `{
	"kind": "InstanceList",
	"instances": {
		"http://bigbluebutton1": "secret1"
	}
}`)
				SetInstancesInstanceManagerMockFunc = func(instances map[string]string) error {
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

func TestCreateTenant(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	tests := []test.Test{
		{
			Name: "an invalid body should return a bad request status and an error",
			Mock: func() {
				request.AddRequestBody(c, "")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "body does not bind Tenant object: EOF", w.Body.String())
			},
		},
		{
			Name: "an error returned by tenant manager should return an internal server error",
			Mock: func() {
				request.AddRequestBody(c, `{
	"kind": "Tenant",
	"spec": {
  		"host": "localhost:8090"
	},
	"instances": []
}`)
				AddTenantTenantManagerMockFunc = func(tenant *Tenant) error {
					return errors.New("manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "failed to add tenant in tenant manager: manager error", w.Body.String())
			},
		},
		{
			Name: "a valid request should return a 201 created status",
			Mock: func() {
				request.AddRequestBody(c, `{
	"kind": "Tenant",
	"spec": {
			"host": "localhost:8090"
	},
	"instances": []
}`)
				AddTenantTenantManagerMockFunc = func(tenant *Tenant) error {
					return nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusCreated, w.Code)
			},
		},
		{
			Name: "an error returned by tenant manager should return an internal server error",
			Mock: func() {
				request.AddRequestBody(c, `{
	"kind": "Tenant",
	"spec": {},
	"instances": []
}`)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				assert.Equal(t, "failed to create tenant. Tenant spec host should not be null", w.Body.String())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.CreateTenant(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestListTenantsHandler(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	tests := []test.Test{
		{
			Name: "an error returned by tenant manager should return an HTTP 500 - Internal Server Error - and a log",
			Mock: func() {
				ListTenantsTenantManagerMockFunc = func() ([]TenantListObject, error) {
					return []TenantListObject{}, errors.New("manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "unable to list all tenants: manager error", w.Body.String())
			},
		},
		{
			Name: "a valid request should return an HTTP 200 - OK - and a valid response",
			Mock: func() {
				ListTenantsTenantManagerMockFunc = func() ([]TenantListObject, error) {
					return []TenantListObject{
						{
							Hostname:      "localhost:8090",
							InstanceCount: 0,
						},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, `{"kind":"TenantList","tenants":[{"hostname":"localhost:8090","instance_count":0}]}`, w.Body.String())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.ListTenants(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})

	tests := []test.Test{
		{
			Name: "an invalid hostname query parameter should return a HTTP 400 - BadRequest - status code",
			Mock: func() {},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			Name: "an error occurred while getting tenant should return a HTTP 500 - Internal Server Error",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return nil, errors.New("redis error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "failed to retrieve tenant: redis error", w.Body.String())
			},
		},
		{
			Name: "no tenant found should return a HTTP 404 - Not Found error",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return nil, nil
				}

			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusNotFound, w.Code)
				assert.Equal(t, "tenant not found for deletion", w.Body.String())
			},
		},
		{
			Name: "an error returned by TenantManager should return a HTTP 500 - InternalServerError - status code and a log",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return &Tenant{}, nil
				}
				DeleteTenantTenantManagerMockFunc = func(hostname string) error {
					return errors.New("manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "unable to delete tenant", w.Body.String())
			},
		},
		{
			Name: "a valid request should return a HTTP 204 - No Content",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				DeleteTenantTenantManagerMockFunc = func(hostname string) error {
					return nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusNoContent, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.DeleteTenant(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestGetConfigurtionHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	config := &config.Config{
		BigBlueSwarm: config.BigBlueSwarm{
			Secret:                 "secret",
			RecordingsPollInterval: "5m",
		},
		Admin: config.AdminConfig{
			APIKey: "key",
		},
		Balancer: config.BalancerConfig{
			MetricsRange: "-5m",
			CPULimit:     99,
			MemLimit:     99,
		},
		Port: 8090,
		RDB: config.RDB{
			Address:  "http://localhost:8500",
			Password: "",
			DB:       0,
		},
		IDB: config.IDB{
			Address:      "http://localhost:8086",
			Token:        "token",
			Organization: "bigblueswarm",
			Bucket:       "bucket",
		},
	}
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, config)

	expected, err := json.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	admin.GetConfiguration(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expected), w.Body.String())
}

func TestGetTenantHandler(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{})
	tests := []test.Test{
		{
			Name: "an error returned by tenant manager should end with a HTTP 500 - Internal Server Error",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return nil, errors.New("manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
				assert.Equal(t, "unable to retrieve tenant", w.Body.String())
			},
		},
		{
			Name: "no tenant found should ends with a HTTP 404 - Not Found",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return nil, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
		{
			Name: "a valid tenant should return a tenant as JSON",
			Mock: func() {
				c.Params = gin.Params{
					{
						Key:   "hostname",
						Value: "localhost",
					},
				}
				GetTenantTenantManagerMockFunc = func(hostname string) (*Tenant, error) {
					return &Tenant{
						Kind: "Tenant",
						Spec: &TenantSpec{
							Host: "localhost",
						},
						Instances: []string{},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, `{"kind":"Tenant","spec":{"host":"localhost"},"instances":[]}`, w.Body.String())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.GetTenant(c)
			test.Validator(t, nil, nil)
		})
	}
}
