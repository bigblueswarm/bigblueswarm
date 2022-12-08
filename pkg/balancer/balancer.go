// Package balancer manage the balancer progress and choose the next server
package balancer

import (
	"context"
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/utils"
	influxdb "github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// Balancer find the right server to use
type Balancer interface {
	Process(instances []string) (string, error)
	ClusterStatus(instances []string) ([]InstanceStatus, error)
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

// Process compute data to find a bigbluebutton server
func (b *InfluxDBBalancer) Process(instances []string) (string, error) {
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
		log.Error("Failed to find a valid server", err)
		return "", err
	}

	return extractBalancerResult(result), nil
}

// ClusterStatus retrieve the cluster status. It returns a list containing all bbb instance with its status
func (b *InfluxDBBalancer) ClusterStatus(instances []string) ([]InstanceStatus, error) {
	req := fmt.Sprintf(`
	from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r["_measurement"] == "bigbluebutton_api" or r["_measurement"] == "bigbluebutton_meetings")
		|> filter(fn: (r) => r["_field"] == "online" or r["_field"] == "active_meetings" or r["_field"] == "participant_count")
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
		log.Error("Failed to retrieve cluster status", err)
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
				Host:               instance,
				CPU:                0,
				Mem:                0,
				ActiveMeeting:      0,
				ActiveParticipants: 0,
				APIStatus:          "Down",
			}
			instanceMap[instance] = status
		} else {
			status = instanceMap[instance]
		}

		res := result.Record().Result()
		switch res {
		case "bbb":
			instanceMap[instance].ActiveMeeting = result.Record().ValueByKey("active_meetings").(int64)
			instanceMap[instance].ActiveParticipants = result.Record().ValueByKey("participant_count").(int64)
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
