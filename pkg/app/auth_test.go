package app

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"

	"github.com/SLedunois/b3lb/pkg/api"

	"github.com/SLedunois/b3lb/pkg/config"

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
			Name: "An invalid checksum should returns 200 with checksum error",
			Mock: func() {
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
			Name: "A valid checksum should returns 200 code",
			Mock: func() {
				test.SetRequestParams(c, "name=simon&checksum=8f0378b9dbb7967c7069c418062d4f486b951b6f")
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
