package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

// BigBlueButton configuration mapping
type BigBlueButton struct {
	Secret                 string `mapstructure:"secret"`
	RecordingsPollInterval string `mapstructure:"recordings_poll_interval"`
}

// RDB represents redis database configuration mapping
type RDB struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"database"`
}

// IDB represents influxdb database configuration mapping
type IDB struct {
	Address      string `mapstructure:"address"`
	Token        string `mapstructure:"token"`
	Organization string `mapstructure:"organization"`
	Bucket       string `mapstructure:"bucket"`
}

// AdminConfig represents the admin configuration
type AdminConfig struct {
	APIKey string `mapstructure:"api_key"`
	URL    string `mapstructure:"url"`
}

// BalancerConfig represents the balancer configuration
type BalancerConfig struct {
	MetricsRange string `mapstructure:"metrics_range"`
}

// Config represents main configuration mapping
type Config struct {
	BigBlueButton BigBlueButton  `mapstructure:"bigbluebutton"`
	Admin         AdminConfig    `mapstructure:"admin"`
	Balancer      BalancerConfig `mapstructure:"balancer"`
	Port          int            `mapstructure:"port"`
	RDB           RDB            `mapstructure:"redis"`
	IDB           IDB            `mapstructure:"influxdb"`
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

	if conf.Port == 0 {
		conf.Port = 8080
	}

	return conf, nil
}
