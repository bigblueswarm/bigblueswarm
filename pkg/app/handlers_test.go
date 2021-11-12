package app

import (
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	type test struct {
		name              string
		url               string
		expectedKey       string
		expectedMessage   string
		expectedCode      string
		expectedSessionID string
	}

	router := launchRouter(&config.Config{
		BigBlueButton: config.BigBlueButton{
			Secret: "secret",
		},
		RDB: config.RDB{
			Address:  redisContainer.URI,
			Password: "",
			DB:       0,
		},
		IDB: config.IDB{
			Address:      fmt.Sprintf("http://%s", influxDBContainer.URI),
			Token:        influxDBToken,
			Bucket:       influxDBBucket,
			Organization: influxDBOrg,
		},
	})

	tests := []test{
		{
			name:              "Create with no name should returns a `no name error`",
			url:               "/bigbluebutton/api/create?checksum=025401fa251bdbcbba29c347a9cf811f29aa15a1",
			expectedCode:      api.FailedReturnCode,
			expectedKey:       api.ValidationErrorMessageKey,
			expectedMessage:   api.EmptyMeetingNameMessage,
			expectedSessionID: "",
		},
		{
			name:              "Create with no meeting id should returns a `no meeting id error`",
			url:               "/bigbluebutton/api/create?name=doe&checksum=291411768c7aeb243819983459755b32b96aae34",
			expectedCode:      api.FailedReturnCode,
			expectedKey:       api.ValidationErrorMessageKey,
			expectedMessage:   api.EmptyMeetingIDMessage,
			expectedSessionID: "",
		},
		{
			name:              "Valid create call should create a meeting",
			url:               "/bigbluebutton/api/create?name=doe&meetingID=id&checksum=62b3012fb88e9c468782e021ea20cded09cec0b4",
			expectedCode:      api.SuccessReturnCode,
			expectedKey:       "",
			expectedMessage:   "",
			expectedSessionID: "id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := executeRequest(router, "GET", test.url, nil)
			var response api.CreateResponse
			if err := xml.Unmarshal(w.Body.Bytes(), &response); err != nil {
				panic(err)
			}

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, test.expectedCode, response.ReturnCode)
			assert.Equal(t, test.expectedKey, response.MessageKey)
			assert.Equal(t, test.expectedMessage, response.Message)
			assert.Equal(t, test.expectedSessionID, response.MeetingID)
		})
	}
}
