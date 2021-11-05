package config

import (
	"testing"

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
