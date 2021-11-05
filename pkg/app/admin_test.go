package app

import (
	"b3lb/pkg/config"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddInstance(t *testing.T) {
	type test struct {
		name           string
		body           io.Reader
		expectedStatus int
		expectedBody   string
	}

	tests := []test{
		{
			name:           "Nil body should trigger a bad request status",
			body:           nil,
			expectedStatus: 400,
			expectedBody:   "",
		},
		{
			name:           "Correct body should create a bbb instance",
			body:           bytes.NewBufferString(`{"url": "http://localhost/bigbluebutton", "secret": "supersecret"}`),
			expectedStatus: 201,
			expectedBody:   `{"url":"http://localhost/bigbluebutton","secret":"supersecret"}`,
		},
		{
			name:           "Adding a duplication should returns a 409 conflict status",
			body:           bytes.NewBufferString(`{"url": "http://localhost/bigbluebutton", "secret": "supersecret"}`),
			expectedStatus: 409,
			expectedBody:   "",
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

	headers := map[string]string{
		"Authorization": defaultAPIKey(),
		"Content-Type":  "application/json",
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := executeRequestWithHeaders(router, "POST", "/admin/servers", test.body, headers)
			assert.Equal(t, test.expectedStatus, w.Code)

			if test.expectedBody != "" {
				assert.Equal(t, test.expectedBody, w.Body.String())
			}
		})
	}
}

func TestListInstances(t *testing.T) {
	router := launchRouter(&config.Config{
		APIKey: defaultAPIKey(),
		RDB: config.RDB{
			Address:  container.URI,
			Password: "",
			DB:       0,
		},
	})

	headers := map[string]string{
		"Authorization": defaultAPIKey(),
	}

	w := executeRequestWithHeaders(router, "GET", "/admin/servers", nil, headers)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `["http://localhost/bigbluebutton"]`, w.Body.String())
}

func TestDeleteInstance(t *testing.T) {

	type test struct {
		name           string
		url            string
		expectedStatus int
	}

	tests := []test{
		{
			name:           "Delete without an url should return a 400",
			url:            "/admin/servers",
			expectedStatus: 400,
		},
		{
			name:           "Delete a non existing instance should return a 404",
			url:            "/admin/servers?url=http://fakebbb",
			expectedStatus: 404,
		},
		{
			name:           "Delete an existing instance should return a 204",
			url:            "/admin/servers?url=http://localhost/bigbluebutton",
			expectedStatus: 204,
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

	headers := map[string]string{
		"Authorization": defaultAPIKey(),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := executeRequestWithHeaders(router, "DELETE", test.url, nil, headers)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
