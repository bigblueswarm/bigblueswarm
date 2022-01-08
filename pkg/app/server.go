package app

import (
	"b3lb/pkg/admin"
	"b3lb/pkg/config"
	"b3lb/pkg/restclient"
	"b3lb/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Server struct represents an object containings the server router and its configuration
type Server struct {
	Router          *gin.Engine
	Config          *config.Config
	InstanceManager admin.InstanceManager
	SessionManager  *SessionManager
	Balancer        *Balancer
}

// NewServer creates a new server based on given configuration
func NewServer(config *config.Config) *Server {
	redisClient := utils.RedisClient(config)
	influxClient := utils.InfluxDBClient(config)

	restclient.Init()

	return &Server{
		Router:          gin.Default(),
		Config:          config,
		InstanceManager: admin.NewInstanceManager(redisClient),
		SessionManager:  NewSessionManager(*redisClient),
		Balancer:        NewBalancer(influxClient),
	}
}

// Run launches the server
func (s *Server) Run() error {
	s.initRoutes()
	err := s.Router.Run(fmt.Sprintf(":%d", s.Config.Port))

	if err != nil {
		return err
	}

	return nil
}
