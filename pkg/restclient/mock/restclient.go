package mock

import "net/http"

// RestClient is a mock implementation of the restclient.Client interface.
type RestClient struct{}

// DoFunc is the function that will be called when the mock rest client is used.
var (
	DoFunc func(req *http.Request) (*http.Response, error)
)

// Do is a mock implementation of the restclient.Client.Do function.
func (m *RestClient) Do(req *http.Request) (*http.Response, error) {
	return DoFunc(req)
}
