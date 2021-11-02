package app

import (
	"b3lb/pkg/config"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func launchRouter(config *config.Config) *gin.Engine {
	server := NewServer(config)
	server.InitRoutes()

	return server.Router
}

func executeRequest(router *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	router.ServeHTTP(w, req)

	return w
}
