package test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// ExecuteRequest performs a request against the given router.
func ExecuteRequest(router *gin.Engine, method string, path string, body io.Reader) *httptest.ResponseRecorder {
	return ExecuteRequestWithHeaders(router, method, path, body, nil)
}

// ExecuteRequestWithHeaders performs a request against the given router with the given headers.
func ExecuteRequestWithHeaders(router *gin.Engine, method string, path string, body io.Reader, headers map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	router.ServeHTTP(w, req)

	return w
}
