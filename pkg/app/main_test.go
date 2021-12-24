package app

import (
	"b3lb/pkg/config"
	"context"
	"fmt"
	"os"
	"testing"
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

func defaultConfig() *config.Config {
	return &config.Config{
		BigBlueButton: config.BigBlueButton{
			Secret: "secret",
		},
		APIKey: defaultAPIKey(),
		RDB: config.RDB{
			Address:  redisContainer.URI,
			Password: "",
			DB:       0,
		},
		IDB: config.IDB{
			Address:      fmt.Sprintf("http://%s", influxDBContainer.URI),
			Token:        influxDBToken,
			Bucket:       influxDBBucket,
			Organization: influxDBOrg,
		},
	}
}
