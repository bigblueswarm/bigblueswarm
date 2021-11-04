package app

import (
	"b3lb/pkg/config"

	"github.com/gin-gonic/gin"
)

// Server struct represents an object containings the server router and its configuration
type Server struct {
	Router *gin.Engine
	Config *config.Config
}

// Create a new server based on given configuration
func NewServer(config *config.Config) *Server {
	return &Server{
		Router: gin.Default(),
		Config: config,
	}
}

// Launch the server
func (s *Server) Run() error {

	s.InitRoutes()
	err := s.Router.Run(":8090")

	if err != nil {
		return err
	}

	return nil
}
