package app

func (s *Server) Init_routes() {
	router := s.router

	base := router.Group("/bigbluebutton")
	{
		base.GET("", s.healthCheck)
		api := base.Group("/api")
		{
			api.GET("", s.healthCheck)
			api.Use(s.checksumValidation)
			api.GET("/getMeetings", s.getMeetings)
		}
	}
}
