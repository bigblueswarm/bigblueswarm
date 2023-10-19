// Package config manages the bigblueswarm config
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConsulPrefix is the configuration path consul prefix
const ConsulPrefix string = "consul:"

// BigBlueSwarm configuration mapping
type BigBlueSwarm struct {
	Secret                 string `yaml:"secret" json:"secret"`
	RecordingsPollInterval string `yaml:"recordingsPollInterval" json:"recordingsPollInterval"`
}

// RDB represents redis database configuration mapping
type RDB struct {
	Address  string `yaml:"address" json:"address"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"database" json:"database"`
}

// PG represents postgresql configuration mapping
type PG struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
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
	MetricsRange        string `yaml:"metricsRange" json:"metricsRange"`
	CPULimit            int    `yaml:"cpuLimit" json:"cpuLimit"`
	MemLimit            int    `yaml:"memLimit" json:"memLimit"`
	AggregationInterval string `yaml:"aggregationInterval" json:"aggregationInterval"`
}

// SetDefaultValues initialize BalancerConfig default values
func (bc *BalancerConfig) SetDefaultValues() {
	if bc.CPULimit == 0 {
		bc.CPULimit = 90
	}

	if bc.MemLimit == 0 {
		bc.MemLimit = 90
	}

	if bc.AggregationInterval == "" {
		bc.AggregationInterval = "10s"
	}
}

// SetDefaultValues initialize BigBlueSwarm default values
func (bbs *BigBlueSwarm) SetDefaultValues() {
	if bbs.RecordingsPollInterval == "" {
		bbs.RecordingsPollInterval = "15m"
	}
}

// Port represents the BigBlueSwarm port configuration
type Port int

// Config represents main configuration mapping
type Config struct {
	BigBlueSwarm BigBlueSwarm   `yaml:"bigblueswarm" json:"bigblueswarm"`
	Admin        AdminConfig    `yaml:"admin" json:"admin"`
	Balancer     BalancerConfig `yaml:"balancer" json:"balancer"`
	Port         Port           `yaml:"port" json:"port"`
	RDB          RDB            `yaml:"redis" json:"redis"`
	IDB          IDB            `yaml:"influxdb" json:"influxdb"`
	PG           PG             `yaml:"postgres" json:"postgres"`
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
