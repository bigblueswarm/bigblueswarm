package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SLedunois/b3lb/v2/pkg/balancer"
	"github.com/SLedunois/b3lb/v2/pkg/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Admin struct manager b3lb administration
type Admin struct {
	InstanceManager InstanceManager
	TenantManager   TenantManager
	Balancer        balancer.Balancer
	Config          *config.Config
}

// CreateAdmin creates a new admin based on given configuration
func CreateAdmin(manager InstanceManager, tenantManager TenantManager, balancer balancer.Balancer, config *config.Config) *Admin {
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
	if err := c.ShouldBindJSON(instanceList); err != nil {
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
	if err := c.ShouldBindJSON(tenant); err != nil {
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

// ListTenants list all tenants in system
func (a *Admin) ListTenants(c *gin.Context) {
	tenants, err := a.TenantManager.ListTenants()
	if err != nil {
		e := fmt.Errorf("Unable to list all tenants: %s", err)
		log.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	list := &TenantList{
		Kind:    "TenantList",
		Tenants: tenants,
	}

	c.JSON(http.StatusOK, list)
}

// DeleteTenant delete a given tenant
func (a *Admin) DeleteTenant(c *gin.Context) {
	hostname, exists := c.Params.Get("hostname")
	if !exists || strings.TrimSpace(hostname) == "" {
		c.String(http.StatusBadRequest, "hostname not found or empty")
		return
	}

	tenant, err := a.TenantManager.GetTenant(hostname)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve tenant %s: %s", hostname, err.Error()))
		return
	}

	if tenant == nil {
		c.String(http.StatusNotFound, fmt.Sprintf("Tenant %s not found for deletion", hostname))
		return
	}

	if err := a.TenantManager.DeleteTenant(hostname); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to delete tenant: %s", err.Error()))
	} else {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

// GetConfiguration render configuration
func (a *Admin) GetConfiguration(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, a.Config)
}

// GetTenant retrieve a tenant based on its hostname
func (a *Admin) GetTenant(c *gin.Context) {
	hostname, _ := c.Params.Get("hostname")
	tenant, err := a.TenantManager.GetTenant(hostname)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to retrieve tenant %s: %s", hostname, err.Error()))
		return
	}

	if tenant == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, tenant)
}
