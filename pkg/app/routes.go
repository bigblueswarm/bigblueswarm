package app

import "b3lb/pkg/admin"

func (s *Server) initRoutes() {
	router := s.Router

	base := router.Group("/bigbluebutton")
	{
		base.GET("", s.HealthCheck)
		api := base.Group("/api")
		{
			api.GET("", s.HealthCheck)
			api.Use(s.ChecksumValidation)
			api.GET("/create", s.Create)
			api.GET("/getMeetings", s.GetMeetings)
			api.GET("/join", s.Join)
		}
	}

	admin.CreateAdmin(s.InstanceManager, &s.Config.Admin).InitRoutes(router)
}
