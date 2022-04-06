package api

import (
	"os"
	"testing"

	TestUtil "github.com/SLedunois/b3lb/v2/internal/test"
	"github.com/SLedunois/b3lb/v2/pkg/restclient"
)

var instance *BigBlueButtonInstance

func TestMain(m *testing.M) {
	instance = &BigBlueButtonInstance{
		URL:    "http://localhost:80/bigbluebutton",
		Secret: TestUtil.DefaultSecret(),
	}

	restclient.Client = &restclient.Mock{}

	status := m.Run()

	os.Exit(status)
}
