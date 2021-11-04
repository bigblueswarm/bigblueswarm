package app

// InitRoutes init server routes
func (s *Server) InitRoutes() {
	router := s.Router

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
