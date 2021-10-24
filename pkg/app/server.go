package app

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{
		router: gin.Default(),
	}
}

func (s *Server) Run() error {

	s.Init_routes()
	err := s.router.Run()

	if err != nil {
		return err
	}

	return nil
}
