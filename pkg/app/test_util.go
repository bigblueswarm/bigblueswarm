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
	server.initRoutes()
	return server.Router
}

func executeRequest(router *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	return executeRequestWithHeaders(router, method, path, body, nil)
}

func executeRequestWithHeaders(router *gin.Engine, method string, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	router.ServeHTTP(w, req)

	return w
}

func defaultAPIKey() string {
	return "supersecret"
}

/*func defaultSecret() string {
	return "supersecret"
}*/
