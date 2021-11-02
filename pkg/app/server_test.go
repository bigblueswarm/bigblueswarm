package app

import (
	"b3lb/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoute(t *testing.T) {
	router := launchRouter(&config.Config{})
	w := executeRequest(router, "GET", "/bigbluebutton/api", nil)

	response := "<response><returncode>SUCCESS</returncode><version>2.0</version></response>"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, response, w.Body.String())
}
