package admin

import (
	TestUtil "b3lb/internal/test"
	"io"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestApiKeyValidation(t *testing.T) {
	type test struct {
		name           string
		headers        map[string]string
		expectedStatus int
		body           io.Reader
	}

	tests := []test{
		{
			name: "An invalid api key should returns an unauthorized error",
			headers: map[string]string{
				"Authorization": TestUtil.DefaultAPIKey() + "dummy",
			},
			expectedStatus: 401,
			body:           nil,
		},
		{
			name:           "An empty api key should returns an unauthorized error",
			headers:        map[string]string{},
			expectedStatus: 401,
			body:           nil,
		},
		{
			name: "A valid api key should go through the api key validation middleware",
			headers: map[string]string{
				"Authorization": TestUtil.DefaultAPIKey(),
			},
			expectedStatus: 400,
			body:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := TestUtil.ExecuteRequestWithHeaders(router, "POST", "/admin/servers", test.body, test.headers)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
