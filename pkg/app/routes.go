package app

func (s *Server) initRoutes() {
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
	admin := router.Group("/admin")
	{
		admin.Use(s.APIKeyValidation)
		servers := admin.Group("/servers")
		servers.Handle("GET", "", s.ListInstances)
		servers.Handle("POST", "", s.AddInstance)
		servers.Handle("DELETE", "", s.DeleteInstance)
	}
}
