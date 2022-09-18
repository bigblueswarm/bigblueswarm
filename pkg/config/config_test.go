package config

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/b3lb/test_utils/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestFSConfigLoad(t *testing.T) {

	type test struct {
		name  string
		path  string
		check func(t *testing.T, config *Config, err error)
	}

	tests := []test{
		{
			name: "Configuration loading does not returns any error with a valid path",
			path: "../../config.yml",
			check: func(t *testing.T, config *Config, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, config)
			},
		},
		{
			name: "Configuration loading returns an error with an invalid path",
			path: "config.yml",
			check: func(t *testing.T, config *Config, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, config)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := Load(test.path)
			test.check(t, config, err)
		})
	}
}

func TestConsulConfigLoad(t *testing.T) {
	var url string
	var bbbConf string
	var adminConf string
	var balancerConf string
	var portConf string
	var rdbConf string
	var idbConf string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.ReplaceAll(r.RequestURI, "/v1/kv/configuration/", "")

		switch key {
		case "bigbluebutton":
			w.Write([]byte(bbbConf))
		case "admin":
			w.Write([]byte(adminConf))
		case "balancer":
			w.Write([]byte(balancerConf))
		case "port":
			w.Write([]byte(portConf))
		case "redis":
			w.Write([]byte(rdbConf))
		case "influxdb":
			w.Write([]byte(idbConf))
		}
	}))

	defer server.Close()

	tests := []test.Test{
		{
			Name: "an invalid url should return an error",
			Mock: func() {
				url = "invalid_url:333333"
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "an error while loading admin configuration should return an error",
			Mock: func() {
				url = server.URL
				bbbConf = `[{"LockIndex":0,"Key":"configuration/bigbluebutton","Flags":0,"Value":"c2VjcmV0OiAwb2w1dDQ0VVIyMXJyUDB4TDVvdTdJQkZ1bVdGM0dFTmViZ1cxUnlUZmJVCnJlY29yZGluZ3NQb2xsSW50ZXJ2YWw6IDFt","CreateIndex":35,"ModifyIndex":1211}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "an error while loading balancer configuration should return an error",
			Mock: func() {
				url = server.URL
				adminConf = `[{"LockIndex":0,"Key":"configuration/admin","Flags":0,"Value":"YXBpS2V5OiBrZ3BxclRpcE0yeWpjWHd6NXBPeEJLVmlFOW9OWDc2Ug==","CreateIndex":40,"ModifyIndex":1219}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "an error while loading port configuration should return an error",
			Mock: func() {
				url = server.URL
				balancerConf = `[{"LockIndex":0,"Key":"configuration/balancer","Flags":0,"Value":"bWV0cmljc1JhbmdlOiAtNW0KY3B1TGltaXQ6IDk5Cm1lbUxpbWl0OiA5OQ==","CreateIndex":38,"ModifyIndex":1214}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "an error while loading redis configuration should return an error",
			Mock: func() {
				url = server.URL
				portConf = `[{"LockIndex":0,"Key":"configuration/port","Flags":0,"Value":"ODA5MA==","CreateIndex":42,"ModifyIndex":42}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "an error while loading influx configuration should return an error",
			Mock: func() {
				url = server.URL
				rdbConf = `[{"LockIndex":0,"Key":"configuration/redis","Flags":0,"Value":"YWRkcmVzczogbG9jYWxob3N0OjYzNzkKcGFzc3dvcmQ6CmRhdGFiYXNlOiAw","CreateIndex":46,"ModifyIndex":46}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			Name: "no error should return a valid configuration",
			Mock: func() {
				url = server.URL
				idbConf = `[{"LockIndex":0,"Key":"configuration/influxdb","Flags":0,"Value":"YWRkcmVzczogaHR0cDovL2xvY2FsaG9zdDo4MDg2CnRva2VuOiBacTl3THNtaG5XNVV0T2lQSkFwVXYxY1RWSmZ3WHNUZ2xfcENraVRpa1EzZzJZR1B0UzVIcXNYZWYtV2Y1cFVVM3dqWTNuVldUWVJJLVdjOExqYkRmZz09Cm9yZ2FuaXphdGlvbjogYjNsYgpidWNrZXQ6IGJ1Y2tldA==","CreateIndex":50,"ModifyIndex":50}]`
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				conf := value.(*Config)
				assert.Nil(t, err)
				expected := &Config{
					Admin: AdminConfig{
						APIKey: "kgpqrTipM2yjcXwz5pOxBKViE9oNX76R",
					},
					Balancer: BalancerConfig{
						MetricsRange: "-5m",
						CPULimit:     99,
						MemLimit:     99,
					},
					BigBlueButton: BigBlueButton{
						Secret:                 "0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU",
						RecordingsPollInterval: "1m",
					},
					Port: 8090,
					IDB: IDB{
						Address:      "http://localhost:8086",
						Token:        "Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg==",
						Organization: "b3lb",
						Bucket:       "bucket",
					},
					RDB: RDB{
						Address:  "localhost:6379",
						Password: "",
						DB:       0,
					},
				}

				assert.Equal(t, expected, conf)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			config, err := Load(fmt.Sprintf("%s%s", ConsulPrefix, url))
			test.Validator(t, config, err)
		})
	}
}

func TestDefaultConfigPath(t *testing.T) {
	assert.Equal(t, "$HOME/.b3lb/b3lb.yaml", DefaultConfigPath())
}

func TestFormalizeConfigPath(t *testing.T) {
	type test struct {
		name     string
		path     string
		expected string
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []test{
		{
			name:     "a custom path should return the custom path",
			path:     "/etc/config/b3lb.yaml",
			expected: "/etc/config/b3lb.yaml",
		},
		{
			name:     "default path should return the home path",
			path:     DefaultConfigPath(),
			expected: fmt.Sprintf("%s/.b3lb/b3lb.yaml", homeDir),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := FormalizeConfigPath(test.path)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, path)
		})
	}
}

func TestBalancerConfigSetDefaultValues(t *testing.T) {
	config := &BalancerConfig{}
	tests := []test.Test{
		{
			Name: "no values for cpu and mem should set 100 on cpu and mem configuration",
			Mock: func() {},
			Validator: func(t *testing.T, value interface{}, err error) {
				conf := value.(*BalancerConfig)
				assert.Equal(t, 100, conf.CPULimit)
				assert.Equal(t, 100, conf.MemLimit)
			},
		},
		{
			Name: "custom values for cpu and mem should not override values",
			Mock: func() {
				config.CPULimit = 30
				config.MemLimit = 30
			},
			Validator: func(t *testing.T, value interface{}, err error) {
				conf := value.(*BalancerConfig)
				assert.Equal(t, 30, conf.CPULimit)
				assert.Equal(t, 30, conf.MemLimit)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			config.SetDefaultValues()
			test.Validator(t, config, nil)
		})
	}
}

func TestConfigPath(t *testing.T) {
	tests := []test.Test{
		{
			Name: "passing no -config flag should return default config path",
			Mock: func() {},
			Validator: func(t *testing.T, value interface{}, err error) {
				path := value.(string)
				assert.Equal(t, DefaultConfigPath(), path)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Mock()
			test.Validator(t, Path(), nil)
		})
	}
}
