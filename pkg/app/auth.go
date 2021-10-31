package app

import (
	"b3lb/pkg/api"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) ChecksumValidation(c *gin.Context) {
	error := &api.ChecksumError{
		Return_code: "FAILED",
		Message_key: "checksumError",
		Message:     "You did not pass the checksum security check",
	}

	checksum_param, exists := c.GetQuery("checksum")
	if !exists {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	params := c.Request.URL.Query()
	params.Del("checksum")

	checksum := &Checksum{
		Secret: s.config.BigBlueButton.Secret,
		Action: strings.TrimPrefix(c.FullPath(), "/bigbluebutton/api/"),
		Params: params,
	}

	fmt.Println("checksum to process", checksum.Value())
	sha := StringToSHA1(checksum.Value())
	fmt.Println("Checksum processed", string(sha))

	if checksum_param != string(sha) {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	c.Next()
}
