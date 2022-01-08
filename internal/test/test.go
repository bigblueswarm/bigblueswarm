package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

// Test is a test
type Test struct {
	Name      string
	Mock      func()
	Validator func(t *testing.T, value interface{}, err error)
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
			URL:    &url.URL{},
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

// SetRequestHeader set the request header
func SetRequestHeader(c *gin.Context, key, value string) {
	initRequestContext(c)
	c.Request.Header.Set(key, value)
}
