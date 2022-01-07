package api

import (
	"b3lb/internal/mock"
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/restclient"
	"os"
	"testing"
)

var instance *BigBlueButtonInstance

func TestMain(m *testing.M) {
	instance = &BigBlueButtonInstance{
		URL:    "http://localhost:80/bigbluebutton",
		Secret: TestUtil.BBBSecret,
	}

	restclient.Client = &mock.MockRestClient{}

	status := m.Run()

	os.Exit(status)
}
