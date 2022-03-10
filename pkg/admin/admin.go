package admin

import (
	"net/http"

	"github.com/SLedunois/b3lb/pkg/api"
	"github.com/SLedunois/b3lb/pkg/balancer"
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Admin struct manager b3lb administration
type Admin struct {
	InstanceManager InstanceManager
	Balancer        balancer.Balancer
	Config          *config.AdminConfig
}

// CreateAdmin creates a new admin based on given configuration
func CreateAdmin(manager InstanceManager, balancer balancer.Balancer, config *config.AdminConfig) *Admin {
	return &Admin{
		InstanceManager: manager,
		Config:          config,
		Balancer:        balancer,
	}
}

// AddInstance insert the body into the database.
func (a *Admin) AddInstance(c *gin.Context) {
	instance := &api.BigBlueButtonInstance{}
	if err := c.ShouldBind(&instance); err != nil || (instance.Secret == "" || instance.URL == "") {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	exists, err := a.InstanceManager.Exists(*instance)

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

	if err := a.InstanceManager.Add(*instance); err != nil {
		log.Error("Failed to add new instance", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		c.JSON(http.StatusCreated, instance)
	}
}

// ListInstances returns Bigbluebutton instance list
func (a *Admin) ListInstances(c *gin.Context) {
	instances, err := a.InstanceManager.ListInstances()
	if err != nil {
		log.Error("Failed to list instances", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, instances)
}

// DeleteInstance deletes an instance
func (a *Admin) DeleteInstance(c *gin.Context) {
	if URL, ok := c.GetQuery("url"); ok {
		exists, err := a.InstanceManager.Exists(api.BigBlueButtonInstance{URL: URL})

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := a.InstanceManager.Remove(URL); err != nil {
			log.Error("Failed to delete instance", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

// ClusterStatus send a status for the cluster. It contains all instances with their status
func (a *Admin) ClusterStatus(c *gin.Context) {
	instances, err := a.InstanceManager.List()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	status, err := a.Balancer.ClusterStatus(instances)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, status)
}
