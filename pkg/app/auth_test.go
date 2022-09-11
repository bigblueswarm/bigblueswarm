package app

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/v2/internal/test"

	"github.com/SLedunois/b3lb/v2/pkg/admin"
	"github.com/SLedunois/b3lb/v2/pkg/api"

	"github.com/SLedunois/b3lb/v2/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestChecksumValidation(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	server := NewServer(&config.Config{
		BigBlueButton: config.BigBlueButton{
			Secret: test.DefaultSecret(),
		},
	})
	server.TenantManager = &admin.TenantManagerMock{}

	tests := []test.Test{
		{
			Name: "No checksum should returns 200 with checksum error",
			Mock: func() {
				test.SetRequestParams(c, "name=simon")
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
				test.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return nil, errors.New("tenant manager error")
				}
				test.SetRequestParams(c, "name=simon&checksum=invalid_checksum")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
		{
			Name: "An unknown tenant should return 403 forbidden",
			Mock: func() {
				test.SetRequestHost(c, "localhost")
				test.SetRequestParams(c, "name=simon&checksum=checksum")
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
				test.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: map[string]string{},
					}, nil
				}
				test.SetRequestParams(c, "name=simon&checksum=invalid_checksum")
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
				test.SetRequestParams(c, "name=simon&checksum=a03a5771d5bd9b0930df4c99599a20dab8319226")
				test.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: map[string]string{
							"secret": "dummy_secret",
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
				test.SetRequestParams(c, "name=simon&checksum=8f0378b9dbb7967c7069c418062d4f486b951b6f")
				test.SetRequestHost(c, "localhost")
				admin.GetTenantTenantManagerMockFunc = func(hostname string) (*admin.Tenant, error) {
					return &admin.Tenant{
						Spec: map[string]string{},
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
