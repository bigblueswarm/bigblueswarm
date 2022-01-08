package api

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/restclient"
	"b3lb/pkg/restclient/mock"
	"os"
	"testing"
)

var instance *BigBlueButtonInstance

func TestMain(m *testing.M) {
	instance = &BigBlueButtonInstance{
		URL:    "http://localhost:80/bigbluebutton",
		Secret: TestUtil.BBBSecret,
	}

	restclient.Client = &mock.RestClient{}

	status := m.Run()

	os.Exit(status)
}
