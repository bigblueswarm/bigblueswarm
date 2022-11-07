package app

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/v2/pkg/admin"
	"github.com/SLedunois/b3lb/v2/pkg/api"
	"github.com/SLedunois/b3lb/v2/pkg/balancer"
	"github.com/b3lb/test_utils/pkg/request"
	"github.com/b3lb/test_utils/pkg/test"

	"github.com/SLedunois/b3lb/v2/pkg/config"
	"github.com/SLedunois/b3lb/v2/pkg/restclient"
	log "github.com/sirupsen/logrus"
	LogTest "github.com/sirupsen/logrus/hooks/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	meetingID = "meeting_id"
	params    = fmt.Sprintf("meetingID=%s", meetingID)
	instance  = "http://localhost:8080/bigbluebutton"
	w         *httptest.ResponseRecorder
	c         *gin.Context
)

func doGenericInitialization() *Server {
	server := NewServer(&config.Config{})
	server.Mapper = mapper
	server.InstanceManager = instanceManager
	server.TenantManager = &admin.TenantManagerMock{}
	server.Balancer = &balancer.Mock{}
	restclient.Client = &restclient.Mock{}

	return server
}

func unMarshallError(body []byte) api.Error {
	var response api.Error
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallIsMeetingRunningResponse(body []byte) api.IsMeetingsRunningResponse {
	var response api.IsMeetingsRunningResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallEndResponse(body []byte) api.EndResponse {
	var response api.EndResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallGetMeetingInfoResponse(body []byte) api.GetMeetingInfoResponse {
	var response api.GetMeetingInfoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallCreateResponse(body []byte) api.CreateResponse {
	var response api.CreateResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallGetMeetingsResponse(body []byte) api.GetMeetingsResponse {
	var response api.GetMeetingsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallJoinRedirectResponse(body []byte) api.JoinRedirectResponse {
	var response api.JoinRedirectResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallGetRecordingsResponse(body []byte) api.GetRecordingsResponse {
	var response api.GetRecordingsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallUpdateRecordingsResponse(body []byte) api.UpdateRecordingsResponse {
	var response api.UpdateRecordingsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallDeleteRecordingsResponse(body []byte) api.DeleteRecordingsResponse {
	var response api.DeleteRecordingsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallPublishRecordingsResponse(body []byte) api.PublishRecordingsResponse {
	var response api.PublishRecordingsResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallGetrecordingsTextTracksResponse(body []byte) api.GetRecordingsTextTracksResponse {
	var response api.GetRecordingsTextTracksResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func unMarshallJSONError(body []byte) api.JSONResponse {
	var response api.JSONResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response
}

func TestHealthCheckRoute(t *testing.T) {
	// Healthcheck has a single test. The method always returns success and the same response.
	t.Run("Healtcheck should returns a valid response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		server := NewServer(&config.Config{})
		server.HealthCheck(c)

		var response api.HealthCheck
		err := xml.Unmarshal(w.Body.Bytes(), &response)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
		assert.Equal(t, "2.0", response.Version)
	})
}

func TestCreate(t *testing.T) {
	creationParams := fmt.Sprintf("%s&name=test_name&attendeePW=pwd&moderatorPW=pwd2", params)
	tests := []test.Test{
		{
			Name: "An error thrown by the TenantManager while getting active tenant should return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestHost(c, "localhost")
				request.SetRequestParams(c, creationParams)
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return nil, errors.New("tenant manager error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "No instances found by InstanceManager shold return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Instances: []string{},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An error thrown by Balancer while processing should returns an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
						Instances: []string{
							"http://localhost/bigbuebutton",
						},
					}, nil
				}
				balancer.BalancerMockProcessFunc = func(instances []string) (string, error) {
					return "", errors.New("balancer error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An error thrown by InstanceManager while getting target instance should return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
						Instances: []string{
							"http://localhost/bigbuebutton",
						},
					}, nil
				}
				balancer.BalancerMockProcessFunc = func(instances []string) (string, error) {
					return instance, nil
				}
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An error thrown by BigBlueButton instance while creation should return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
						Instances: []string{
							"http://localhost/bigbuebutton",
						},
					}, nil
				}
				balancer.BalancerMockProcessFunc = func(instances []string) (string, error) {
					return instance, nil
				}
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("bbb error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An error thrown by Mapper while adding session should return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
						Instances: []string{
							"http://localhost/bigbuebutton",
						},
					}, nil
				}
				balancer.BalancerMockProcessFunc = func(instances []string) (string, error) {
					return instance, nil
				}
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					createResponse := &api.CreateResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
						MeetingID:   meetingID,
						AttendeePW:  "pwd",
						ModeratorPW: "pwd2",
					}

					response, err := xml.Marshal(createResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
				redisMock.ExpectSet(MeetingMapKey(meetingID), instance, 0).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A valid request should return a valid response",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: creationParams,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, creationParams)
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
						Instances: []string{
							"http://localhost/bigbuebutton",
						},
					}, nil
				}
				balancer.BalancerMockProcessFunc = func(instances []string) (string, error) {
					return instance, nil
				}
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					createResponse := &api.CreateResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
						MeetingID:   meetingID,
						AttendeePW:  "pwd",
						ModeratorPW: "pwd2",
					}

					response, err := xml.Marshal(createResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
				redisMock.ExpectSet(MeetingMapKey(meetingID), instance, 0).SetVal(meetingID)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallCreateResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, meetingID, response.MeetingID)
				assert.Equal(t, "pwd", response.AttendeePW)
				assert.Equal(t, "pwd2", response.ModeratorPW)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			server := doGenericInitialization()
			server.Create(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []test.Test{
		{
			Name: "No provided meeting id should returns an empty meeting id reponse",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: "",
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().ValidationError, response.MessageKey)
				assert.Equal(t, api.Messages().EmptyMeetingID, response.Message)
			},
		},
		{
			Name: "Providing a meeting id that does not exists should return a not found error",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal("")
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
				assert.Equal(t, api.Messages().NotFound, response.Message)
			},
		},
		{
			Name: "Providing a meeting id that does not match an instance should return a not found error",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal("")
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
				assert.Equal(t, api.Messages().NotFound, response.Message)
			},
		},
		{
			Name: "A valid request should redirect to the meeting url",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code) // Test returs a 200 status instead of 302
				assert.Equal(t, w.Header().Get("Location"), fmt.Sprintf("%s/api/join?meetingID=%s&checksum=3326f2a7090212891651d6da31a608ec88f3ca58", instance, meetingID))
			},
		},
		{
			Name: "An error return by BigBlueButton instance while calling join api with `redirect=false` parameter set should return an internal server error status code",
			Mock: func() {
				request.SetRequestParams(c, fmt.Sprintf("%s&redirect=false", params))
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: fmt.Sprintf("%s&redirect=false", params),
					Action: api.Join,
				}
				c.Set("api_ctx", checksum)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("instance error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "Calling join api with `redirect=false` parameter set should return a valid JoinRedirectResponse",
			Mock: func() {
				request.SetRequestParams(c, fmt.Sprintf("%s&redirect=false", params))
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: fmt.Sprintf("%s&redirect=false", params),
					Action: api.Join,
				}
				c.Set("api_ctx", checksum)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					joinResponse := &api.JoinRedirectResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
						URL: "https://localhost:8080/html5client/join?sessionToken=ai1wqj8wb6s7rnk0",
					}

					response, err := xml.Marshal(joinResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallJoinRedirectResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, "https://localhost:8080/html5client/join?sessionToken=ai1wqj8wb6s7rnk0", response.URL)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			server := doGenericInitialization()
			server.Join(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestEnd(t *testing.T) {
	/*
		TestIsMeetingsRunning tests the following cases:
		- No provided meeting id should returns an empty meeting id reponse
		- Providing a meeting id that does not exists should return a not found error
		- Providing a meeting id that does not match an instance should return a not found error
		- An error thrown by remote bigbluebutton instance should return an http internal error status
	*/

	tests := []test.Test{
		{
			Name: "An error thrown by Mapper removing session should return an http interal server error",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				redisMock.ExpectDel(MeetingMapKey(meetingID)).SetErr(errors.New("error"))
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					endResponse := &api.EndResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
					}

					response, err := xml.Marshal(endResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A valid end call should return a success response",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				redisMock.ExpectDel(MeetingMapKey(meetingID)).SetVal(1)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					endResponse := &api.EndResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
					}

					response, err := xml.Marshal(endResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallEndResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			server := doGenericInitialization()
			server.End(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestIsMeetingRunning(t *testing.T) {
	tests := []test.Test{
		{
			Name: "No provided meeting id should returns an empty meeting id reponse",
			Mock: func() {
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: "",
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().ValidationError, response.MessageKey)
				assert.Equal(t, api.Messages().EmptyMeetingID, response.Message)
			},
		},
		{
			Name: "Providing a meeting id that does not exists should return a not found error",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal("")
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
				assert.Equal(t, api.Messages().NotFound, response.Message)
			},
		},
		{
			Name: "Providing a meeting id that does not match an instance should return a not found error",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal("")
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, response.MessageKey)
				assert.Equal(t, api.Messages().NotFound, response.Message)
			},
		},
		{
			Name: "An error thrown by remote bigbluebutton instance should return an http internal error status",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("http error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A valid request should return a valid response",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					isRunningReponse := &api.IsMeetingsRunningResponse{
						Running:    true,
						ReturnCode: api.ReturnCodes().Success,
					}

					response, err := xml.Marshal(isRunningReponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallIsMeetingRunningResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, true, response.Running)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			server := doGenericInitialization()
			server.IsMeetingRunning(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestGetMeetingInfo(t *testing.T) {
	/*
		TestIsMeetingsRunning tests the following cases:
		- No provided meeting id should returns an empty meeting id reponse
		- Providing a meeting id that does not exists should return a not found error
		- Providing a meeting id that does not match an instance should return a not found error
		- An error thrown by remote bigbluebutton instance should return an http internal error status
	*/
	tests := []test.Test{
		{
			Name: "A valid end call should return a success response",
			Mock: func() {
				request.SetRequestParams(c, params)
				redisMock.ExpectGet(MeetingMapKey(meetingID)).SetVal(instance)
				redisMock.ExpectHGet(admin.B3LBInstances, instance).SetVal(test.DefaultSecret())
				checksum := &api.Checksum{
					Secret: test.DefaultSecret(),
					Params: params,
					Action: api.IsMeetingRunning,
				}
				c.Set("api_ctx", checksum)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					meetingResponse := &api.GetMeetingInfoResponse{
						ReturnCode: api.ReturnCodes().Success,
						MeetingInfo: api.MeetingInfo{
							MeetingID: meetingID,
						},
					}

					response, err := xml.Marshal(meetingResponse)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetMeetingInfoResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, response.MeetingID, meetingID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.GetMeetingInfo(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestGetMeetings(t *testing.T) {
	tests := []test.Test{
		{
			Name: "An error thrown by instance manager should return an http internal error status",
			Mock: func() {
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A valid request should return a valid response",
			Mock: func() {
				instances := map[string]string{
					"http://localhost/bigbluebutton": test.DefaultSecret(),
				}
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetVal(instances)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					meetings := &api.GetMeetingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Meetings: []api.MeetingInfo{
							{
								MeetingID: "meeting-id",
							},
						},
					}

					response, err := xml.Marshal(meetings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetMeetingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, 1, len(response.Meetings))
				assert.Equal(t, "meeting-id", response.Meetings[0].MeetingID)
			},
		},
		{
			Name: "An error thrown by a remte instance should not block the response",
			Mock: func() {
				instances := map[string]string{
					"http://localhost/bigbluebutton":      test.DefaultSecret(),
					"http://localhost:8080/bigbluebutton": test.DefaultSecret(),
				}
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetVal(instances)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					if req.URL.Host == "localhost:8080" {
						return nil, errors.New("remote error")
					}

					meetings := &api.GetMeetingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Meetings: []api.MeetingInfo{
							{
								MeetingID: "meeting-id",
							},
						},
					}

					response, err := xml.Marshal(meetings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetMeetingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, 1, len(response.Meetings))
				assert.Equal(t, "meeting-id", response.Meetings[0].MeetingID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.GetMeetings(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestGetRecodings(t *testing.T) {
	checksum := &api.Checksum{
		Secret: test.DefaultSecret(),
		Params: "",
		Action: api.GetRecordings,
	}

	instances := map[string]string{
		"http://localhost/bigbluebutton": test.DefaultSecret(),
	}

	logHook := LogTest.NewGlobal()
	log.AddHook(logHook)

	tests := []test.Test{
		{
			Name: "An error returned by the instance manager ListInstance method should return a no recordings response",
			Mock: func() {
				c.Set("api_ctx", checksum)
				mock := redisMock.ExpectHGetAll(admin.B3LBInstances)
				mock.SetErr(errors.New("redis error"))
				mock.SetVal(map[string]string{})
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NoRecordings, response.MessageKey)
				assert.Equal(t, api.Messages().NoRecordings, response.Message)
				assert.Equal(t, 0, len(response.Recordings))
			},
		},
		{
			Name: "An error returned by the remote instance should be logged",
			Mock: func() {
				c.Set("api_ctx", checksum)
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetVal(instances)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("remote error")
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, "Instance http://localhost/bigbluebutton failed to retrieve recordings. remote error", logHook.LastEntry().Message)
			},
		},
		{
			Name: "An empty response from the remote instance should return a no recordings response",
			Mock: func() {
				c.Set("api_ctx", checksum)
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetVal(instances)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := api.GetRecordingsResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
							MessageKey: api.MessageKeys().NoRecordings,
							Message:    api.Messages().NoRecordings,
						},
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, api.MessageKeys().NoRecordings, response.MessageKey)
				assert.Equal(t, api.Messages().NoRecordings, response.Message)
				assert.Equal(t, 0, len(response.Recordings))
			},
		},
		{
			Name: "An non empty response from the remote instance should return a valid response",
			Mock: func() {
				c.Set("api_ctx", checksum)
				redisMock.ExpectHGetAll(admin.B3LBInstances).SetVal(instances)
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := api.GetRecordingsResponse{
						Response: api.Response{
							ReturnCode: api.ReturnCodes().Success,
						},
						Recordings: []api.Recording{
							{
								RecordID:  "record-id",
								MeetingID: "meeting-id",
							},
						},
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				response := unMarshallGetRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, 1, len(response.Recordings))
				assert.Equal(t, "record-id", response.Recordings[0].RecordID)
				assert.Equal(t, "meeting-id", response.Recordings[0].MeetingID)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.GetRecordings(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestUpdateRecordings(t *testing.T) {
	tests := []test.Test{
		{
			Name: "A missing recordID parameter should return a missing parameter response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "",
					Secret: test.DefaultSecret(),
					Action: api.UpdateRecordings,
				}

				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				err := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, err.ReturnCode)
				assert.Equal(t, api.MessageKeys().MissingRecordIDParameter, err.MessageKey)
				assert.Equal(t, api.Messages().MissingRecordIDParameter, err.Message)
			},
		},
		{
			Name: "An error returned by the mapper should return a failed response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.UpdateRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				mock := redisMock.ExpectGet(RecordingMapKey("record-id"))
				mock.SetErr(errors.New("redis error"))
				mock.SetVal("")
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				err := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, err.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, err.MessageKey)
				assert.Equal(t, api.Messages().RecordingNotFound, err.Message)
			},
		},
		{
			Name: "An error returned by the instance manager should return a failed response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.UpdateRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				mock := redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton")
				mock.SetVal("")
				mock.SetErr(errors.New("redis error"))
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				err := unMarshallError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, err.ReturnCode)
				assert.Equal(t, api.MessageKeys().NotFound, err.MessageKey)
				assert.Equal(t, api.Messages().RecordingNotFound, err.Message)
			},
		},
		{
			Name: "An error returned by the remote instance should return an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.UpdateRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("http error")
				}
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A successful update should return a a valid UpdateRecordings response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.UpdateRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := &api.UpdateRecordingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Updated:    true,
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				response := unMarshallUpdateRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, true, response.Updated)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.UpdateRecordings(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestDeleteRecordings(t *testing.T) {
	// Because the DeleteRecordings uses the proxyRecordings method that is already tested by UpdateRecording,
	// DeleteRecordings test will only test the end process method and a valid test case
	tests := []test.Test{
		{
			Name: "A valid request with deleted=false should not delete the recording",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.DeleteRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := &api.DeleteRecordingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Deleted:    false,
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				response := unMarshallDeleteRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, false, response.Deleted)
			},
		},
		{
			Name: "An error returned by the mapper when deleting the recordings should returns an internal server error",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.DeleteRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := &api.DeleteRecordingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Deleted:    true,
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
				mock := redisMock.ExpectDel(RecordingMapKey("record-id"))
				mock.SetErr(errors.New("redis error"))
				mock.SetVal(0)
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "A valid request should return an http 200 code and a valid response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.DeleteRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := &api.DeleteRecordingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Deleted:    true,
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
				redisMock.ExpectDel(RecordingMapKey("record-id")).SetVal(1)
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				response := unMarshallDeleteRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, true, response.Deleted)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.DeleteRecordings(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestPublishRecordings(t *testing.T) {
	// Because PublishRecordings only use proxyRecordings method
	// we test the success scenario
	tests := []test.Test{
		{
			Name: "A valid request should return an http 200 code and a valid response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id&published=true",
					Secret: test.DefaultSecret(),
					Action: api.PublishRecordings,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id&published=true")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					recordings := &api.PublishRecordingsResponse{
						ReturnCode: api.ReturnCodes().Success,
						Published:  true,
					}

					response, err := xml.Marshal(recordings)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				response := unMarshallPublishRecordingsResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, response.ReturnCode)
				assert.Equal(t, true, response.Published)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.PublishRecordings(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestGetRecordingsTextTracks(t *testing.T) {
	// Because GetRecordingsTextTracks only use proxyRecordings method
	// we test the success scenario
	tests := []test.Test{
		{
			Name: "A valid request should return an http 200 code and a valid response",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.GetRecordingsTextTracks,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost:8080/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost:8080/bigbluebutton").SetVal(test.DefaultSecret())
				restclient.RestClientMockDoFunc = func(req *http.Request) (*http.Response, error) {
					tracks := &api.GetRecordingsTextTracksResponse{
						Response: api.RecordingsTextTrackResponseType{
							ReturnCode: api.ReturnCodes().Success,
							Tracks: []api.Track{
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

					response, err := json.Marshal(tracks)
					if err != nil {
						panic(err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader(response)),
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, _ error) {
				tracks := unMarshallGetrecordingsTextTracksResponse(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Success, tracks.Response.ReturnCode)
				assert.Equal(t, len(tracks.Response.Tracks), 1)
				assert.Equal(t, tracks.Response.Tracks[0].Href, "http://localhost/client/api/v1/recordings/recording-id/texttracks/track-id")
				assert.Equal(t, tracks.Response.Tracks[0].Kind, "subtitles")
				assert.Equal(t, tracks.Response.Tracks[0].Lang, "en")
				assert.Equal(t, tracks.Response.Tracks[0].Label, "English")
				assert.Equal(t, tracks.Response.Tracks[0].Source, "upload")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.GetRecordingsTextTracks(c)
			test.Validator(t, nil, nil)
		})
	}
}

func TestPutRecordingTextTrack(t *testing.T) {
	tests := []test.Test{
		{
			Name: "No recording identifier should return a param error",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordinID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.GetRecordingsTextTracks,
				}

				c.Set("api_ctx", checksum)
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				e := unMarshallJSONError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.MessageKeys().ParamError, e.Response.MessageKey)
				assert.Equal(t, api.Messages().MissingRecordIDParameter, e.Response.Message)
			},
		},
		{
			Name: "An error returned by the mapper should return a not found error",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.GetRecordingsTextTracks,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				mock := redisMock.ExpectGet(RecordingMapKey("record-id"))
				mock.SetVal("")
				mock.SetErr(errors.New("redis error"))

			},
			Validator: func(t *testing.T, value interface{}, err error) {
				e := unMarshallJSONError(w.Body.Bytes())
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.MessageKeys().NoRecordings, e.Response.MessageKey)
				assert.Equal(t, api.Messages().RecordingTextTrackNotFound, e.Response.Message)
			},
		},
		{
			Name: "A valid request should redirect to the right bigbluebutton instance with a valid url",
			Mock: func() {
				checksum := &api.Checksum{
					Params: "recordID=record-id",
					Secret: test.DefaultSecret(),
					Action: api.PutRecordingTextTrack,
				}

				c.Set("api_ctx", checksum)
				request.SetRequestParams(c, "recordID=record-id")

				redisMock.ExpectGet(RecordingMapKey("record-id")).SetVal("http://localhost/bigbluebutton")
				redisMock.ExpectHGet(admin.B3LBInstances, "http://localhost/bigbluebutton").SetVal(test.DefaultSecret())
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusFound, w.Code)
				assert.Equal(t, "http://localhost/bigbluebutton/api/putRecordingTextTrack?recordID=record-id&checksum=d770e21b9c3077df658f9162a95d1a11cb35854a", c.Writer.Header().Get("Location"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			server := doGenericInitialization()
			test.Mock()
			server.PutRecordingTextTrack(c)
			test.Validator(t, nil, nil)
		})
	}
}
