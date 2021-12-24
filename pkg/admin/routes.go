package admin

import "github.com/gin-gonic/gin"

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
