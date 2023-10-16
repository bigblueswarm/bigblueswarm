// Package admin manages the bigblueswarm admin part
package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bigblueswarm/bigblueswarm/v3/pkg/balancer"
	"github.com/bigblueswarm/bigblueswarm/v3/pkg/config"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Admin struct manager bigblueswarm administration
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
		e := fmt.Errorf("failed to list instances: %s", err)
		log.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	c.JSON(http.StatusOK, instances)
}

// ClusterStatus send a status for the cluster. It contains all instances with their status
func (a *Admin) ClusterStatus(c *gin.Context) {
	instances, err := a.InstanceManager.List()
	if err != nil {
		e := fmt.Errorf("failed to retrieve instances: %s", err.Error())
		log.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	status, err := a.Balancer.ClusterStatus(instances)
	if err != nil {
		e := fmt.Errorf("failed to retrieve balancer cluster status: %s", err)
		log.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	c.JSON(http.StatusOK, status)
}

// SetInstances set all instances. It takes InstanceList object in body
func (a *Admin) SetInstances(c *gin.Context) {
	defer c.Request.Body.Close()

	instanceList := &InstanceList{}
	if err := c.ShouldBindJSON(instanceList); err != nil {
		e := fmt.Errorf("body does not bind InstanceList object: %s", err)
		log.Error(e)
		c.String(http.StatusBadRequest, e.Error())
		return
	}

	if err := a.InstanceManager.SetInstances(instanceList.Instances); err != nil {
		e := fmt.Errorf("failed to set instances in instance manager: %s", err)
		log.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}

// CreateTenant create a tenant from a configuraion YAML body
func (a *Admin) CreateTenant(c *gin.Context) {
	defer c.Request.Body.Close()

	tenant := &Tenant{}
	if err := c.ShouldBindJSON(tenant); err != nil {
		e := fmt.Errorf("body does not bind Tenant object: %s", err)
		log.Error(e)
		c.String(http.StatusBadRequest, e.Error())
		return
	}

	logger := log.WithField("tenant", tenant.Spec.Host)
	if tenant.Spec.Host == "" {
		m := "failed to create tenant. Tenant spec host should not be null"
		logger.Warn(m)
		c.String(http.StatusBadRequest, m)
		return
	}

	if err := a.TenantManager.AddTenant(tenant); err != nil {
		e := fmt.Errorf("failed to add tenant in tenant manager: %s", err)
		logger.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	logger.Info("tenant successfully created")
	c.AbortWithStatus(http.StatusCreated)
}

// ListTenants list all tenants in system
func (a *Admin) ListTenants(c *gin.Context) {
	tenants, err := a.TenantManager.ListTenants()
	if err != nil {
		e := fmt.Errorf("unable to list all tenants: %s", err)
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
		m := "hostname not found or empty"
		log.Warn(m)
		c.String(http.StatusBadRequest, m)
		return
	}

	logger := log.WithField("tenant", hostname)
	tenant, err := a.TenantManager.GetTenant(hostname)
	if err != nil {
		e := fmt.Errorf("failed to retrieve tenant: %s", err.Error())
		logger.Error(e)
		c.String(http.StatusInternalServerError, e.Error())
		return
	}

	if tenant == nil {
		m := "tenant not found for deletion"
		logger.Info(m)
		c.String(http.StatusNotFound, m)
		return
	}

	if err := a.TenantManager.DeleteTenant(hostname); err != nil {
		m := "unable to delete tenant"
		logger.Error(m, err)
		c.String(http.StatusInternalServerError, m)
	} else {
		logger.Info("tenant successfully deleted")
		c.AbortWithStatus(http.StatusNoContent)
	}
}

// GetConfiguration render configuration
func (a *Admin) GetConfiguration(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, a.Config)
}

// GetTenant retrieve a tenant based on its hostname
func (a *Admin) GetTenant(c *gin.Context) {
	hostname, exists := c.Params.Get("hostname")
	if !exists || strings.TrimSpace(hostname) == "" {
		m := "hostname not found or empty"
		log.Warn(m)
		c.String(http.StatusBadRequest, m)
		return
	}

	logger := log.WithField("tenant", hostname)
	tenant, err := a.TenantManager.GetTenant(hostname)
	if err != nil {
		m := "unable to retrieve tenant"
		logger.Error(m, err)
		c.String(http.StatusInternalServerError, m)
		return
	}

	if tenant == nil {
		logger.Info("tenant not found")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, tenant)
}
