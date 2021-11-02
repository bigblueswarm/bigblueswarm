package app

import (
	"b3lb/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	Config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Router: gin.Default(),
		Config: config,
	}
}

func (s *Server) Run() error {

	s.InitRoutes()
	err := s.Router.Run(":8090")

	if err != nil {
		return err
	}

	return nil
}
