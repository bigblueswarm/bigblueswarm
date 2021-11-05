package app

import (
	"b3lb/pkg/admin"

	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AddInstance insert the body into the database.
func (s *Server) AddInstance(c *gin.Context) {
	instance := &admin.BigBlueButtonInstance{}
	if err := c.ShouldBind(&instance); err != nil || (instance.Secret == "" || instance.URL == "") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	exists, err := s.Manager.Exists(*instance)

	if err != nil {
		log.Error("Failed to check if instance already exists", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if exists {
		log.Warn("Instance already exists")
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	if err := s.Manager.Add(*instance); err != nil {
		log.Error("Failed to add new instance", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusCreated, instance)
	}
}

// ListInstances returns Bigbluebutton instance list
func (s *Server) ListInstances(c *gin.Context) {
	instances, err := s.Manager.List()
	if err != nil {
		log.Error("Failed to list instances", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, instances)
}

// DeleteInstance deletes an instance
func (s *Server) DeleteInstance(c *gin.Context) {
	if URL, ok := c.GetQuery("url"); ok {
		exists, err := s.Manager.Exists(admin.BigBlueButtonInstance{URL: URL})

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := s.Manager.Remove(URL); err != nil {
			log.Error("Failed to delete instance", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusNoContent)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
