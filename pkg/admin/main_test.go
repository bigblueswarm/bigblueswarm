package admin

import (
	"os"
	"testing"

	TestUtil "github.com/SLedunois/b3lb/internal/test"
	"github.com/SLedunois/b3lb/pkg/balancer"
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

var router *gin.Engine
var instanceManager InstanceManager
var tenantManager TenantManager
var redisMock redismock.ClientMock
var redisClient *redis.Client

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	client, mock := redismock.NewClientMock()
	redisClient = client
	redisMock = mock

	instanceManager = NewInstanceManager(*client)
	tenantManager = NewTenantManager(*client)

	router = gin.Default()
	CreateAdmin(instanceManager, tenantManager, &balancer.Mock{}, &config.AdminConfig{
		APIKey: TestUtil.DefaultAPIKey(),
	})

	status := m.Run()
	if err := redisMock.ExpectationsWereMet(); err != nil {
		panic(err)
	}

	os.Exit(status)
}
