package app

import (
	"context"
	"fmt"

	redisApi "github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// Balancer find the right server to use
type Balancer struct {
	Client redisApi.QueryAPI
}

// NewBalancer creates a new Balancer object
func NewBalancer(idb redisApi.QueryAPI) *Balancer {
	return &Balancer{
		Client: idb,
	}
}

func (b *Balancer) formatInstancesFilter(instances []string) string {
	var result string
	for i, instance := range instances {
		filter := fmt.Sprintf(`r["b3lb_host"] == "%s"`, instance)
		result = fmt.Sprintf("%s %s", result, filter)

		if i != (len(instances) - 1) {
			result = fmt.Sprintf("%s or ", result)
		}
	}

	return result
}

// Process compute data to find a bigbluebutton server
func (b *Balancer) Process(instances []string) (string, error) {
	req := fmt.Sprintf(`from(bucket: "bucket")
	|> range(start: -5m)
	|> filter(fn: (r) => r["_measurement"] == "cpu" or r["_measurement"] == "mem")
	|> filter(fn: (r) => r["_field"] == "usage_system" or r["_field"] == "used_percent")
	|> filter(fn: (r) => r["cpu"] == "cpu-total" or r["_measurement"] == "mem")
	|> filter(fn: (r) => %s)
	|> group(columns: ["b3lb_host"])
	|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
	|> map(fn: (r) => ({ r with _value: r["usage_system"] + r["used_percent"] }))
	|> lowestAverage(n: 1, column: "_value", groupColumns: ["b3lb_host", "_time"])`, b.formatInstancesFilter(instances))
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
