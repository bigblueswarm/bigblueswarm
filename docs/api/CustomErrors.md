# Custom errors

The multi-client management and the different restrictions possible on these clients have created custom errors in addition to the classic BigBlueButton errors.

The following list describes these custom errors:

- `noInstanceFound`: if the balancer process returns an error or can't find an instance to create a meeting, it returns a `noInstanceFound`
```xml
<response>
    <returncode>FAILED</returncode>
    <messageKey>noInstanceFound</messageKey>
    <message>BigBlueSwarm do not find a valid BigBlueButton instance for your request</message>
</response>
```

- `tenantNotFound`: on each request, BigBlueSwarm checks if the tenant you try to call exists in the configuration. If it does not, it returns a `tenantNotFound` error. To fix this error, refer to the cli documentation and add a new tenant.
```xml
<response>
    <returncode>FAILED</returncode>
    <messageKey>tenantNotFound</messageKey>
    <message>BigBlueSwarm does not find the requesting  tenant</message>
</response>
```

- `meetingPoolReached`: this error appears when your tenant reaches the meeting pool limitation configured in the tenant configuration.
```xml
<response>
    <returncode>FAILED</returncode>
    <messageKey>meetingPoolReached</messageKey>
    <message>Your tenant reached the meeting pool limit and can't create a new one.</message>
</response>
```

- `userPoolReached`: this error appears when your tenant reaches the user pool limitation configured in the tenant configuration.
```xml
<response>
    <returncode>FAILED</returncode>
    <messageKey>userPoolReached</messageKey>
    <message>Your tenant reached the user pool limit.</message>
</response>
```

- `internalError`: each time BigBlueSwarm encounters an internal error, it returns an `internalError` message. Please refers to the message and check to logs. For example:
```xml
<response>
    <returncode>FAILED</returncode>
    <messageKey>internalError</messageKey>
    <message>BigBlueSwarm failed to retrieve the requesting tenant</message>
</response>
```
