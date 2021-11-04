package app

import (
	"b3lb/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returning an health check response
func (s *Server) HealthCheck(c *gin.Context) {
	status := &api.HealtCheck{
		ReturnCode: "SUCCESS",
		Version:    "2.0",
	}

	c.XML(http.StatusOK, status)
}

// Handler returning the getMeetings API. See https://docs.bigbluebutton.org/dev/api.html#getmeetings.
func (s *Server) GetMeetings(c *gin.Context) {
	c.String(http.StatusOK, c.FullPath())
}
