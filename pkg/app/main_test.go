package app

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/admin"
	"b3lb/pkg/api"
	"b3lb/pkg/config"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var cluster *TestUtil.Cluster

func TestMain(m *testing.M) {
	ctx := context.Background()

	cluster = &TestUtil.Cluster{
		Redis:    TestUtil.InitRedisContainer("tc_redis_app"),
		InfluxDB: TestUtil.InitInfluxDBContainer("tc_influxdb_app"),
		BigBlueButtons: []*TestUtil.Container{
			TestUtil.InitBigBlueButtonContainer("tc_bbb1", "80/tcp"),
			TestUtil.InitBigBlueButtonContainer("tc_bbb2", "80/tcp"),
		},
	}

	TestUtil.WriteIDBData(cluster)
	insertBBBInstances(cluster)
	TestUtil.SetBBBSecret(cluster.BigBlueButtons)

	status := m.Run()

	rTErr := cluster.Redis.Container.Terminate(ctx)
	if rTErr != nil {
		panic(rTErr)
	}

	iDBTErr := cluster.InfluxDB.Container.Terminate(ctx)
	if iDBTErr != nil {
		panic(iDBTErr)
	}

	for _, bbb := range cluster.BigBlueButtons {
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
		Admin: config.AdminConfig{
			APIKey: TestUtil.DefaultAPIKey(),
		},
		RDB: config.RDB{
			Address:  cluster.Redis.URI,
			Password: "",
			DB:       0,
		},
		IDB: config.IDB{
			Address:      fmt.Sprintf("http://%s", cluster.InfluxDB.URI),
			Token:        TestUtil.InfluxDBToken,
			Bucket:       TestUtil.InfluxDBBucket,
			Organization: TestUtil.InfluxDBOrg,
		},
	}
}

func launchRouter(config *config.Config) *gin.Engine {
	server := NewServer(config)
	server.initRoutes()
	return server.Router
}

// InsertBBBInstances inserts bigbluebutton instances into the database
func insertBBBInstances(cluster *TestUtil.Cluster) {
	client := redis.NewClient(&redis.Options{
		Addr:     cluster.Redis.URI,
		Password: "",
		DB:       0,
	})

	manager := admin.NewInstanceManager(client)

	for _, bbb := range cluster.BigBlueButtons {
		instance := api.BigBlueButtonInstance{
			URL:    TestUtil.FormatBBBInstanceURL(bbb.URI),
			Secret: TestUtil.BBBSecret,
		}

		if err := manager.Add(instance); err != nil {
			panic(err)
		}
	}
}
