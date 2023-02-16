// Package app is the bigblueswarm core
package app

import (
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/balancer"
	log "github.com/sirupsen/logrus"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/restclient"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// SentryEnabled tells if sentry is enabled or not. If it is, we add a gin hook for performance monitoring
var SentryEnabled = false

// Server struct represents an object containings the server router and its configuration
type Server struct {
	Router          *gin.Engine
	Config          *config.Config
	InstanceManager admin.InstanceManager
	TenantManager   admin.TenantManager
	Mapper          Mapper
	Balancer        balancer.Balancer
}

// NewServer creates a new server based on given configuration
func NewServer(config *config.Config) *Server {
	redisClient := utils.RedisClient(config)
	influxClient := utils.InfluxDBClient(config)

	restclient.Init()

	router := gin.Default()
	if SentryEnabled {
		log.Info("Sentry enabled: adding gin middleware")
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	return &Server{
		Router:          router,
		Config:          config,
		InstanceManager: admin.NewInstanceManager(*redisClient),
		TenantManager:   admin.NewTenantManager(*redisClient),
		Mapper:          NewMapper(*redisClient),
		Balancer:        balancer.New(influxClient, &config.Balancer, &config.IDB),
	}
}

// Run launches the server
func (s *Server) Run() error {
	s.initRoutes()
	go s.launchRecordingPoller()
	err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port))

	if err != nil {
		return err
	}

	return nil
}
