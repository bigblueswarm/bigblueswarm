package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
)

func isConsulEnabled(path string) bool {
	return strings.HasPrefix(path, ConsulPrefix)
}

func getConsulConfig(path string) *api.Config {
	config := api.DefaultConfig()
	addr := strings.ReplaceAll(path, ConsulPrefix, "")
	config.Address = addr

	return config
}

func loadConfigFromConsul(path string) (*Config, error) {
	// Get a new consul client
	client, err := api.NewClient(getConsulConfig(path))
	if err != nil {
		return nil, err
	}

	// Get a handle to the KV API
	kv := client.KV()

	conf := &Config{}
	if err := conf.loadBBBConf(kv); err != nil {
		return nil, err
	}

	if err := conf.loadAdminConf(kv); err != nil {
		return nil, err
	}

	if err := conf.loadBalancerConf(kv); err != nil {
		return nil, err
	}

	if err := conf.loadPortConf(kv); err != nil {
		return nil, err
	}

	if err := conf.loadRedisConf(kv); err != nil {
		return nil, err
	}

	if err := conf.loadInfluxDBConf(kv); err != nil {
		return nil, err
	}

	return conf, nil
}

func consulKey(conf string) string {
	return fmt.Sprintf("configuration/%s", conf)
}

func loadKey(kv *api.KV, key string) (interface{}, error) {
	pair, _, err := kv.Get(consulKey(key), nil)

	if err != nil {
		return nil, err
	}

	var result interface{}
	if key != "port" {
		result = getConfigType(key)
		if err := yaml.Unmarshal(pair.Value, result); err != nil {
			return nil, err
		}
	} else {
		result, err = strconv.Atoi(string(pair.Value))
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func getConfigType(key string) interface{} {
	switch key {
	case "bigbluebutton":
		return &BigBlueButton{}
	case "admin":
		return &AdminConfig{}
	case "balancer":	
		return &BalancerConfig{}
	case "redis":
		return &RDB{}
	case "influxdb":
		return &IDB{}
	default:
		return nil
	}
}

func watchChanges(key string, handler func(value []byte)) error {
	params := map[string]interface{}{
		"type": "key",
		"key":  consulKey(key),
	}

	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}

	plan.Handler = func(u uint64, raw interface{}) {
		var pair *api.KVPair
		if raw == nil {
			pair = nil
		} else {
			var ok bool
			if pair, ok = raw.(*api.KVPair); !ok {
				log.Error("unable to cast handled object as KVPair")
				return // ignore
			}
		}

		handler(pair.Value)
	}

	go func() {
		if err := plan.Run(api.DefaultConfig().Address); err != nil {
			log.Error(fmt.Errorf("err watching consul key: %v", err))
		}
	}()

	return nil
}

func (c *Config) loadBBBConf(kv *api.KV) error {
	key := "bigbluebutton"
	conf, err := loadKey(kv, key)
	if err != nil {
		return err
	}

	if value, ok := conf.(*BigBlueButton); ok {
		c.BigBlueButton = *value
	}

	return watchChanges(key, func(value []byte) {
		var conf BigBlueButton
		if err := yaml.Unmarshal(value, &conf); err != nil {
			log.Error(fmt.Errorf("unable to parse new config value: %s", err))
			return
		}

		c.BigBlueButton = conf
	})
}

func (c *Config) loadAdminConf(kv *api.KV) error {
	key := "admin"
	conf, err := loadKey(kv, key)
	if err != nil {
		return err
	}

	if value, ok := conf.(*AdminConfig); ok {
		c.Admin = *value
	}

	return watchChanges(key, func(value []byte) {
		var conf AdminConfig
		if err := yaml.Unmarshal(value, &conf); err != nil {
			log.Error(fmt.Errorf("unable to parse new config value: %s", err))
			return
		}

		c.Admin = conf
	})
}

func (c *Config) loadBalancerConf(kv *api.KV) error {
	key := "balancer"
	conf, err := loadKey(kv, key)
	if err != nil {
		return err
	}

	if value, ok := conf.(*BalancerConfig); ok {
		c.Balancer = *value
	}

	return watchChanges(key, func(value []byte) {
		var conf BalancerConfig
		if err := yaml.Unmarshal(value, &conf); err != nil {
			log.Error(fmt.Errorf("unable to parse new config value: %s", err))
			return
		}

		c.Balancer = conf
	})
}

func (c *Config) loadPortConf(kv *api.KV) error {
	conf, err := loadKey(kv, "port")
	if err != nil {
		return err
	}

	c.Port = Port(conf.(int))

	return nil
}

func (c *Config) loadRedisConf(kv *api.KV) error {
	conf, err := loadKey(kv, "redis")
	if err != nil {
		return err
	}

	if value, ok := conf.(*RDB); ok {
		c.RDB = *value
	}

	return nil
}

func (c *Config) loadInfluxDBConf(kv *api.KV) error {
	conf, err := loadKey(kv, "influxdb")
	if err != nil {
		return err
	}

	if value, ok := conf.(*IDB); ok {
		c.IDB = *value
	}

	return nil
}
