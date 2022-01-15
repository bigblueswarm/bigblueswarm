package utils

import (
	"github.com/SLedunois/b3lb/pkg/config"

	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// RedisClient initilize a redis client based on provided configuration
func RedisClient(conf *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.RDB.Address,
		Password: conf.RDB.Password,
		DB:       conf.RDB.DB,
	})
}

// InfluxDBClient initilize an influxdb client based on provided configuration
func InfluxDBClient(conf *config.Config) api.QueryAPI {
	client := influxdb2.NewClient(conf.IDB.Address, conf.IDB.Token)
	return client.QueryAPI(conf.IDB.Organization)
}
