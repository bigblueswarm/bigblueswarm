# b3lb

## Possible influxdb query
Proof of concept:
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