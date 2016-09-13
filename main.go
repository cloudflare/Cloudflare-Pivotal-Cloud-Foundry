package main

import (
    "net/http"

    "github.com/pivotal-cf/brokerapi"
    "code.cloudfoundry.org/lager"
)

type myServiceBroker struct {}

func (*myServiceBroker) Services() []brokerapi.Service {
    // Return a []brokerapi.Service here, describing your service(s) and plan(s)
    return []brokerapi.Service{}
}

func (*myServiceBroker) Provision(
    instanceID string,
    details brokerapi.ProvisionDetails,
    asyncAllowed bool,
) (brokerapi.ProvisionedServiceSpec, error) {
    // Provision a new instance here. If async is allowed, the broker can still
    // chose to provision the instance synchronously.

    return brokerapi.ProvisionedServiceSpec{}, nil
}

func (*myServiceBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
    // If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
    // for the status of the provisioning operation.
    // This also applies to deprovisioning (work in progress).
    return brokerapi.LastOperation{}, nil
}

func (*myServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
    // Deprovision a new instance here. If async is allowed, the broker can still
    // chose to deprovision the instance synchronously, hence the first return value.
    return brokerapi.DeprovisionServiceSpec{}, nil
}

func (*myServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
    // Bind to instances here
    // Return a binding which contains a credentials object that can be marshalled to JSON,
    // and (optionally) a syslog drain URL.

    return brokerapi.Binding{}, nil
}

func (*myServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
    // Unbind from instances here
    return nil
}

func (*myServiceBroker) Update(instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
    // Update instance here
    return brokerapi.UpdateServiceSpec{}, nil
}

func main() {
    serviceBroker := &myServiceBroker{}
    logger := lager.NewLogger("my-service-broker")
    credentials := brokerapi.BrokerCredentials{
        Username: "username",
        Password: "password",
    }

    brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
    http.Handle("/", brokerAPI)
    http.ListenAndServe(":3000", nil)
}