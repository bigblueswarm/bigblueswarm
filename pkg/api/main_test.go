package api

import (
	"os"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/restclient"
	"github.com/bigblueswarm/test_utils/pkg/test"
)

var instance *BigBlueButtonInstance

func TestMain(m *testing.M) {
	instance = &BigBlueButtonInstance{
		URL:    "http://localhost:80/bigbluebutton",
		Secret: test.DefaultSecret(),
	}

	restclient.Client = &restclient.Mock{}

	status := m.Run()

	os.Exit(status)
}
