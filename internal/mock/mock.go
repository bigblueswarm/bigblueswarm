package mock

import "net/http"

// MockRestClient is a mock implementation of the restclient.Client interface.
type MockRestClient struct {
	// DoFunc is the function that will be called when the Do function is called.
	DoFunc func(req *http.Request) (*http.Response, error)
}

// DoFunc is the function that will be called when the mock rest client is used.
var (
	DoFunc func(req *http.Request) (*http.Response, error)
)

// Do is a mock implementation of the restclient.Client.Do function.
func (m *MockRestClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
