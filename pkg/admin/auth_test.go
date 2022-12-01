package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/balancer"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"

	"github.com/bigblueswarm/test_utils/pkg/request"
	"github.com/bigblueswarm/test_utils/pkg/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApiKeyValidation(t *testing.T) {
	var w *httptest.ResponseRecorder
	var c *gin.Context
	admin := CreateAdmin(&InstanceManagerMock{}, &TenantManagerMock{}, &balancer.Mock{}, &config.Config{Admin: config.AdminConfig{APIKey: test.DefaultAPIKey()}})
	tests := []test.Test{
		{
			Name: "An empty api key should returns an unauthorized error",
			Mock: func() {
				request.SetRequestHeader(c, "Authorization", "")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			Name: "An invalid api key should returns an unauthorized error",
			Mock: func() {
				request.SetRequestHeader(c, "Authorization", "invalid_key")
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, http.StatusUnauthorized, w.Code)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
			test.Mock()
			admin.APIKeyValidation(c)
			test.Validator(t, nil, nil)
		})
	}
}
