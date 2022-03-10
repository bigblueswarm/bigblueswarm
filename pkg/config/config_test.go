package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/SLedunois/b3lb/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {

	type test struct {
		name     string
		filename string
		check    func(t *testing.T, config *Config, err error)
	}

	tests := []test{
		{
			name:     "Configuration loading does not returns any error with a valid path",
			filename: "../../config.yml",
			check: func(t *testing.T, config *Config, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, config)
			},
		},
		{
			name:     "Configuration loading returns an error with an invalid path",
			filename: "config.yml",
			check: func(t *testing.T, config *Config, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, config)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := Load(test.filename)
			test.check(t, config, err)
		})
	}
}

func TestDefaultConfigPath(t *testing.T) {
	assert.Equal(t, "$HOME/.b3lb.yaml", DefaultConfigPath())
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
			expected: fmt.Sprintf("%s/.b3lb.yaml", homeDir),
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
