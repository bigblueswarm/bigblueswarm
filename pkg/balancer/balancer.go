// Package balancer manage the balancer progress and choose the next server
package balancer

import (
	"context"
	"errors"
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"
	influxdb "github.com/influxdata/influxdb-client-go/v2/api"
)

// Balancer fcheck the cluster status
type Balancer interface {
	// Process compute data to find a bigbluebutton server
	Process(instances []string) (string, error)
	// ClusterStatus retrieve the cluster status. It returns a list containing all bbb instance with its status
	ClusterStatus(instances []string) ([]InstanceStatus, error)
	// GetCurrentState retrieve the measurement state in cluster
	GetCurrentState(measurement string, field string) (int64, error)
}

// InfluxDBBalancer is the InfluxDB implementation of Balancer
type InfluxDBBalancer struct {
	Client    influxdb.QueryAPI
	Config    *config.BalancerConfig
	IDBConfig *config.IDB
}

// New creates a new Balancer object
func New(idb influxdb.QueryAPI, config *config.BalancerConfig, idbConfig *config.IDB) Balancer {
	return &InfluxDBBalancer{
		Client:    idb,
		Config:    config,
		IDBConfig: idbConfig,
	}
}

func (b *InfluxDBBalancer) filterOnlineInstances(instances []string) ([]string, error) {
	req := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "bigbluebutton" and r["_field"] == "online")
		|> filter(fn: (r) => %s)
		|> filter(fn: (r) => r["_value"] == 1)
		|> last()
		|> yield(name: "online")
	`,
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances))

	result, err := b.Client.Query(context.Background(), req)
	if err != nil {
		return nil, err
	}

	values := []string{}
	for result.Next() {
		if result.Record().Result() == "online" {
			values = append(values, result.Record().ValueByKey("bigblueswarm_host").(string))
		}
	}

	return values, nil
}

// Process compute data to find a bigbluebutton server
func (b *InfluxDBBalancer) Process(instances []string) (string, error) {
	instances, err := b.filterOnlineInstances(instances)
	if err != nil {
		return "", err
	}

	if len(instances) == 0 {
		return "", errors.New("no instance online to process a balancer request")
	}

	req := fmt.Sprintf(`
	cpuFilter = from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "cpu" and r["_field"] == "usage_system" and r["cpu"] == "cpu-total")
		|> filter(fn: (r) => %s)
		|> group(columns: ["bigblueswarm_host"])
		|> mean(column: "_value")
		|> yield(name: "cpu")
  
	memFilter = from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "mem" and r["_field"] == "used_percent")
		|> filter(fn: (r) => %s)
		|> group(columns: ["bigblueswarm_host"])
		|> mean(column: "_value")
		|> yield(name: "mem")
	
	join(
		tables: {mem: memFilter, cpu: cpuFilter},
		on: ["bigblueswarm_host", "_start", "_stop"],
	)
	|> filter(fn: (r) => r["_value_cpu"] <= %d and r["_value_mem"] <= %d)
	|> map(fn: (r) => ({ r with _value: r["_value_cpu"] + r["_value_mem"] }))
	|> lowestAverage(n: 1, column: "_value", groupColumns: ["bigblueswarm_host", "_time"])
	|> yield(name: "balancer")
	`,
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances),
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances),
		b.Config.CPULimit,
		b.Config.MemLimit,
	)
	result, err := b.Client.Query(context.Background(), req)
	if err != nil || result.Err() != nil {
		return "", err
	}

	return extractBalancerResult(result), nil
}

// ClusterStatus retrieve the cluster status. It returns a list containing all bbb instance with its status
func (b *InfluxDBBalancer) ClusterStatus(instances []string) ([]InstanceStatus, error) {
	req := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "bigbluebutton")
		|> filter(fn: (r) => r["_field"] == "online" or r["_field"] == "meetings" or r["_field"] == "participants")
		|> filter(fn: (r) => %s)
		|> group(columns:["bigblueswarm_host", "_start"])
		|> pivot(rowKey: ["_start"], columnKey: ["_field"], valueColumn: "_value")
		|> last(column: "_start")
		|> yield(name: "bbb")
		
	from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "cpu" and r["_field"] == "usage_system" and r["cpu"] == "cpu-total")
		|> filter(fn: (r) => %s)
		|> group(columns:["bigblueswarm_host", "_start"])
		|> mean()
		|> yield(name: "cpu")
	
	from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "mem" and r["_field"] == "used_percent")
		|> filter(fn: (r) => %s)
		|> group(columns:["bigblueswarm_host", "_start"])
		|> mean()
		|> yield(name: "mem")
	`,
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances),
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances),
		b.IDBConfig.Bucket,
		b.Config.MetricsRange,
		utils.FormatInstancesFilter(instances),
	)

	result, err := b.Client.Query(context.Background(), req)
	if err != nil || result.Err() != nil {
		return nil, err
	}

	return parseClusterStatusResult(result), nil
}

func parseClusterStatusResult(result *influxdb.QueryTableResult) []InstanceStatus {
	instanceMap := make(map[string]*InstanceStatus)
	for result.Next() {
		instance := result.Record().ValueByKey("bigblueswarm_host").(string)
		var status *InstanceStatus
		if _, ok := instanceMap[instance]; !ok {
			status = &InstanceStatus{
				Host:         instance,
				CPU:          0,
				Mem:          0,
				Meetings:     0,
				Participants: 0,
				APIStatus:    "Down",
			}
			instanceMap[instance] = status
		} else {
			status = instanceMap[instance]
		}

		res := result.Record().Result()
		switch res {
		case "bbb":
			instanceMap[instance].Meetings = result.Record().ValueByKey("meetings").(int64)
			instanceMap[instance].Participants = result.Record().ValueByKey("participants").(int64)
			instanceMap[instance].APIStatus = apiStatusToString(result.Record().ValueByKey("online").(int64))
		case "mem":
			instanceMap[instance].Mem = utils.Round2Digits(result.Record().Value().(float64))
		case "cpu":
			instanceMap[instance].CPU = utils.Round2Digits(result.Record().Value().(float64))
		}
	}

	instances := []InstanceStatus{}
	for _, instance := range instanceMap {
		instances = append(instances, *instance)
	}

	return instances
}

func apiStatusToString(status int64) string {
	if status == 1 {
		return "Up"
	}

	return "Down"
}

func extractBalancerResult(result *influxdb.QueryTableResult) string {
	instance := ""
	for result.Next() {
		if result.Record().Result() == "balancer" {
			return result.Record().ValueByKey("bigblueswarm_host").(string)
		}
	}

	return instance
}

// GetCurrentState retrieve the measurement state in cluster
func (b *InfluxDBBalancer) GetCurrentState(measurement string, field string) (int64, error) {
	q := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: -60s)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["_field"] == "%s")
		|> aggregateWindow(every: %s, fn: sum, createEmpty: false)
		|> last()
	`, b.IDBConfig.Bucket, measurement, field, b.Config.AggregationInterval)

	result, err := b.Client.Query(context.Background(), q)

	if err != nil {
		return -1, fmt.Errorf("failed to retrieve current state for measurement %s and field %s: %s", measurement, field, err)
	}

	val := int64(0)
	if result.Next() {
		val = result.Record().Value().(int64)
	}

	if result.Err() != nil {
		return -1, fmt.Errorf("get current state returns an error for measurement %s and field %s: %s", measurement, field, result.Err())
	}

	return val, nil
}
