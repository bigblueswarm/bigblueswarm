package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

// BigBlueButton configuration mapping
type BigBlueButton struct {
	Secret string `mapstructure:"secret"`
}

// Config represents main configuration mapping
type Config struct {
	BigBlueButton BigBlueButton `mapstructure:"bigbluebutton"`
}

// Load the configuration from the given path
func Load(path string) (*Config, error) {
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(path)

	if err != nil {
		return nil, err
	}

	conf := &Config{}

	if err := config.BindStruct("", &conf); err != nil {
		return nil, err
	}

	return conf, nil
}
