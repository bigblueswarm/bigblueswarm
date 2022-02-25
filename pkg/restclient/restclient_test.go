package restclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/stretchr/testify/assert"
)

type HttpTest struct {
	test.Test
	Validator      func(t *testing.T, req *http.Request)
	PerformRequest func(s *httptest.Server)
}

func TestGet(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing GET request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				Get(fmt.Sprintf("%s/test_get", s.URL))
			},
		},
		{
			Test: test.Test{
				Name: "Performing GET request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_get", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				Get(fmt.Sprintf("%s/test_get", s.URL))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}

func TestGetWithHeaders(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing GetWithHeaders request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				GetWithHeaders(fmt.Sprintf("%s/test_get_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing GetWithHeaders request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_get_with_headers", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				GetWithHeaders(fmt.Sprintf("%s/test_get_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing GetWithHeaders request with no headers should not set any headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, 2, len(req.Header))
			},
			PerformRequest: func(s *httptest.Server) {
				GetWithHeaders(fmt.Sprintf("%s/test_get_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing GetWithHeaders request with headers should use the correct headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "test_value", req.Header.Get("test_header"))
			},
			PerformRequest: func(s *httptest.Server) {
				GetWithHeaders(fmt.Sprintf("%s/test_get_with_headers", s.URL), map[string]string{
					"test_header": "test_value",
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}

func TestPost(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing POST request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				Post(fmt.Sprintf("%s/test_post", s.URL))
			},
		},
		{
			Test: test.Test{
				Name: "Performing POST request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_post", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				Post(fmt.Sprintf("%s/test_post", s.URL))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}

func TestPostWithHeaders(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing PostWithHeaders request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodPost, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				PostWithHeaders(fmt.Sprintf("%s/test_get_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing PostWithHeaders request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_post_with_headers", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				PostWithHeaders(fmt.Sprintf("%s/test_post_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing PostWithHeaders request with no headers should not set any headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, 3, len(req.Header))
			},
			PerformRequest: func(s *httptest.Server) {
				PostWithHeaders(fmt.Sprintf("%s/test_post_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing PostWithHeaders request with headers should use the correct headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "test_value", req.Header.Get("test_header"))
			},
			PerformRequest: func(s *httptest.Server) {
				PostWithHeaders(fmt.Sprintf("%s/test_post_with_headers", s.URL), map[string]string{
					"test_header": "test_value",
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}

func TestDelete(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing DELETE request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodDelete, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				Delete(fmt.Sprintf("%s/test_delete", s.URL))
			},
		},
		{
			Test: test.Test{
				Name: "Performing DELETE request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_delete", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				Delete(fmt.Sprintf("%s/test_delete", s.URL))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}

func TestDeleteWithHeaders(t *testing.T) {
	tests := []HttpTest{
		{
			Test: test.Test{
				Name: "Performing DeleteWithHeaders request should use the correct method",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, http.MethodDelete, req.Method)
			},
			PerformRequest: func(s *httptest.Server) {
				DeleteWithHeaders(fmt.Sprintf("%s/test_delete_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing DeleteWithHeaders request should use the correct request uri",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "/test_delete_with_headers", req.RequestURI)
			},
			PerformRequest: func(s *httptest.Server) {
				DeleteWithHeaders(fmt.Sprintf("%s/test_delete_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing DeleteWithHeaders request with no headers should not set any headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, 2, len(req.Header))
			},
			PerformRequest: func(s *httptest.Server) {
				DeleteWithHeaders(fmt.Sprintf("%s/test_delete_with_headers", s.URL), map[string]string{})
			},
		},
		{
			Test: test.Test{
				Name: "Performing DeleteWithHeaders request with headers should use the correct headers",
			},
			Validator: func(t *testing.T, req *http.Request) {
				assert.Equal(t, "test_value", req.Header.Get("test_header"))
			},
			PerformRequest: func(s *httptest.Server) {
				DeleteWithHeaders(fmt.Sprintf("%s/test_delete_with_headers", s.URL), map[string]string{
					"test_header": "test_value",
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.Validator(t, r)
			}))

			test.PerformRequest(server)
		})

	}
}
