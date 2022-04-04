package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyValidation check that the request contains an api key provided by Authorization header
func (a *Admin) APIKeyValidation(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	auth = strings.TrimSpace(auth)
	if auth == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if auth != a.Config.Admin.APIKey {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
