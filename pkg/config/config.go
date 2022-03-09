package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	CPULimit     int    `mapstructure:"cpu_limit"`
	MemLimit     int    `mapstructure:"mem_limit"`
}

// SetDefaultValues initialize BalancerConfig default values
func (bc *BalancerConfig) SetDefaultValues() {
	if bc.CPULimit == 0 {
		bc.CPULimit = 100
	}

	if bc.MemLimit == 0 {
		bc.MemLimit = 100
	}
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

const defaultConfigFileName = ".b3lb.yaml"

// DefaultConfigPath return the default config path file
func DefaultConfigPath() string {
	return fmt.Sprintf("$HOME/%s", defaultConfigFileName)
}

// FormalizeConfigPath formalize config path. If config path is the default config path (home directory),
// it returns a computed path
func FormalizeConfigPath(path string) (string, error) {
	if path == DefaultConfigPath() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(homeDir, defaultConfigFileName)
	}

	return path, nil
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

	conf.Balancer.SetDefaultValues()

	return conf, nil
}
