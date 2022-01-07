package restclient

import "net/http"

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
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	return Client.Do(request)
}
