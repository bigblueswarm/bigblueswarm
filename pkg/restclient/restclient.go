package restclient

import (
	"bytes"
	"net/http"
)

var (
	// Client is the http client used to make requests
	Client HTTPClient
)

// Init initializes the restclient package.
func Init() {
	Client = &http.Client{}
}

// HTTPClient defines an HTTP client interface to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Get performs an HTTP GET request.
func Get(url string) (*http.Response, error) {
	return perform(http.MethodGet, url, map[string]string{}, nil)
}

func perform(method string, url string, headers map[string]string, body []byte) (*http.Response, error) {
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	for k, v := range headers {
		request.Header.Add(k, v)
	}

	return Client.Do(request)
}

// GetWithHeaders performs an HTTP GET request with headers.
func GetWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	return perform(http.MethodGet, url, headers, nil)
}

// Post performs an HTTP POST request.
func Post(url string, body []byte) (*http.Response, error) {
	return perform(http.MethodPost, url, map[string]string{}, body)
}

// PostWithHeaders performs an HTTP POST request with headers.
func PostWithHeaders(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return perform(http.MethodPost, url, headers, body)
}

// Delete performs an HTTP DELETE request.
func Delete(url string) (*http.Response, error) {
	return perform(http.MethodDelete, url, map[string]string{}, nil)
}

// DeleteWithHeaders performs an HTTP DELETE request with headers.
func DeleteWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	return perform(http.MethodDelete, url, headers, nil)
}
