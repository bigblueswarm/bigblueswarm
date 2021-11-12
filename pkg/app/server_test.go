package app

import (
	"b3lb/pkg/config"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var redisContainer Container
var influxDBContainer Container
var bbbContainers []*Container

func TestMain(m *testing.M) {
	ctx := context.Background()

	redisContainer = *initRedisContainer()
	influxDBContainer = *initInfluxDBContainer()
	bbbContainers = []*Container{
		initBigBlueButtonContainer("tc_bbb1", "80/tcp"),
		initBigBlueButtonContainer("tc_bbb2", "80/tcp"),
	}

	writeIDBData(fmt.Sprintf("http://%s", influxDBContainer.URI), bbbContainers)
	insertBBBInstances(redisContainer, bbbContainers)
	setBBBSecret(bbbContainers)

	status := m.Run()

	rTErr := redisContainer.Container.Terminate(ctx)
	if rTErr != nil {
		panic(rTErr)
	}

	iDBTErr := influxDBContainer.Container.Terminate(ctx)
	if iDBTErr != nil {
		panic(iDBTErr)
	}

	for _, bbb := range bbbContainers {
		bbbErr := bbb.Container.Terminate(ctx)
		if bbbErr != nil {
			panic(bbbErr)
		}
	}

	os.Exit(status)
}

func TestHealthCheckRoute(t *testing.T) {
	router := launchRouter(&config.Config{})
	w := executeRequest(router, "GET", "/bigbluebutton/api", nil)

	response := "<response><returncode>SUCCESS</returncode><version>2.0</version></response>"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, response, w.Body.String())
}
