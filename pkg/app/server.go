// Package app is the bigblueswarm core
package app

import (
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/admin"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/balancer"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/restclient"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"

	"github.com/gin-gonic/gin"
)

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

	return &Server{
		Router:          gin.Default(),
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
