package api

import (
	TestUtil "b3lb/internal/test"
	"os"
	"testing"
	"time"
)

var instance *BigBlueButtonInstance

func TestMain(m *testing.M) {
	bbb := TestUtil.InitBigBlueButtonContainer("tc_bbb1_api", "80/tcp")
	time.Sleep(30 * time.Second)
	TestUtil.SetBBBSecret([]*TestUtil.Container{bbb})

	instance = &BigBlueButtonInstance{
		URL:    TestUtil.FormatBBBInstanceURL(bbb.URI),
		Secret: TestUtil.BBBSecret,
	}

	status := m.Run()

	os.Exit(status)
}
