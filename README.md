# b3lb
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/c4c4627abd1f474fb2200f9831dfe502)](https://www.codacy.com/gh/SLedunois/b3lb/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=SLedunois/b3lb&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/c4c4627abd1f474fb2200f9831dfe502)](https://www.codacy.com/gh/SLedunois/b3lb/dashboard?utm_source=github.com&utm_medium=referral&utm_content=SLedunois/b3lb&utm_campaign=Badge_Coverage)

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
