// Package app is the bigblueswarm core
package app

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/api"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func processParameters(query string) string {
	reg := regexp.MustCompile(`&?checksum(=[^&]*&?)`)
	return reg.ReplaceAllString(query, "")
}

// ChecksumValidation handler validate all requests checksum and returns an error if the checksum is not int the request or if the checksum is invalid
func (s *Server) ChecksumValidation(c *gin.Context) {
	error := api.DefaultChecksumError()

	checksumParam, exists := c.GetQuery("checksum")
	if !exists {
		log.Warn("checksum not found in request")
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	logger := log.WithFields(log.Fields{
		"checksum": checksumParam,
		"tenant":   utils.GetHost(c),
	})

	secret := s.Config.BigBlueButton.Secret
	tenant, err := s.TenantManager.GetTenant(utils.GetHost(c))
	if err != nil {
		logger.Error("tenant manager can't retrieve tenant: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if tenant == nil {
		logger.Info("tenant not found")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if tenant.Spec.Secret != "" {
		secret = tenant.Spec.Secret
	}

	checksum := &api.Checksum{
		Secret: secret,
		Action: strings.TrimPrefix(c.FullPath(), "/bigbluebutton/api/"),
		Params: processParameters(c.Request.URL.RawQuery),
	}

	sha, err := checksum.Process()

	if err != nil {
		logger.Error("failed to process checksum validation", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if checksumParam != string(sha) {
		logger.WithFields(log.Fields{
			"tenant":         utils.GetHost(c),
			"checksum":       checksumParam,
			"expected_value": string(sha),
			"params":         checksum.Params,
		}).Warn("checksum does not pass the checksum validation")
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	setAPIContext(c, checksum)

	c.Next()
}

func setAPIContext(c *gin.Context, checksum *api.Checksum) {
	c.Set("api_ctx", checksum)
}

func getAPIContext(c *gin.Context) *api.Checksum {
	return c.MustGet("api_ctx").(*api.Checksum)
}
