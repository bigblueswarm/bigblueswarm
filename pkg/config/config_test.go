package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {

	t.Run("Configuration loading does not returns any error with a valid path", func(t *testing.T) {
		conf, err := Load("../../config.yml")

		assert.Nil(t, err)
		assert.NotNil(t, conf)
	})

	t.Run("Configuration loading returns an error with an invalid path", func(t *testing.T) {
		conf, err := Load("config.yml")

		assert.Nil(t, conf)
		assert.NotNil(t, err)
	})
}
