// Package config manages the bigblueswarm config
package config

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func loadConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error(fmt.Sprintf("unable to close config file: %s", err))
		}
	}()

	conf := &Config{}

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
