package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type BigBlueButton struct {
	Secret string `mapstructure:"secret"`
}

type Config struct {
	BigBlueButton BigBlueButton `mapstructure:"bigbluebutton"`
}

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
