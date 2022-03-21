package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var conf *Config

// BigBlueButton configuration mapping
type BigBlueButton struct {
	Secret                 string `yaml:"secret"`
	RecordingsPollInterval string `yaml:"recordingsPollInterval"`
}

// RDB represents redis database configuration mapping
type RDB struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"database"`
}

// IDB represents influxdb database configuration mapping
type IDB struct {
	Address      string `yaml:"address"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}

// AdminConfig represents the admin configuration
type AdminConfig struct {
	APIKey string `yaml:"apiKey"`
	URL    string `yaml:"url"`
}

// BalancerConfig represents the balancer configuration
type BalancerConfig struct {
	MetricsRange string `yaml:"metricsRange"`
	CPULimit     int    `yaml:"cpuLimit"`
	MemLimit     int    `yaml:"memLimit"`
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

// Port represents the B3LB port configuration
type Port int

// Config represents main configuration mapping
type Config struct {
	BigBlueButton BigBlueButton  `yaml:"bigbluebutton"`
	Admin         AdminConfig    `yaml:"admin"`
	Balancer      BalancerConfig `yaml:"balancer"`
	Port          Port           `yaml:"port"`
	RDB           RDB            `yaml:"redis"`
	IDB           IDB            `yaml:"influxdb"`
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error(fmt.Sprintf("unable to close config file: %s", err))
		}
	}()

	conf = &Config{}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(b, &conf)

	if conf.Port == 0 {
		conf.Port = 8080
	}

	conf.Balancer.SetDefaultValues()

	return conf, nil
}
