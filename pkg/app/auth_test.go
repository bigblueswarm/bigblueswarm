package app

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/api"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/test_utils/pkg/request"
	"github.com/bigblueswarm/test_utils/pkg/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChecksumValidation(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	server := NewServer(&config.Config{
		BigBlueSwarm: config.BigBlueSwarm{
			Secret: test.DefaultSecret(),
		},
	})
	server.TenantManager = &admin.TenantManagerMock{}

	tests := []test.Test{
		{
			Name: "No checksum should returns 200 with checksum error",
			Mock: func() {
				request.SetRequestParams(c, "name=simon")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				var response api.Error
				e := xml.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, e)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, "checksumError", response.MessageKey)
			},
		},
		{
			Name: "An error returned by tenant manager should returns 500 status code",
			Mock: func() {
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return nil, errors.New("tenant manager error")
				}
				request.SetRequestParams(c, "name=simon&checksum=invalid_checksum")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An unknown tenant should return 403 forbidden",
			Mock: func() {
				request.SetRequestHost(c, "localhost")
				request.SetRequestParams(c, "name=simon&checksum=checksum")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return nil, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusForbidden, w.Code)
			},
		},
		{
			Name: "An invalid checksum should returns 200 with checksum error",
			Mock: func() {
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
					}, nil
				}
				request.SetRequestParams(c, "name=simon&checksum=invalid_checksum")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				var response api.Error
				e := xml.Unmarshal(w.Body.Bytes(), &response)
				assert.Nil(t, e)
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, api.ReturnCodes().Failed, response.ReturnCode)
				assert.Equal(t, "checksumError", response.MessageKey)
			},
		},
		{
			Name: "A valid custom tenant checksum should returns 200 code",
			Mock: func() {
				request.SetRequestParams(c, "name=simon&checksum=f0ce59033b7468b690112cd2e715c698e35c8e2b")
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host:   "localhost",
							Secret: "mydummysecret",
						},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, w.Body.String(), "") //Next handler returns an empty string
			},
		},
		{
			Name: "A valid checksum should returns 200 code",
			Mock: func() {
				request.SetRequestParams(c, "name=simon&checksum=8f0378b9dbb7967c7069c418062d4f486b951b6f")
				request.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: &admin.TenantSpec{
							Host: "localhost",
						},
					}, nil
				}
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, w.Body.String(), "") //Next handler returns an empty string
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			server.ChecksumValidation(c)
			test.Validator(t, nil, nil)
		})
	}
}
