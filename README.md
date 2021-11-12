# b3lb
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/c4c4627abd1f474fb2200f9831dfe502)](https://www.codacy.com/gh/SLedunois/b3lb/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=SLedunois/b3lb&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/c4c4627abd1f474fb2200f9831dfe502)](https://www.codacy.com/gh/SLedunois/b3lb/dashboard?utm_source=github.com&utm_medium=referral&utm_content=SLedunois/b3lb&utm_campaign=Badge_Coverage)

`scripts` folder contains every script you need to build and develop.
*   `init.sh`: Download and install go dependencies
*   `build_image.sh`: Build docker image used in your local cluster
*   `cluster.sh`: Manage your local cluster

## Prerequisite

You need:
*   docker 19+
*   golang 1.15+

## Installation

First, launch the init script:
 ```sh
./scripts/init.sh
 ```

Once you download dependencies, build the docker image
```shell
./scripts/build_image.sh
```

Then launch your local cluster with:
```shell
./scripts/cluster.sh --start
```

Init the cluster with:
```shell
./scripts/cluster.sh --init
```
It initializes the influx db cluster by creating a user, an organization and a bucket. It also set the token into BigBlueButton containers.
By default, it creates the user `admin` with the password `password`

Call api using `api.sh` script. For example create a room using:
```shell
./scripts/api.sh --create
```
For more information, use the help method.
```shell
./scripts/api.sh --help
```

Call admin function using `admin.sh` script. For example, add a bigbluebutton instance using:
```shell
./scripts/admin.sh --create
```

For more information, use the help method.
```shell
./scripts/admin.sh --help
```

## POC

### Benchmarks

*   [Octopuce benchmark](https://www.octopuce.fr/retour-dexperience-sur-bigbluebutton-a-fort-charge/)
*   [Aufood benchmark](https://www.aukfood.fr/faire-un-stress-test-sur-bigbluebutton/)

### Query example
```influxQL
from(bucket: "bucket")
  |> range(start: -3m)
  |> filter(fn: (r) => r["_measurement"] == "cpu")
  |> filter(fn: (r) => r["_field"] == "usage_system")
  |> filter(fn: (r) => r["cpu"] == "cpu-total")
  |> mean(column: "_value")
  |> group(columns: ["_time"])
  |> top(n:1, columns:["_value"])
```

```influxQL
from(bucket: "bucket")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "cpu" or r["_measurement"] == "mem")
  |> filter(fn: (r) => r["_field"] == "usage_system" or r["_field"] == "used_percent")
  |> filter(fn: (r) => r["cpu"] == "cpu-total" or r["_measurement"] == "mem")
  |> filter(fn: (r) => r["host"] == "28fca396fba6" or r["host"] == "d25f727605e0")
  |> group(columns: ["host"])
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> map(fn: (r) => ({ r with _value: r["usage_system"] + r["used_percent"] }))
  |> lowestAverage(n: 1, column: "_value", groupColumns: ["host", "_time"])
```

```influxQL
from(bucket: "bucket")
  |> range(start: -5m)
  |> filter(fn: (r) => r["_measurement"] == "cpu" or r["_measurement"] == "mem")
  |> filter(fn: (r) => r["_field"] == "usage_system" or r["_field"] == "used_percent")
  |> filter(fn: (r) => r["cpu"] == "cpu-total" or r["_measurement"] == "mem")
  |> filter(fn: (r) => r["b3lb_host"] == "http://localhost/bigbluebutton" or r["b3lb_host"] == "http://localhost:8080/bigbluebutton")
  |> group(columns: ["b3lb_host"])
  |> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
  |> map(fn: (r) => ({ r with _value: r["usage_system"] + r["used_percent"] }))
  |> lowestAverage(n: 1, column: "_value", groupColumns: ["b3lb_host", "_time"])
```

Check join method: https://docs.influxdata.com/influxdb/v2.1/reference/syntax/flux/flux-vs-influxql/#joins

## InfluxDB client

https://github.com/influxdata/influxdb-client-go#how-to-use
