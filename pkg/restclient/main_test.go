package restclient

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server *httptest.Server
)

func TestMain(m *testing.M) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write([]byte(r.RequestURI))
		if err != nil {
			panic(err)
		}
	}))

	Init()

	status := m.Run()

	server.Close()
	os.Exit(status)
}
