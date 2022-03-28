package admin

import (
	"fmt"
	"net/http"

	"github.com/SLedunois/b3lb/pkg/balancer"
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Admin struct manager b3lb administration
type Admin struct {
	InstanceManager InstanceManager
	TenantManager   TenantManager
	Balancer        balancer.Balancer
	Config          *config.AdminConfig
}

// CreateAdmin creates a new admin based on given configuration
func CreateAdmin(manager InstanceManager, tenantManager TenantManager, balancer balancer.Balancer, config *config.AdminConfig) *Admin {
	return &Admin{
		InstanceManager: manager,
		TenantManager:   tenantManager,
		Config:          config,
		Balancer:        balancer,
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

// SetInstances set all instances. It takes InstanceList object in body
func (a *Admin) SetInstances(c *gin.Context) {
	defer c.Request.Body.Close()

	instanceList := &InstanceList{}
	if err := c.ShouldBindYAML(instanceList); err != nil {
		e := fmt.Errorf("Body does not bind InstanceList object: %s", err)
		log.Error(e)
		c.String(http.StatusBadRequest, e.Error())
		return
	}

	if err := a.InstanceManager.SetInstances(instanceList.Instances); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}

// CreateTenant create a tenant from a configuraion YAML body
func (a *Admin) CreateTenant(c *gin.Context) {
	defer c.Request.Body.Close()

	tenant := &Tenant{}
	if err := c.ShouldBindYAML(tenant); err != nil {
		e := fmt.Errorf("Body does not bind Tenant object: %s", err)
		log.Error(e)
		c.String(http.StatusBadRequest, e.Error())
		return
	}

	if err := a.TenantManager.AddTenant(tenant); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}
