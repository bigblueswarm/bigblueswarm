# Configuration

BigBlueSwarm has two types of configurations:
* Application configuration: configuration needed for the application to work properly. They are read when BigBlueSwarm is launched;
* Launching configuration: technical configurations that allow BigBlueSwarm to be launched

## Application configurations

> For the configuration using Consul as a configuration provider, please refer to the [official documentation](https://www.hashicorp.com/products/consul) for installation.

### Configuration list

By default, BigBlueSwarm uses a [YAML](https://yaml.org/) configuration file. The following chapters list the different configurations that this file must contain.

#### BigBlueButton

* `secret` - __String__ - Secret BigBlueButton. As BigBlueSwarm works as a proxy, it reproduces the behavior of a BigBlueButton server and its authentication system. This `secret` configuration represents the key used by BigBlueButton clients to authenticate requests.

Exemple:
```yml
bigbluebutton:
  secret: 0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU
```

#### Admin

* `api_key` - __String__ - API key used to consume the administration API. The configuration is also used by the `bbsctl` cli.

Exemple:

```yml
admin:
  api_key: kgpqrTipM2yjcXwz5pOxBKViE9oNX76R
```

#### Balancer

* `metrics_range` - __String__ - Server utilization calculation period. This period is used to calculate the utilization rate of BigBlueButton servers. The server selection algorithm is not based on CPU and memory usage values at a given time but on an average calculated from the `metrics_range` configuration. For example, a `5m` configuration will define an average CPU and memory usage percentage over the last 5 minutes.
* `cpu_limit` - __Integer__ - Maximum CPU usage of BigBlueButton servers. If the BigBlueButton server reaches this average CPU value over the configured `metrics_range` period, then the server is no longer available to create a new meeting.
* `mem_limit` - __Integer__ - Maximum memory usage of BigBlueButton servers. If the BigBlueButton server reaches this average memory value over the configured `metrics_range` pŕiode, then the server is no longer available when creating a new meeting.

Example:
```yml
balancer:
  metrics_range: -5m 
  cpu_limit: 100
  mem_limit: 100
```

#### Port

* __Integer__ - Listening port of BigBlueSwarm.

Exemple:
```yml
port: 8090
```

#### Redis

* `address` - __String__ - Address to access the Redis server.
* `password` - __String__ - Password to access the Redis database. If you do not use a password on the Redis database, leave this field blank.
* `database` - __Integer__ - Redis databases are numbered from 0 to 15. The default Redis database is 0. If you are not using a particular database, set the value to 0.

Exemple:
```yml
redis:
  address: localhost:6379
  password:
  database: 0
```

#### InfluxDB

* `address` - __String__ - Access address to the InfluxDB server.
* token` - __String__ - Access token to the InfluxDB server api.
* `organization` - __String__ - Organization configured during InfluxDB server installation.
* `bucket` - __String__ - InfluxDB Bucket configured during InfluxDB server installation.

Exemple:

```yml
influxdb:
  address: http://localhost:8086
  token: Zq9wLsmhnW5UtOiPJApUv1cTVJfwXsTgl_pCkiTikQ3g2YGPtS5HqsXef-Wf5pUU3wjY3nVWTYRI-Wc8LjbDfg==
  organization: bigblueswarm
  bucket: bucket
```

#### Sample configuration file

```yml
bigbluebutton:
  secret: 0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU
admin:
  api_key: kgpqrTipM2yjcXwz5pOxBKViE9oNX76R
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

### Consul

BigBlueSwarm can also use Consul as a configuration provider by using the Consul KV function. The use of Consul is preferred when deploying in cluster mode.

Here is the configuration mapping for using Consul:

| Configuration   | Endpoint                      | Type      | Autorefresh*       | Example                                                                                                                 |
| --------------- | ----------------------------- | --------- | ------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| `bigbluebutton` | `configuration/bigbluebutton` | code/YAML | :heavy_check_mark: | <pre><code>secret: 0ol5t44UR21rrP0xL5ou7IBFumWF3GENebgW1RyTfbU</code></pre>                                             |
| `admin`         | `configuration/admin`         | code/YAML | :heavy_check_mark: | <pre><code>api_key: kgpqrTipM2yjcXwz5pOxBKViE9oNX76R</code></pre>                                                       |
| `balancer`      | `configuration/balancer`      | code/YAML | :heavy_check_mark: | <pre><code>metrics_range: -5m</code><br /><code>cpu_limit: 100</code><br /><code>mem_limit: 100</code></pre>            |
| `port`          | `configuration/port`          | none      |                    | <pre><code>8090</code></pre>                                                                                            |
| `redis`         | `configuration/redis`         | code/YAML |                    | <pre><code>address: </code><br /><code>password:</code><br /><code>database: 0</code></pre>                             |
| `influxdb`      | `configuration/influxdb`      | code/YAML |                    | <pre><code>address: </code><br /><code>token:</code><br /><code>organization: 0</code><br /><code>bucket: </code></pre> |

> Autorefresh*: when the value on Consul is changed, it is automatically updated in BigBlueSwarm. This feature does not work for the port used by BigBlueSwarm and the database accesses.

## Configuration at launch

Here is the list of possible configurations at the launch of BigBlueSwarm:
* `-config` - path to the configuration. In the case of a YAML configuration file, point to the file. In the case of a Consul server, prefix the address with `consul:`, e.g. `-config consul:http://localhost:8500`
* `log.level` - The desired log level. By default, the log level is set to INFO. Accepted values: `panic`, `fatal`, `error`, `warn`, `info`, `debug`, or `trace`.
* `log.path` __Optional__ - Path to the log file. If the option is not set then the logs are displayed in the STDOUT.
* `sentry.dsn` - Data Source Name. The Sentry DSN is required for error reporting to Sentry.
* `sentry.rates` - Sentry trace sampling rate. Takes a float value between 0.0 and 1.0.

[Next page](initialization.md)
