package admin

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/config"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var router *gin.Engine
var instanceManager InstanceManager
var redisMock redismock.ClientMock
var redisClient *redis.Client

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	client, mock := redismock.NewClientMock()
	redisClient = client
	redisMock = mock

	instanceManager = NewInstanceManager(client)

	router = gin.Default()
	CreateAdmin(instanceManager, &config.AdminConfig{
		APIKey: TestUtil.DefaultAPIKey(),
	}).InitRoutes(router)

	status := m.Run()
	if err := redisMock.ExpectationsWereMet(); err != nil {
		panic(err)
	}

	os.Exit(status)
}
