package config

import (
	"flag"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type Config struct {
	BigBlueButton struct {
		Secret string `mapstructure: "secret"`
	} `mapstructure:"bigbluebutton"`
}

var config_path string

func Init() (*Config, error) {
	flag.StringVar(&config_path, "config", "config.yml", "Config file path")
	flag.Parse()

	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(config_path)

	if err != nil {
		return nil, err
	}

	var conf = &Config{}
	config.BindStruct("", &conf)

	return conf, nil
}
