# CloudFlare-Pivotal-Cloud-Foundry

## Usage

Be sure you have `$GOPATH` setup.

Set environment variables
```
export BROKER_USERNAME=username
export BROKER_PASSWORD=password
export PORT=9000
```

`go run main.go` runs the service on localhost.
`go test ./...` runs tests.

## Example cURL calls

In the example the following attributes are replaced with the corresponding value.

* username is `username`
* password is `password`
* broker-url is `localhost`
* port is `9000`
* instance_id is `1`.

### Service
```
curl -H "X-Broker-API-Version: 2.10" http://username:password@localhost:9000/v2/catalog
```

### Provision

```
curl http://username:password@localhost:9000/v2/service_instances/1 -d '{
  "organization_guid": "org-guid-here",
  "plan_id":           "plan-guid-here",
  "service_id":        "service-guid-here",
  "space_guid":        "space-guid-here",
  "parameters":        {
    "x-auth-key": "mykey",
    "x-auth-email": "email@email.com"
  }
}' -X PUT -H "X-Broker-API-Version: 2.10" -H "Content-Type: application/json"
```

### Deprovision

```
curl 'http://username:password@localhost:9000/v2/service_instances/1?service_id=service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.10"
```

### Bind

* Assumed binding_id as `2`

```
curl http://username:password@localhost:9000/v2/service_instances/1/service_bindings/2 -d '{
  "plan_id":      "plan-guid-here",
  "service_id":   "service-guid-here",
  "app_guid":     "app-guid-here",
  "bind_resource":     {
    "app_guid": "app-guid-here"
  },
  "parameters":        {
    "domain": "domain.com"
  }
}' -X PUT
```

### Unbind

* Assumed binding_id as `2`

```
curl 'http://username:password@localhost:9000/v2/service_instances/1/service_bindings/2?service_id=service-id-here&plan_id=plan-id-here' -X DELETE -H "X-Broker-API-Version: 2.10"
```

### Last Operation (Not supported)
```
curl http://username:password@localhost:9000/v2/service_instances/1/last_operation
```


### Update (Not supported)

```
curl http://username:password@localhost:9000/v2/service_instances/1 -d '{
  "service_id": "service-guid-here",
  "plan_id": "plan-guid-here",
  "parameters":        {
    "parameter1": 1,
    "parameter2": "value"
  },
  "previous_values": {
    "plan_id": "old-plan-guid-here",
    "service_id": "service-guid-here",
    "organization_id": "org-guid-here",
    "space_id": "space-guid-here"
  }
}' -X PATCH -H "X-Broker-API-Version: 2.10" -H "Content-Type: application/json"
```