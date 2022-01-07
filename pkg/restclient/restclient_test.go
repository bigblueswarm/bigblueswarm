package restclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("Performing GET request should returns the requestURI and a 200 http status code", func(t *testing.T) {
		url := fmt.Sprintf("%s/test_get", server.URL)

		response, err := Get(url)
		if err != nil {
			t.Error(err)
			return
		}

		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, response.StatusCode, http.StatusOK)
		assert.Equal(t, string(bytes), "/test_get")
	})
}
