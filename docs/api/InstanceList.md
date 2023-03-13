# InstanceList

The `InstanceList` entity tells BigBlueSwarm the list of BigBlueButton instances available in its cluster.

## Instances

The `instances` property is a mandatory property that defines the list of BigBlueButton instances present in the cluster. It is presented as a map taking in key the public URLs of the instances (String) and in key the BigBlueButton secret of the instance.

Example:

```yml
instances:
  http://bbb1.com/bigbluebutton: my_dummy_secret1
  http://bbb2.com/bigbluebutton: my_dummy_secret2
  http://bbb3.com/bigbluebutton: my_dummy_secret3
```

## Initialization

The list of instances can be initialized using the command [`bbsctl init instances`](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_init_instances.md).

## Sample file

```yml
kind: InstanceList
instances:
  http://bbb1.com/bigbluebutton: my_dummy_secret1
  http://bbb2.com/bigbluebutton: my_dummy_secret2
  http://bbb3.com/bigbluebutton: my_dummy_secret3
```
