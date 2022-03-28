package restclient

import "net/http"

// Mock is a mock implementation of the restclient.Client interface.
type Mock struct{}

var (
	// RestClientMockDoFunc is the function that will be called when the mock rest client is used.
	RestClientMockDoFunc func(req *http.Request) (*http.Response, error)
)

// Do is a mock implementation of the restclient.Client.Do function.
func (m *Mock) Do(req *http.Request) (*http.Response, error) {
	return RestClientMockDoFunc(req)
}
