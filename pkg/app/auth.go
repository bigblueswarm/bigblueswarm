package app

import (
	"b3lb/pkg/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ChecksumValidation handler validate all requests checksum and returns an error if the checksum is not int the request or if the checksum is invalid
func (s *Server) ChecksumValidation(c *gin.Context) {
	error := api.DefaultChecksumError()

	checksumParam, exists := c.GetQuery("checksum")
	if !exists {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	params := c.Request.URL.Query()
	params.Del("checksum")

	checksum := &Checksum{
		Secret: s.Config.BigBlueButton.Secret,
		Action: strings.TrimPrefix(c.FullPath(), "/bigbluebutton/api/"),
		Params: params,
	}

	sha, err := StringToSHA1(checksum.Value())

	if err != nil {
		panic(err)
	}

	if checksumParam != string(sha) {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	c.Next()
}

// APIKeyValidation check that the request contains an api key provided by Authorization header
func (s *Server) APIKeyValidation(c *gin.Context) {

	auth := c.Request.Header.Get("Authorization")
	auth = strings.TrimSpace(auth)
	if auth == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if auth != s.Config.APIKey {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
