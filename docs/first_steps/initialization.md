# Initialize your cluster

To manage the BigBlueSwarm cluster, it is best to install the [dedicated cli](https://github.com/bigblueswarm/bbsctl).

## Configuring the CLI

Once the dedicated cli is installed, you can initialize it using the `bbsctl init config` command with the `--bbs <your_bigblueswarm_host>` and `--key <your_bigblueswarm_api_key>` options (see the [documentation](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_init_config.md)).

For example:
```sh
bbsctl init config --bbs http://localhost:8090 --key kgpqrTipM2yjcXwz5pOxBKViE9oNX76R
```

This command will generate the following configuration file `$HOME/.bigblueswarm/.bbsctl`:
```yml
bbs: http://localhost:8090
apiKey: kgpqrTipM2yjcXwz5pOxBKViE9oNX76R
```

<a href="https://asciinema.org/a/2qz4H250QCzMCbioMqEsnkuVE" target="_blank"><img src="https://asciinema.org/a/2qz4H250QCzMCbioMqEsnkuVE.svg" height="300" /></a>


## Add BigBlueButton instances

The addition of BigBlueButton instances is done through a YAML configuration file. To create the file, run the command ([see documentation](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_init_instances.md)):
```sh
bbsctl init instances
```

The command will generate an `instances.yml` file in the BigBlueSwarm configuration folder `$HOME/.bigblueswarm` (default configuration folder):

```yml
kind: InstanceList
instances: {}
```

Complete the `instances` part with your BigBlueButton instances using the host of the instance as key and its secret as value, for example:
```yml
kind: InstanceList
instances:
  http://bbb1.com/bigbluebutton: my_dummy_secret
```

Then add the instances to your cluster using the `apply` command ([see documentation](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_apply.md)):

```sh
bbsctl apply -f ~/.bigblueswarm/instances.yml
```

Validate that your instances have been added ([see documentation](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_get_instances.md)):

```sh
bbsctl get instances
```

<a href="https://asciinema.org/a/nlueByPpY8E9g7DElk85jRcqH" target="_blank"><img src="https://asciinema.org/a/nlueByPpY8E9g7DElk85jRcqH.svg" height="300" /></a>

When you change your instance list, you can update it using the `apply` command.

## Create a client

BigBlueSwarm is a multi-client load balancer and cannot work without a client. So we will initialize one. To do this, issue the command [`init tenant`](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_init_tenant.md):

```sh
bbsctl init tenant --host my_public_host
```

This command will generate the following configuration file `my_public_host.tenant.yml`:
```yml
kind: Tenant
spec:
    host: my_public_host
instances: []
``` 

To configure the `spec` part, go to the [dedicated page](../api/TENANTS.md).
By default the list of instances available on your client is empty: this means that the client is allowed to use all the instances of the cluster. However, if you want to assign particular instances to it, you can add them to the `instances` list, for example:
```yml
kind: Tenant
spec:
    host: my_public_host
instances:
  - http://bbb1.com/bigbluebutton
  - http://bbb2.com/bigbluebutton
```

Once the client is configured, we need to add it to the cluster using the `apply` command:
```sh
bbsctl apply -f ~/.bigblueswarm/my_public_host.tenant.yml
```

Check that your client is added using the commands:
```sh
bbstcl get tenants
bbsctl describe tenant my_public_host
```

When you change your client, you can update it using the `apply` command.

<a href="https://asciinema.org/a/gB55DijzRQCl40or7bgpa3sv5" target="_blank"><img src="https://asciinema.org/a/gB55DijzRQCl40or7bgpa3sv5.svg" height="300" /></a>

## Monitoring your cluster

Once your cluster is configured and in use, you can monitor it using the [`cluster-info`](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_cluster-info.md) command.

```sh
bbsctl cluster-info
```

<a href="https://asciinema.org/a/Nqec46FDprZpUzbP43oa940Xb" target="_blank"><img src="https://asciinema.org/a/Nqec46FDprZpUzbP43oa940Xb.svg" height="300" /></a>


## Checking the BigBlueSwarm configuration

If you want to check the BigBlueSwarm configuration remotely, you can use the [`describe config`](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_describe_config.md) command.

```sh
bbsctl describe config
```

<a href="https://asciinema.org/a/idAmn6AjWZIh77x1oWcaK8Yiv" target="_blank"><img src="https://asciinema.org/a/idAmn6AjWZIh77x1oWcaK8Yiv.svg" height="300" /></a>
