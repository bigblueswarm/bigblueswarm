package admin

import (
	"net/http"

	"github.com/SLedunois/b3lb/pkg/api"

	"github.com/gin-gonic/gin"
)

// InitRoutes initialize the admin routes
func (a *Admin) InitRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.Use(a.APIKeyValidation)
		servers := admin.Group("/servers")
		servers.Handle("GET", "", a.ListInstances)
		servers.Handle("POST", "", a.AddInstance)
		servers.Handle("DELETE", "", a.DeleteInstance)
	}
}

// Routes returns admin routes
func (a *Admin) Routes() *[]api.EndpointGroup {
	return &[]api.EndpointGroup{
		{
			Path: "/admin",
			Endpoints: []interface{}{
				api.Endpoint{
					Handler: a.APIKeyValidation,
				},
				api.EndpointGroup{
					Path: "/servers",
					Endpoints: []interface{}{
						api.Endpoint{
							Method:  http.MethodGet,
							Handler: a.ListInstances,
						},
						api.Endpoint{
							Method:  http.MethodPost,
							Handler: a.AddInstance,
						},
						api.Endpoint{
							Method:  http.MethodDelete,
							Handler: a.DeleteInstance,
						},
					},
				},
			},
		},
	}
}
