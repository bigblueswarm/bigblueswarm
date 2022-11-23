package utils

import "github.com/gin-gonic/gin"

func GetHost(ctx *gin.Context) string {
	if host := ctx.Request.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}

	return ctx.Request.Host
}
