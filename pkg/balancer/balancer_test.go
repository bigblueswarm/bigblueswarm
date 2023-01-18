package balancer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"
	"github.com/bigblueswarm/test_utils/pkg/test"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New(nil, nil, nil))
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

	defer server.Close()

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
,result,table,_time,bigblueswarm_host,_value
,balancer,0,2022-01-10T16:57:30Z,http://localhost:8080,`
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
		IDBConfig: &config.IDB{
			Bucket: "bucket",
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

func TestClusterStatus(t *testing.T) {
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
			Name: "No result returned by influxDB should return an empty array",
			Mock: func() {
				statusCode = http.StatusOK
			},
			Validator: func(t *testing.T, result interface{}, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 0, len(result.([]InstanceStatus)))
			},
		},
		{
			Name: "A valid result returned by influxDB should return the result",
			Mock: func() {
				statusCode = http.StatusOK
				body = `#group,false,false,true,true,true,false
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,string,double
#default,cpu,,,,,
,result,table,_start,_stop,bigblueswarm_host,_value
,,0,2022-03-10T13:41:21.808246343Z,2022-03-10T13:46:21.808246343Z,http://localhost/bigbluebutton,8.840855095953396

#group,false,false,true,true,true,false
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,string,double
#default,mem,,,,,
,result,table,_start,_stop,bigblueswarm_host,_value
,,0,2022-03-10T13:41:21.808246343Z,2022-03-10T13:46:21.808246343Z,http://localhost/bigbluebutton,68.82938000133251

#group,false,false,true,true,true,false,false,false
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,string,long,long,long
#default,bbb,,,,,,,
,result,table,_start,_stop,bigblueswarm_host,online,meetings,participants
,,0,2022-03-10T13:41:21.808246343Z,2022-03-10T13:46:21.808246343Z,http://localhost/bigbluebutton,1,0,0`
			},
			Validator: func(t *testing.T, result interface{}, err error) {
				assert.Nil(t, err)
				status := result.([]InstanceStatus)[0]
				assert.Equal(t, "http://localhost/bigbluebutton", status.Host)
				assert.Equal(t, float64(8.84), status.CPU)
				assert.Equal(t, float64(68.83), status.Mem)
				assert.Equal(t, "Up", status.APIStatus)
				assert.Equal(t, int64(0), status.Meetings)
				assert.Equal(t, int64(0), status.Participants)
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
		IDBConfig: &config.IDB{
			Bucket: "bucket",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			status, err := balancer.ClusterStatus([]string{"http://localhost:8080"})
			test.Validator(t, status, err)
		})
	}
}

func TestApiStatusToString(t *testing.T) {
	assert.Equal(t, "Up", apiStatusToString(int64(1)))
	assert.Equal(t, "Down", apiStatusToString(int64(0)))
}

func TestGetCurrentState(t *testing.T) {
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

	defer server.Close()

	tests := []test.Test{
		{
			Name: "an error returned by influxdb should be returned",
			Mock: func() {
				statusCode = http.StatusInternalServerError
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Error(t, err)
			},
		},
		{
			Name: "a valid result should be parsed and returned",
			Mock: func() {
				statusCode = http.StatusOK
				body = `#group,false,false,true,true,false,false,true,true
#datatype,string,long,dateTime:RFC3339,dateTime:RFC3339,dateTime:RFC3339,long,string,string
#default,_result,,,,,,,
,result,table,_start,_stop,_time,_value,_field,_measurement
,,0,2022-12-13T18:32:35.99098962Z,2022-12-13T18:38:35.99098962Z,2022-12-13T18:38:35.99098962Z,218,meetings,bigbluebutton:localhost:8090`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.Equal(t, int64(218), value.(int64))
				assert.Nil(t, err)
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
			AggregationInterval: "10s",
		},
		IDBConfig: &config.IDB{
			Bucket: "bucket",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			state, err := balancer.GetCurrentState("bigbluebutton:localhost:8090", "meetings")
			test.Validator(t, state, err)
		})
	}
}
