package app

import (
	"os"
	"testing"

	"github.com/SLedunois/b3lb/pkg/admin"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var (
	redisMock       redismock.ClientMock
	redisClient     *redis.Client
	sessionManager  SessionManager
	instanceManager admin.InstanceManager
)

func TestMain(m *testing.M) {
	//ctx := context.Background()
	gin.SetMode(gin.TestMode)
	client, rMock := redismock.NewClientMock()
	redisClient = client
	redisMock = rMock
	sessionManager = NewSessionManager(*redisClient)
	instanceManager = admin.NewInstanceManager(*redisClient)

	/*cluster = &TestUtil.Cluster{
		InfluxDB: TestUtil.InitInfluxDBContainer("tc_influxdb_app"),
		BigBlueButtons: []*TestUtil.Container{
			TestUtil.InitBigBlueButtonContainer("tc_bbb1", "80/tcp"),
			TestUtil.InitBigBlueButtonContainer("tc_bbb2", "80/tcp"),
		},
	}

	TestUtil.WriteIDBData(cluster)
	insertBBBInstances(cluster)
	TestUtil.SetBBBSecret(cluster.BigBlueButtons)*/

	status := m.Run()

	/*iDBTErr := cluster.InfluxDB.Container.Terminate(ctx)
	if iDBTErr != nil {
		panic(iDBTErr)
	}

	for _, bbb := range cluster.BigBlueButtons {
		bbbErr := bbb.Container.Terminate(ctx)
		if bbbErr != nil {
			panic(bbbErr)
		}
	}*/

	if err := redisMock.ExpectationsWereMet(); err != nil {
		panic(err)
	}

	os.Exit(status)
}
