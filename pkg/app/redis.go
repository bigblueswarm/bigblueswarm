package app

import (
	"b3lb/pkg/config"

	"github.com/go-redis/redis/v8"
)

func redisClient(conf *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.RDB.Address,
		Password: conf.RDB.Password,
		DB:       conf.RDB.DB,
	})
}
