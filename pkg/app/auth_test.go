package app

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"encoding/xml"
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
			w := TestUtil.ExecuteRequest(router, "GET", test.url, nil)
			assert.Equal(t, 200, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
