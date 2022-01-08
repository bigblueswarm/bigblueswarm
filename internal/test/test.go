package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Test is a test
type Test struct {
	Name      string
	Mock      func()
	Validator func(value interface{}, err error) bool
}

// StringToJSONArray convert a json string to a string array
func StringToJSONArray(value string) []string {
	var data []string
	_ = json.Unmarshal([]byte(value), &data)
	return data
}

func initRequestContext(c *gin.Context) {
	if c.Request == nil {
		c.Request = &http.Request{
			Header: make(http.Header),
		}
	}
}

// AddRequestBody create a request body from a json string
func AddRequestBody(c *gin.Context, data string) {
	initRequestContext(c)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(data)))
}

// SetRequestMethod set the request method
func SetRequestMethod(c *gin.Context, method string) {
	initRequestContext(c)
	c.Request.Method = method
}

// SetRequestContentType set the request content type
func SetRequestContentType(c *gin.Context, contentType string) {
	initRequestContext(c)
	c.Request.Header.Set("Content-Type", contentType)
}

// SetRequestParams set the request url
func SetRequestParams(c *gin.Context, value string) {
	initRequestContext(c)
	c.Request.URL = &url.URL{
		RawQuery: value,
	}
}
