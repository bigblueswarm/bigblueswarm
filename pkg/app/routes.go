package app

func (s *Server) InitRoutes() {
	router := s.router

	base := router.Group("/bigbluebutton")
	{
		base.GET("", s.HealthCheck)
		api := base.Group("/api")
		{
			api.GET("", s.HealthCheck)
			api.Use(s.ChecksumValidation)
			api.GET("/getMeetings", s.GetMeetings)
		}
	}
}
