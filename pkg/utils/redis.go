// Package utils provide few utilies functions
package utils

import "github.com/go-redis/redis/v8"

// ComputeErr manager redis error return
func ComputeErr(err error) error {
	if err == redis.Nil {
		return nil
	}

	return err
}
