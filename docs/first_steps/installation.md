# Installation

## BigBlueButton

For the installation of BigBlueButton, please refer to the [official documentation of BigBlueButton](https://docs.bigbluebutton.org/).

## InfluxDB

To install InfluxDB, please refer to the [official InfluxDB documentation](https://docs.influxdata.com/influxdb).

Once the installation of InfluxDB is done, create an API token through the [UI](https://docs.influxdata.com/influxdb/cloud/security/tokens/create-token/) or by following the instructions below:

```bash
export INFLUXDB_ORG=bigblueswarm # BigBlueSwarm InfluxDB organization
export INFLUXDB_BUCKET=bucket # BigBlueSwarm InfluxDB bucket
export INFLUXDB_TOKEN=Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg== # InfluxDB API token.
influx setup --name bigblueswarm --org $INFLUXDB_ORG --username admin --password password --token $INFLUX_TOKEN --bucket $INFLUXDB_BUCKET --retention 0 --force
``` 

## Telegraf

The installation of Telegraf must be done on each BigBlueButton server of the cluster. Please refer to the official documentation of [Telegraf](https://docs.influxdata.com/telegraf) to perform the installation.

BigBlueSwarm also requires the installation of the external Telegraf plugin for [BigBlueButton](https://github.com/bigblueswarm/bigbluebutton-telegraf-plugin). Follow the installation documentation. This plugin allows storing in InfluxDB the functional data of the cluster such as the number of meetings or the number of users.

Once installed, add the following minimal configuration:
```toml
[global_tags]
  bigblueswarm_host="${BIGBLUESWARM_HOST}" # don't miss this tag

# Configuration for telegraf agent
[agent]
  ## Default data collection interval for all inputs
  interval = "10s"
  ## Rounds collection interval to 'interval'
  ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
  round_interval = true

  ## Telegraf will send metrics to outputs in batches of at most
  ## metric_batch_size metrics.
  ## This controls the size of writes that Telegraf sends to output plugins.
  metric_batch_size = 1000

  ## For failed writes, telegraf will cache metric_buffer_limit metrics for each
  ## output, and will flush this buffer on a successful write. Oldest metrics
  ## are dropped first when this buffer fills.
  ## This buffer only fills when writes fail to output plugin(s).
  metric_buffer_limit = 10000

  ## Collection jitter is used to jitter the collection by a random amount.
  ## Each plugin will sleep for a random time within jitter before collecting.
  ## This can be used to avoid many plugins querying things like sysfs at the
  ## same time, which can have a measurable effect on the system.
  collection_jitter = "0s"

  ## Default flushing interval for all outputs. Maximum flush_interval will be
  ## flush_interval + flush_jitter
  flush_interval = "10s"
  ## Jitter the flush interval by a random amount. This is primarily to avoid
  ## large write spikes for users running a large number of telegraf instances.
  ## ie, a jitter of 5s and interval 10s means flushes will happen every 10-15s
  flush_jitter = "0s"

  ## By default or when set to "0s", precision will be set to the same
  ## timestamp order as the collection interval, with the maximum being 1s.
  ##   ie, when interval = "10s", precision will be "1s"
  ##       when interval = "250ms", precision will be "1ms"
  ## Precision will NOT be used for service inputs. It is up to each individual
  ## service input to set the timestamp at the appropriate precision.
  ## Valid time units are "ns", "us" (or "µs"), "ms", "s".
  precision = ""

  ## Logging configuration:
  ## Run telegraf with debug log messages.
  debug = false
  ## Run telegraf in quiet mode (error log messages only).
  quiet = false
  ## Specify the log file name. The empty string means to log to stderr.
  logfile = ""

  ## Override default hostname, if empty use os.Hostname()
  hostname = ""
  ## If set to true, do no set the "host" tag in the telegraf agent.
  omit_hostname = false
[[outputs.influxdb_v2]]
  urls = ["${INFLUXDB_URL}"]

  ## Token for authentication.
  token = "${INFLUXDB_TOKEN}"

  ## Organization is the name of the organization you wish to write to; must exist.
  organization = "${INFLUXDB_ORG}"

  ## Destination bucket to write into.
  bucket = "${INFLUXDB_BUCKET}"
[[inputs.cpu]]
  ## Whether to report per-cpu stats or not
  percpu = true
  ## Whether to report total system cpu stats or not
  totalcpu = true
  ## If true, collect raw CPU time metrics.
  collect_cpu_time = false
  ## If true, compute and report the sum of all non-idle CPU states.
  report_active = false
[[inputs.mem]]
[[inputs.net]]

[[inputs.execd]]
 command = ["/path/to/bbb-telegraf", "-config", "/path/to/bbb-telegraf/config", "-poll_interval", "10s"]
 signal = "none"
```

Edit the `/etc/default/telegraf` file and add the following variables:
```bash
BIGBLUESWARM_HOST= # Your public bigbluebutton url like https://yourbbbhost/bigbluebutton
INFLUXDB_URL= # InfluxDB url
INFLUXDB_TOKEN=Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg== # Generated InfluxDB api token
INFLUXDB_ORG=bigblueswarm # InfluxDB organization
INFLUXDB_BUCKET=bucket # InfluxDB bucket
```

Restart the Telegraf probe:
```bash
systemctl restart telegraf
``` 

## Redis

For the installation of Redis, please refer to the [official documentation](https://redis.io/topics/quickstart). You must install Redis in persistent mode.

## BigBlueSwarm load balancer

Download the [latest version](https://github.com/bigblueswarm/bigblueswarm/releases) binary on Github.

Since version 2, BigBlueSwarm provides two configuration systems:
- yaml file
- network configuration service [Consul](https://www.hashicorp.com/products/consul)

### Configuration file

Add a configuration file on your server:
```yaml
bigblueswarm:
  secret: <your_bigbluebutton_secret>
admin:
  api_key: <your_admin_api_key>
balancer:
  metrics_range: -5m 
  cpu_limit: 100
  mem_limit: 100
port: 8090
redis:
  address:
  password:
  database: 0
influxdb:
  address:
  token: 
  organization:
  bucket:
```

Start the BigBlueSwarm server with the `-config` option pointing to the previously created file.

> BigBlueSwarm also supports [Consul](https://www.hashicorp.com/products/consul) as a configuration provider. Refer to the page dedicated to [configuration](CONFIGURATION.md).

### Docker

Alternatively, BigBlueSwarm can run with Docker

It is also possible to run [BigBlueSwarm with Docker](https://hub.docker.com/r/sledunois/bigblueswarm).

```sh
docker run -it --mount type=bind,source="$(pwd)/config.yml",target=/config.yml,readonly -p 8090:8090 sledunois/bigblueswarm:latest -config /config.yml
```


[Next page](configuration.md)
