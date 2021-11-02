package app

import (
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecksumValidation(t *testing.T) {
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

	t.Run("No checksum returns 200 with checksum error", func(t *testing.T) {
		w := executeRequest(router, "GET", "/bigbluebutton/api/getMeetings?name=doe", nil)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, w.Body.String(), string(bytes))
	})

	t.Run("Invalid checksum returns 200 with checksum error", func(t *testing.T) {
		w := executeRequest(router, "GET", "/bigbluebutton/api/getMeetings?name=doe&checksum=dummychecksum", nil)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, w.Body.String(), string(bytes))
	})

	t.Run("Valid checksum goes through checksum validation middleware", func(t *testing.T) {
		w := executeRequest(router, "GET", "/bigbluebutton/api/getMeetings?name=doe&checksum=80207d6781a83ac95b86d3c3884809fcfb8040fc", nil)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "/bigbluebutton/api/getMeetings", w.Body.String())
	})
}
