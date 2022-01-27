package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/SLedunois/b3lb/pkg/config"
	redisApi "github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// Balancer find the right server to use
type Balancer interface {
	Process(instances []string) (string, error)
}

// InfluxDBBalancer is the InfluxDB implementation of Balancer
type InfluxDBBalancer struct {
	Client redisApi.QueryAPI
	Config *config.BalancerConfig
}

// NewBalancer creates a new Balancer object
func NewBalancer(idb redisApi.QueryAPI, config *config.BalancerConfig) Balancer {
	return &InfluxDBBalancer{
		Client: idb,
		Config: config,
	}
}

func (b *InfluxDBBalancer) formatInstancesFilter(instances []string) string {
	var result string
	for i, instance := range instances {
		filter := fmt.Sprintf(`r["b3lb_host"] == "%s"`, instance)
		result = fmt.Sprintf("%s %s", result, filter)

		if i != (len(instances) - 1) {
			result = fmt.Sprintf("%s or", result)
		}
	}

	return strings.TrimSpace(result)
}

// Process compute data to find a bigbluebutton server
func (b *InfluxDBBalancer) Process(instances []string) (string, error) {
	req := fmt.Sprintf(`from(bucket: "bucket")
	|> range(start: %s)
	|> filter(fn: (r) => r["_measurement"] == "cpu" or r["_measurement"] == "mem")
	|> filter(fn: (r) => r["_field"] == "usage_system" or r["_field"] == "used_percent")
	|> filter(fn: (r) => r["cpu"] == "cpu-total" or r["_measurement"] == "mem")
	|> filter(fn: (r) => %s)
	|> group(columns: ["b3lb_host"])
	|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	|> map(fn: (r) => ({ r with _value: r["usage_system"] + r["used_percent"] }))
	|> lowestAverage(n: 1, column: "_value", groupColumns: ["b3lb_host", "_time"])`, b.Config.MetricsRange, b.formatInstancesFilter(instances))
	result, err := b.Client.Query(context.Background(), req)
	if err != nil || result.Err() != nil {
		log.Error("Failed to find a valid server", err)
		return "", err
	}

	if result.Next() {
		return fmt.Sprintf("%v", result.Record().ValueByKey("b3lb_host")), nil
	}

	return "", nil
}
