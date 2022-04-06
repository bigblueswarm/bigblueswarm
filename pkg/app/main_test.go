package app

import (
	"os"
	"testing"

	"github.com/SLedunois/b3lb/v2/pkg/admin"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var (
	redisMock       redismock.ClientMock
	redisClient     *redis.Client
	mapper          Mapper
	instanceManager admin.InstanceManager
)

func TestMain(m *testing.M) {
	//ctx := context.Background()
	gin.SetMode(gin.TestMode)
	client, rMock := redismock.NewClientMock()
	redisClient = client
	redisMock = rMock
	mapper = NewMapper(*redisClient)
	instanceManager = admin.NewInstanceManager(*redisClient)

	status := m.Run()

	if err := redisMock.ExpectationsWereMet(); err != nil {
		panic(err)
	}

	os.Exit(status)
}
