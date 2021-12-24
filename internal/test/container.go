package test

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// InfluxDBToken define the default influxdb token
const InfluxDBToken string = "Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg=="

// InfluxDBOrg define the default influxdb org
const InfluxDBOrg string = "b3lb"

// InfluxDBBucket define the default influxdb bucket
const InfluxDBBucket string = "bucket"

// BBBSecret define the default bigbluebutton secret
const BBBSecret = "0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU"

// Container represents a test container object
type Container struct {
	testcontainers.Container
	URI string
}

// InitRedisContainer creates a redis container
func InitRedisContainer(name string) *Container {
	req := testcontainers.ContainerRequest{
		Name:         name,
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

// InitInfluxDBContainer creates an influxdb container
func InitInfluxDBContainer(name string) *Container {
	req := testcontainers.ContainerRequest{
		Name:         name,
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
	createInfluxDBAccessToken(name)

	endpoint, endpointErr := influxdb.Endpoint(context.Background(), "")

	if endpointErr != nil {
		panic(err)
	}

	return &Container{
		Container: influxdb,
		URI:       endpoint,
	}
}

func createInfluxDBAccessToken(containerName string) {
	_, cmdErr := exec.Command("/bin/sh", "-c", fmt.Sprintf(`sudo docker exec %s sh -c "influx setup --name b3lbconfig --org %s --username admin --password password --token %s --bucket %s --retention 0 --force"`, containerName, InfluxDBOrg, InfluxDBToken, InfluxDBBucket)).Output()

	if cmdErr != nil {
		panic(cmdErr)
	}
}

// InitBigBlueButtonContainer creates a bigbluebutton container
func InitBigBlueButtonContainer(name string, port string) *Container {
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

	time.Sleep(90 * time.Second)

	endpoint, endpointErr := bbb.Endpoint(context.Background(), "")

	if endpointErr != nil {
		panic(err)
	}

	return &Container{
		Container: bbb,
		URI:       endpoint,
	}
}

// FormatBBBInstanceURL format bbb uri to a valid url
func FormatBBBInstanceURL(uri string) string {
	return fmt.Sprintf("http://%s/bigbluebutton", uri)
}

// SetBBBSecret set secret of Bigbluebutton containers
func SetBBBSecret(bbbs []*Container) {
	for _, bbb := range bbbs {
		name, err := bbb.Container.Name(context.Background())
		if err != nil {
			panic(err)
		}

		_, cmdErr := exec.Command("/bin/sh", "-c", fmt.Sprintf(`sudo docker exec %s sh -c "bbb-conf --setsecret %s"`, name, BBBSecret)).Output()
		if cmdErr != nil {
			panic(err)
		}
	}
}

// WriteIDBData write some custom metrics for both bbb containers in influxdb
func WriteIDBData(cluster *Cluster) {
	client := influxdb2.NewClient(fmt.Sprintf("http://%s", cluster.InfluxDB.URI), InfluxDBToken)
	writeAPI := client.WriteAPI(InfluxDBOrg, InfluxDBBucket)
	bbb1 := cluster.BigBlueButtons[0]
	bbb2 := cluster.BigBlueButtons[1]

	// Write some custom cpu usage points
	writeAPI.WritePoint(idbCPUUsagePoint(fmt.Sprintf("%s_cpu_1", "bbb1"), bbb1.URI, 10, FormatBBBInstanceURL(bbb1.URI)))
	writeAPI.WritePoint(idbCPUUsagePoint(fmt.Sprintf("%s_cpu_2", "bbb2"), bbb2.URI, 80, FormatBBBInstanceURL(bbb2.URI)))

	// Wite some custom mem used_percent points
	writeAPI.WritePoint(idbMemPoint(fmt.Sprintf("%s_mem_1", "bbb1"), bbb1.URI, 25, FormatBBBInstanceURL(bbb1.URI)))
	writeAPI.WritePoint(idbMemPoint(fmt.Sprintf("%s_mem_2", "bbb2"), bbb2.URI, 30, FormatBBBInstanceURL(bbb2.URI)))

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
