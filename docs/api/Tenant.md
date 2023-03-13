# Tenant

A `Tenant` is an entity identifying a client by its public hostname, through the public entry point of BigBlueSwarm. At least one tenant is mandatory for BigBlueSwarm to work.

Sample tenant file:

```yml
kind: Tenant
spec:
    host: localhost
instances: []
```

## Specifications

Described by the `spec` object, client specifications may contain the following properties:
  * `host` - String - __Required__ - public host of the client. URL on which the client will respond.
  * `secret` - String - BigBlueSwarm's dedicated secret. BigBLueSwarm is configured with a default secret. However you can override it for a particular client by defining it in the client's specifications.
  * `meeting_pool` - Integer - Meeting limit for the client. Once this limit is reached, it will not be possible to create new meetings on the client.
  * user_pool - Integer - User limit for the client. Once this limit is reached, users will not be able to join meetings.

Example:
```yml
spec:
  host: localhost
  secret: dummy_secret
  meeting_pool: 100
  user_pool: 1000
```

## Metadata

Metadata is an optional part of a tenant. It allows you to add data that have no functional use for BigBlueSwarm but that help to identify the tenant. It is a map that take a string as a key and a string as a value.

Example:
```yml
metadata:
  name: dummy_name
  description: "this is a dummy description"
```

## Instances

Instances describe the BigBlueButton instances accessible by the tenant. This configuration allows you to restrict (or not) the access of the instances to certain tenants.

> It is important to note that leaving the list of instances empty gives the tenant access to all instances of the cluster.

Example:
```yml
# Providing instances for the tenant
instances:
  - http://bbb1.example.com/bigbluebutton
  - http://bbb2.example.com/bigbluebutton
  - http://bbb3.example.com/bigbluebutton

# Giving access to all instances of the cluster
instances: []
```

## Initialization

A tenant can be initialized using the command [`bbsctl init tenant --host my_tenant_hostname`](https://github.com/bigblueswarm/bbsctl/blob/main/docs/bbsctl_init_tenant.md).

## Example file

The following is an example of a complete file containing all the different parts filled in:

```yml
kind: Tenant
spec:
  host: localhost
  secret: dummy_secret
  meeting_pool: 100
  user_pool: 1000
metadata:
  name: dummy_name
  description: "this is a dummy description"
instances:
  - http://bbb1.example.com/bigbluebutton
  - http://bbb2.example.com/bigbluebutton
  - http://bbb3.example.com/bigbluebutton
```
