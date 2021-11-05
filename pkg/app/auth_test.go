package app

import (
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"encoding/xml"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumValidation(t *testing.T) {

	type test struct {
		name         string
		url          string
		expectedBody string
	}

	router := launchRouter(&config.Config{
		BigBlueButton: config.BigBlueButton{
			Secret: "supersecret",
		},
	})

	bytes, err := xml.Marshal(api.DefaultChecksumError())

	if err != nil {
		t.Fatal(err)
		return
	}

	tests := []test{
		{
			name:         "No checksum returns 200 with checksum error",
			url:          "/bigbluebutton/api/getMeetings?name=doe",
			expectedBody: string(bytes),
		},
		{
			name:         "Invalid checksum returns 200 with checksum error",
			url:          "/bigbluebutton/api/getMeetings?name=doe&checksum=dummychecksum",
			expectedBody: string(bytes),
		},
		{
			name:         "Valid checksum goes through checksum validation middleware",
			url:          "/bigbluebutton/api/getMeetings?name=doe&checksum=80207d6781a83ac95b86d3c3884809fcfb8040fc",
			expectedBody: "/bigbluebutton/api/getMeetings",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := executeRequest(router, "GET", test.url, nil)
			assert.Equal(t, 200, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

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
				"Authorization": defaultAPIKey() + "dummy",
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
				"Authorization": defaultAPIKey(),
			},
			expectedStatus: 400,
			body:           nil,
		},
	}

	router := launchRouter(&config.Config{
		APIKey: defaultAPIKey(),
		RDB: config.RDB{
			Address:  container.URI,
			Password: "",
			DB:       0,
		},
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := executeRequestWithHeaders(router, "POST", "/admin/servers", test.body, test.headers)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
