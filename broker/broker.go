package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
)

const BROKER_PASSWORD = "BROKER_PASSWORD"
const BROKER_PORT = "BROKER_PORT"
const BROKER_USERNAME = "BROKER_USERNAME"
const CLOUDFLARE_EMAIL = "CLOUDFLARE_EMAIL"
const CLOUDFLARE_API_KEY = "CLOUDFLARE_API_KEY"
const X_AUTH_EMAIL_HEADER = "X-Auth-Email"
const X_AUTH_KEY_HEADER = "X-Auth-Key"

const ENDPOINT_NOT_AVAILABLE = "This endpoint is not available"

type CloudflareBroker struct {
	logger lager.Logger
	Zones  map[string]Zone
	Auth   AuthHeaders
}

type AuthHeaders struct {
	XAuthEmail string `json:"x-auth-email"`
	XAuthKey   string `json:"x-auth-key"`
}

type Zone struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	NameServers []string `json:"name_servers"`
}

type ZoneCreateResponse struct {
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   Zone          `json:"result"`
	Success  bool          `json:"success"`
}

func (b *CloudflareBroker) getAuthHeaders() AuthHeaders {
	return AuthHeaders{
		XAuthEmail: os.Getenv(CLOUDFLARE_EMAIL),
		XAuthKey:   os.Getenv(CLOUDFLARE_API_KEY),
	}
}

func (b *CloudflareBroker) setAuthHeaders(authHeaders AuthHeaders) AuthHeaders {
	os.Setenv(CLOUDFLARE_EMAIL, authHeaders.XAuthEmail)
	os.Setenv(CLOUDFLARE_API_KEY, authHeaders.XAuthKey)

	return b.getAuthHeaders()
}

func (*CloudflareBroker) Services(context context.Context) []brokerapi.Service {
	return []brokerapi.Service{
		{
			ID:            "31e38e96-df7e-4a38-b3cb-f489fc8ab421",
			Name:          "Cloudflare",
			Description:   "Give us five minutes and we’ll supercharge your website.",
			Bindable:      true,
			PlanUpdatable: false,
			Tags:          []string{"Cloudflare", "HTTP2", "SSL", "TLS", "CDN"}, //TODO
			Plans: []brokerapi.ServicePlan{
				{
					ID:          "e5c2ef96-fda2-417a-92af-dee310081600",
					Name:        "cloudflare-free",
					Description: "Fast site performance. Broad security protection. SSL. Powerful stats about your visitors. Peace of mind about running your website so you can get back to what you love.",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Cloudflare Free Plan",
					},
				},
			},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "Cloudflare",
				ImageUrl:            "TODO image url",
				SupportUrl:          "TODO put github issue link",
				DocumentationUrl:    "TODO put github README link (maybe KB article)",
				ProviderDisplayName: "Cloudflare Inc.",
				LongDescription:     "Fast site performance. Broad security protection. SSL. Powerful stats about your visitors. Peace of mind about running your website so you can get back to what you love.",
			},
		},
	}
}

func (b *CloudflareBroker) Provision(context context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	var authHeaders AuthHeaders

	if err := json.Unmarshal(details.RawParameters, &authHeaders); err != nil {
		// b.logger.Error("Error decoding details.RawParameters", err)
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	if authHeaders == (AuthHeaders{}) {
		return brokerapi.ProvisionedServiceSpec{}, errors.New("Error: Auth parameters are empty")
	}

	b.setAuthHeaders(authHeaders)

	return brokerapi.ProvisionedServiceSpec{}, nil
}

func (b *CloudflareBroker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	// Clear environment variables
	envVariables := [5]string{BROKER_PASSWORD, BROKER_PORT, BROKER_USERNAME, CLOUDFLARE_EMAIL, CLOUDFLARE_API_KEY}
	for _, envVariable := range envVariables {
		if err := os.Unsetenv(envVariable); err != nil {
			return brokerapi.DeprovisionServiceSpec{}, errors.New("Failed Devprovisioning")
		}
	}

	// Clear set data
	b.Zones = map[string]Zone{}

	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *CloudflareBroker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	// Bind to instances here
	// Return a binding which contains a credentials object that can be marshalled to JSON,
	// and (optionally) a syslog drain URL.
	paramDomain, ok := details.Parameters["domain"]
	if !ok {
		return brokerapi.Binding{}, errors.New("key 'domain' not found in BindDetails.Parameters.")
	}

	domain, ok := paramDomain.(string)
	if !ok {
		return brokerapi.Binding{}, errors.New("key 'domain' is not a in type string.")
	}
	var jsonBody = []byte(`{"name" : "` + domain + `"}`)
	request, err := http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones", bytes.NewReader(jsonBody))
	if err != nil {
		return brokerapi.Binding{}, errors.New("Failed API Call. Please check your connection")
	}

	var authHeaders = b.getAuthHeaders()
	request.Header.Add(X_AUTH_EMAIL_HEADER, authHeaders.XAuthEmail)
	request.Header.Add(X_AUTH_KEY_HEADER, authHeaders.XAuthKey)

	var client = http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return brokerapi.Binding{}, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return brokerapi.Binding{}, err
	}

	var ZoneCreateResponse ZoneCreateResponse
	if err := json.Unmarshal(data, &ZoneCreateResponse); err != nil {
		b.logger.Error("Error decoding details.RawParameters", err)
		return brokerapi.Binding{}, err
	}

	if ZoneCreateResponse.Success == false {
		// TODO: Error should be updated to ZoneCreateResponse.error
		b.logger.Error("Error from CloudFlare Client API", errors.New(fmt.Sprintf("%+v", ZoneCreateResponse)))
		return brokerapi.Binding{}, errors.New(fmt.Sprintf("%+v", ZoneCreateResponse))
	}

	// TODO: It can be changed to something like "instanceID:bindingID"
	b.Zones[bindingID] = ZoneCreateResponse.Result

	return brokerapi.Binding{
		Credentials: ZoneCreateResponse.Result,
	}, nil
}

func (b *CloudflareBroker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	zone, ok := b.Zones[bindingID]
	if !ok {
		return errors.New("Zone does not exist")
	}

	// Delete Zone from Cloudflare
	requestURL := "https://api.cloudflare.com/client/v4/zones/" + zone.ID
	request, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	var authHeaders = b.getAuthHeaders()
	request.Header.Add(X_AUTH_EMAIL_HEADER, authHeaders.XAuthEmail)
	request.Header.Add(X_AUTH_KEY_HEADER, authHeaders.XAuthKey)

	var client = http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	// Remove from local Zone List
	delete(b.Zones, bindingID)

	return nil
}

func (b *CloudflareBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	b.logger.Info("LastOperation", lager.Data{"log_message": ENDPOINT_NOT_AVAILABLE})

	return brokerapi.LastOperation{State: "failed", Description: ENDPOINT_NOT_AVAILABLE}, nil
}

func (b *CloudflareBroker) Update(context context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	b.logger.Info("Update", lager.Data{"log_message": ENDPOINT_NOT_AVAILABLE})

	return brokerapi.UpdateServiceSpec{OperationData: ENDPOINT_NOT_AVAILABLE}, nil
}

func New(logger lager.Logger, zones map[string]Zone) CloudflareBroker {
	return CloudflareBroker{
		Zones:  zones,
		logger: logger,
	}
}
