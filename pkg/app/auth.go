package app

import (
	"b3lb/pkg/api"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func string_to_sha1(value string) string {
	hasher := sha1.New()
	hasher.Write([]byte("getmeetings"))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *Server) checksumValidation(c *gin.Context) {
	error := &api.ChecksumError{
		Return_code: "FAILED",
		Message_key: "checksumError",
		Message:     "You did not pass the checksum security check",
	}

	checksum, exists := c.GetQuery("checksum")
	if !exists {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	sha := string_to_sha1("getmeetings")

	fmt.Println(string(sha))
	fmt.Println(checksum)

	if checksum != string(sha) {
		c.XML(http.StatusOK, error)
		c.Abort()
		return
	}

	c.Next()
}
