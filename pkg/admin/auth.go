// Package admin manages the bigblueswarm admin part
package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// APIKeyValidation check that the request contains an api key provided by Authorization header
func (a *Admin) APIKeyValidation(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	auth = strings.TrimSpace(auth)
	logger := log.WithField("auth", auth)
	if auth == "" {
		logger.Warn("auth key can't be an empty string")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if auth != a.Config.Admin.APIKey {
		logger.Error("auth key does not match the configured admin key")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
