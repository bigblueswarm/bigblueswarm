# b3lb
[![codecov](https://codecov.io/gh/SLedunois/b3lb/branch/main/graph/badge.svg?token=PUP086DAUJ)](https://codecov.io/gh/SLedunois/b3lb)

`scripts` folder contains every script you need to build and develop.
1. `init.sh`: Download and install go dependencies
2. `build_image.sh`: Build docker image used in your local cluster
3. `cluster.sh`: Manage your local cluster

## Prerequisite

You need:
* docker 19+
* golang 1.15+

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

Connect on InfluxDB server and create an authentication token then configure your cluster to use it
```shell
./scripts/cluster.sh --set-token [token]
```

## POC

### Benchmarks

* [Octopuce benchmark](https://www.octopuce.fr/retour-dexperience-sur-bigbluebutton-a-fort-charge/)
* [Aufood benchmark](https://www.aukfood.fr/faire-un-stress-test-sur-bigbluebutton/)

### Query example
```
from(bucket: "bucket")
  |> range(start: -3m)
  |> filter(fn: (r) => r["_measurement"] == "cpu")
  |> filter(fn: (r) => r["_field"] == "usage_system")
  |> filter(fn: (r) => r["cpu"] == "cpu-total")
  |> mean(column: "_value")
  |> group(columns: ["_time"])
  |> top(n:1, columns:["_value"])
```
