package app

import (
	"b3lb/pkg/config"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	testcontainers.Container
	URI string
}

var container redisContainer

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Name:         "tc_redis",
		Image:        "redis:6.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redis, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	endpoint, endpointErr := redis.Endpoint(ctx, "")

	if endpointErr != nil {
		panic(err)
	}

	container = redisContainer{
		Container: redis,
		URI:       endpoint,
	}

	status := m.Run()

	terminationErr := redis.Terminate(ctx)
	if terminationErr != nil {
		panic(err)
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
