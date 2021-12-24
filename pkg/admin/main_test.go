package admin

import (
	TestUtil "b3lb/internal/test"
	"b3lb/pkg/config"
	"b3lb/pkg/utils"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	cluster := &TestUtil.Cluster{
		Redis: TestUtil.InitRedisContainer("tc_redis_admin"),
	}

	time.Sleep(20 * time.Second)

	instanceManager := NewInstanceManager(utils.RedisClient(&config.Config{
		RDB: config.RDB{
			Address:  cluster.Redis.URI,
			DB:       0,
			Password: "",
		},
	}))

	router = gin.Default()
	CreateAdmin(instanceManager, &config.AdminConfig{
		APIKey: TestUtil.DefaultAPIKey(),
	}).InitRoutes(router)

	status := m.Run()

	rTErr := cluster.Redis.Container.Terminate(ctx)
	if rTErr != nil {
		panic(rTErr)
	}

	os.Exit(status)
}
