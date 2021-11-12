package app

import (
	"b3lb/pkg/admin"
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const influxDBToken string = "Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg=="
const influxDBOrg string = "b3lb"
const influxDBBucket string = "bucket"

const bbbSecret = "0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU"

func launchRouter(config *config.Config) *gin.Engine {
	server := NewServer(config)
	server.initRoutes()
	return server.Router
}

func executeRequest(router *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	return executeRequestWithHeaders(router, method, path, body, nil)
}

func executeRequestWithHeaders(router *gin.Engine, method string, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	router.ServeHTTP(w, req)

	return w
}

func defaultAPIKey() string {
	return "supersecret"
}

// Container represents a test container object
type Container struct {
	testcontainers.Container
	URI string
}

func initRedisContainer() *Container {
	req := testcontainers.ContainerRequest{
		Name:         "tc_redis",
		Image:        "redis:6.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
		AutoRemove:   true,
	}
	redis, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	endpoint, endpointErr := redis.Endpoint(context.Background(), "")

	if endpointErr != nil {
		panic(err)
	}

	return &Container{
		Container: redis,
		URI:       endpoint,
	}
}

func initInfluxDBContainer() *Container {
	req := testcontainers.ContainerRequest{
		Name:         "tc_influxdb",
		Image:        "influxdb:2.0",
		ExposedPorts: []string{"8086/tcp"},
		AutoRemove:   true,
	}

	influxdb, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	time.Sleep(15 * time.Second)
	createInfluxDBAccessToken()

	endpoint, endpointErr := influxdb.Endpoint(context.Background(), "")

	if endpointErr != nil {
		panic(err)
	}

	return &Container{
		Container: influxdb,
		URI:       endpoint,
	}
}

func createInfluxDBAccessToken() {
	_, cmdErr := exec.Command("/bin/sh", "-c", fmt.Sprintf(`sudo docker exec tc_influxdb sh -c "influx setup --name b3lbconfig --org %s --username admin --password password --token %s --bucket %s --retention 0 --force"`, influxDBOrg, influxDBToken, influxDBBucket)).Output()

	if cmdErr != nil {
		panic(cmdErr)
	}
}

func initBigBlueButtonContainer(name string, port string) *Container {
	req := testcontainers.ContainerRequest{
		Name:         name,
		Image:        "sledunois/bbb-dev:2.4-develop",
		ExposedPorts: []string{port},
		AutoRemove:   true,
		Privileged:   true,
	}

	bbb, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Minute)

	endpoint, endpointErr := bbb.Endpoint(context.Background(), "")

	if endpointErr != nil {
		panic(err)
	}

	return &Container{
		Container: bbb,
		URI:       endpoint,
	}
}

func insertBBBInstances(redisContainer Container, bbbs []*Container) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisContainer.URI,
		Password: "",
		DB:       0,
	})

	manager := admin.NewInstanceManager(client)

	for _, bbb := range bbbs {
		instance := api.BigBlueButtonInstance{
			URL:    formatBBBInstanceURL(bbb.URI),
			Secret: bbbSecret,
		}

		if err := manager.Add(instance); err != nil {
			panic(err)
		}
	}
}

func formatBBBInstanceURL(uri string) string {
	return fmt.Sprintf("http://%s/bigbluebutton", uri)
}

func setBBBSecret(bbbs []*Container) {
	for _, bbb := range bbbs {
		name, err := bbb.Container.Name(context.Background())
		if err != nil {
			panic(err)
		}

		_, cmdErr := exec.Command("/bin/sh", "-c", fmt.Sprintf(`sudo docker exec %s sh -c "bbb-conf --setsecret %s"`, name, bbbSecret)).Output()
		if cmdErr != nil {
			panic(err)
		}
	}
}

func writeIDBData(address string, bbbs []*Container) {
	client := influxdb2.NewClient(address, influxDBToken)
	writeAPI := client.WriteAPI(influxDBOrg, influxDBBucket)
	bbb1 := bbbs[0]
	bbb2 := bbbs[1]

	// Write some custom cpu usage points
	writeAPI.WritePoint(idbCPUUsagePoint(fmt.Sprintf("%s_cpu_1", "bbb1"), bbb1.URI, 10, formatBBBInstanceURL(bbb1.URI)))
	writeAPI.WritePoint(idbCPUUsagePoint(fmt.Sprintf("%s_cpu_2", "bbb2"), bbb2.URI, 80, formatBBBInstanceURL(bbb2.URI)))

	// Wite some custom mem used_percent points
	writeAPI.WritePoint(idbMemPoint(fmt.Sprintf("%s_mem_1", "bbb1"), bbb1.URI, 25, formatBBBInstanceURL(bbb1.URI)))
	writeAPI.WritePoint(idbMemPoint(fmt.Sprintf("%s_mem_2", "bbb2"), bbb2.URI, 30, formatBBBInstanceURL(bbb2.URI)))

	writeAPI.Flush()
	client.Close()
}

func idbCPUUsagePoint(id string, host string, value int, b3lbHost string) *write.Point {
	return influxdb2.NewPoint(
		"cpu",
		map[string]string{
			"id":        id,
			"hostname":  host,
			"b3lb_host": b3lbHost,
			"cpu":       "cpu-total",
		},
		map[string]interface{}{
			"usage_system": value,
		},
		time.Now())
}

func idbMemPoint(id string, host string, value int, b3lbHost string) *write.Point {
	return influxdb2.NewPoint(
		"mem",
		map[string]string{
			"id":        id,
			"hostname":  host,
			"b3lb_host": b3lbHost,
		},
		map[string]interface{}{
			"used_percent": value,
		},
		time.Now())
}
