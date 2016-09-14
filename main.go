package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

const BROKER_PASSWORD = "BROKER_PASSWORD"
const BROKER_PORT = "PORT"
const BROKER_USERNAME = "BROKER_USERNAME"
const CLOUDFLARE_EMAIL = "CLOUDFLARE_EMAIL"
const CLOUDFLARE_API_KEY = "CLOUDFLARE_API_KEY"
const X_AUTH_EMAIL_HEADER = "X-Auth-Email"
const X_AUTH_KEY_HEADER = "X-Auth-Key"

type myServiceBroker struct {
	logger lager.Logger
}

func (*myServiceBroker) Services() []brokerapi.Service {
	return []brokerapi.Service{
		brokerapi.Service{
			Name:        "cloudflare",
			ID:          "cloudflare",
			Description: "Give us five minutes and we’ll supercharge your website.",
			Bindable:    true,
			Plans: []brokerapi.ServicePlan{
				brokerapi.ServicePlan{
					ID:          "cloudflare-free",
					Name:        "cloudflare-free",
					Description: "Fast site performance. Broad security protection. SSL. Powerful stats about your visitors. Peace of mind about running your website so you can get back to what you love.",
				},
			},
		},
	}
}

func (m *myServiceBroker) Provision(
	instanceID string,
	details brokerapi.ProvisionDetails,
	asyncAllowed bool,
) (brokerapi.ProvisionedServiceSpec, error) {
	// Provision a new instance here. If async is allowed, the broker can still
	// chose to provision the instance synchronously.

	var authHeaders AuthHeaders
	if error := json.Unmarshal(details.RawParameters, &authHeaders); error != nil {
		m.logger.Error("Error decoding details.RawParameters", error)
		return brokerapi.ProvisionedServiceSpec{}, error
	}

	setAuthHeaders(authHeaders)

	return brokerapi.ProvisionedServiceSpec{
		IsAsync:       false,
		DashboardURL:  "",
		OperationData: "",
	}, nil
}

type AuthHeaders struct {
	XAuthEmail string `json:"x-auth-email"`
	XAuthKey   string `json:"x-auth-key"`
}

func getAuthHeaders() AuthHeaders {
	return AuthHeaders{
		XAuthEmail: os.Getenv(CLOUDFLARE_EMAIL),
		XAuthKey:   os.Getenv(CLOUDFLARE_API_KEY),
	}
}

func setAuthHeaders(authHeaders AuthHeaders) AuthHeaders {
	os.Setenv(CLOUDFLARE_EMAIL, authHeaders.XAuthEmail)
	os.Setenv(CLOUDFLARE_API_KEY, authHeaders.XAuthKey)

	return getAuthHeaders()
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

func (m *myServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// Bind to instances here
	// Return a binding which contains a credentials object that can be marshalled to JSON,
	// and (optionally) a syslog drain URL.

	_, isDomainInDictionary := details.Parameters["domain"]
	if isDomainInDictionary == false {
		return brokerapi.Binding{}, errors.New("key 'domain' not found in BindDetails.Parameters.")
	}
	var domain = details.Parameters["domain"].(string)

	var client = http.Client{}

	var jsonBody = []byte(`{"name" : "` + domain + `"}`)
	request, error := http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones", bytes.NewReader(jsonBody))

	var authHeaders = getAuthHeaders()

	request.Header.Add(X_AUTH_EMAIL_HEADER, authHeaders.XAuthEmail)
	request.Header.Add(X_AUTH_KEY_HEADER, authHeaders.XAuthKey)

	response, error := client.Do(request)

	if error != nil {
		return brokerapi.Binding{}, error
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(response.Body)
	var zoneCreateClientResponse ZoneCreateClientResponse

	if error := json.Unmarshal(buffer.Bytes(), &zoneCreateClientResponse); error != nil {
		m.logger.Error("Error decoding details.RawParameters", error)
		return brokerapi.Binding{}, error
	}

	if zoneCreateClientResponse.Success == false {
		m.logger.Error("Error from CloudFlare Client API", errors.New(fmt.Sprintf("%+v", zoneCreateClientResponse)))
		return brokerapi.Binding{}, errors.New(fmt.Sprintf("%+v", zoneCreateClientResponse))
	}

	return brokerapi.Binding{}, nil
}

type ZoneCreateClientResponse struct {
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   struct {
			 ID          string   `json:"id"`
			 Name        string   `json:"name"`
			 NameServers []string `json:"name_servers"`
		 } `json:"result"`
	Success bool `json:"success"`
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
	logger := lager.NewLogger("my-service-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))
	serviceBroker := &myServiceBroker{
		logger: logger,
	}
	credentials := brokerapi.BrokerCredentials{
		Username: os.Getenv(BROKER_USERNAME),
		Password: os.Getenv(BROKER_PASSWORD),
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)
	http.ListenAndServe(":"+os.Getenv(BROKER_PORT), nil)
}
