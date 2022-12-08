// Package utils provide few utilies functions
package utils

import "github.com/gin-gonic/gin"

// GetHost get the gin request host. It returns the X-Forwarded-Host header if present however it returns the request host
func GetHost(ctx *gin.Context) string {
	if host := ctx.Request.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}

	return ctx.Request.Host
}
