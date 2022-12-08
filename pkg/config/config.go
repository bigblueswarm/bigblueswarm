// Package config manages the bigblueswarm config
package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// ConsulPrefix is the configuration path consul prefix
const ConsulPrefix string = "consul:"

// BigBlueButton configuration mapping
type BigBlueButton struct {
	Secret                 string `yaml:"secret" json:"secret"`
	RecordingsPollInterval string `yaml:"recordingsPollInterval" json:"recordingsPollInterval"`
}

// RDB represents redis database configuration mapping
type RDB struct {
	Address  string `yaml:"address" json:"address"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"database" json:"database"`
}

// IDB represents influxdb database configuration mapping
type IDB struct {
	Address      string `yaml:"address" json:"address"`
	Token        string `yaml:"token" json:"token"`
	Organization string `yaml:"organization" json:"organization"`
	Bucket       string `yaml:"bucket" json:"bucket"`
}

// AdminConfig represents the admin configuration
type AdminConfig struct {
	APIKey string `yaml:"apiKey" json:"apiKey"`
}

// BalancerConfig represents the balancer configuration
type BalancerConfig struct {
	MetricsRange string `yaml:"metricsRange" json:"metricsRange"`
	CPULimit     int    `yaml:"cpuLimit" json:"cpuLimit"`
	MemLimit     int    `yaml:"memLimit" json:"memLimit"`
}

// SetDefaultValues initialize BalancerConfig default values
func (bc *BalancerConfig) SetDefaultValues() {
	if bc.CPULimit == 0 {
		bc.CPULimit = 90
	}

	if bc.MemLimit == 0 {
		bc.MemLimit = 90
	}
}

// Port represents the BigBlueSwarm port configuration
type Port int

// Config represents main configuration mapping
type Config struct {
	BigBlueButton BigBlueButton  `yaml:"bigbluebutton" json:"bigbluebutton"`
	Admin         AdminConfig    `yaml:"admin" json:"admin"`
	Balancer      BalancerConfig `yaml:"balancer" json:"balancer"`
	Port          Port           `yaml:"port" json:"port"`
	RDB           RDB            `yaml:"redis" json:"redis"`
	IDB           IDB            `yaml:"influxdb" json:"influxdb"`
}

const defaultConfigFileName = "bigblueswarm.yaml"

// DefaultConfigFolder is the default config folder path
const DefaultConfigFolder = "$HOME/.bigblueswarm"

// DefaultConfigPath return the default config path file
func DefaultConfigPath() string {
	return fmt.Sprintf("%s/%s", DefaultConfigFolder, defaultConfigFileName)
}

// FormalizeConfigPath formalize config path. If config path is the default config path (home directory),
// it returns a computed path
func FormalizeConfigPath(path string) (string, error) {
	if path == DefaultConfigPath() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(homeDir, ".bigblueswarm", defaultConfigFileName)
	}

	return path, nil
}

// Load the configuration from the given path
func Load(path string) (*Config, error) {
	if IsConsulEnabled(path) {
		return LoadConfigFromConsul(path)
	}

	return loadConfigFromFile(path)
}

// Path return the flag config path
func Path() string {
	var configPath string

	flag.StringVar(&configPath, "config", DefaultConfigPath(), "Config file path")
	flag.Parse()

	return configPath
}
