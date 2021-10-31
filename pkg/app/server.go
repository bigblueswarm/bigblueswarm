package app

import (
	"b3lb/pkg/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		router: gin.Default(),
		config: config,
	}
}

func (s *Server) Run() error {

	s.InitRoutes()
	err := s.router.Run(":8090")

	if err != nil {
		return err
	}

	return nil
}
