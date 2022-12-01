package utils

import (
	"testing"

	"github.com/bigblueswarm/test_utils/pkg/request"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetHost(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)

	t.Run("No X-Forwarded-Host host should return the host request", func(t *testing.T) {
		host := "localhost"
		request.SetRequestHost(c, host)
		assert.Equal(t, host, GetHost(c))
	})

	t.Run("A X-Forwarded-Host header should be returned instead of request host", func(t *testing.T) {
		headerValue := "mydummyheaderhostname"
		request.SetRequestHeader(c, "X-Forwarded-Host", headerValue)
		assert.Equal(t, headerValue, GetHost(c))
	})
}
