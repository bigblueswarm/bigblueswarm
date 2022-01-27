package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SLedunois/b3lb/pkg/utils"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestFormatInstanceFilter(t *testing.T) {
	type Test struct {
		name      string
		parameter []string
		expected  string
	}
	tests := []Test{
		{
			name:      "No parameter should returns empty string",
			parameter: []string{},
			expected:  "",
		},
		{
			name:      "One parameter should returns a valid filter",
			parameter: []string{"http://localhost:8080"},
			expected:  `r["b3lb_host"] == "http://localhost:8080"`,
		},
		{
			name:      "Multiple parameters should returns a valid filter",
			parameter: []string{"http://localhost:8080", "http://localhost:8081"},
			expected:  `r["b3lb_host"] == "http://localhost:8080" or r["b3lb_host"] == "http://localhost:8081"`,
		},
	}

	balancer := &InfluxDBBalancer{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, balancer.formatInstancesFilter(test.parameter))
		})
	}
}

func TestProcess(t *testing.T) {
	var (
		statusCode int
		body       string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		} else {
			_, err := w.Write([]byte(body))
			if err != nil {
				panic(err)
			}
		}
	}))

	tests := []test.Test{
		{
			Name: "An error thrown by influxDB should return an error",
			Mock: func() {
				statusCode = http.StatusInternalServerError
			},
			Validator: func(t *testing.T, result interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "No result returned by influxDB should return an empty string",
			Mock: func() {
				statusCode = http.StatusOK
			},
			Validator: func(t *testing.T, result interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "", result)
			},
		},
		{
			Name: "A valid result returned by influxDB should return the result",
			Mock: func() {
				statusCode = http.StatusOK
				body = `#datatype,string,long,dateTime:RFC3339,string,double
#group,false,false,false,false,false
#default,_result,,,,
,result,table,_time,b3lb_host,_value
,,0,2022-01-10T16:57:30Z,http://localhost:8080,`
			},
			Validator: func(t *testing.T, result interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "http://localhost:8080", result)
			},
		},
	}

	balancer := &InfluxDBBalancer{
		Client: utils.InfluxDBClient(&config.Config{
			IDB: config.IDB{
				Address: server.URL,
			},
		}),
		Config: &config.BalancerConfig{
			MetricsRange: "-5m",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			instance, err := balancer.Process([]string{"http://localhost:8080", "http://localhost:8081"})
			test.Validator(t, instance, err)
		})
	}
}
