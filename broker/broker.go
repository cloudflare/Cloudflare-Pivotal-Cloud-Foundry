package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"code.cloudfoundry.org/lager"
	"github.com/cloudflare/Cloudflare-Pivotal-Cloud-Foundry/api"
	"github.com/pivotal-cf/brokerapi"
)

const BROKER_PASSWORD = "BROKER_PASSWORD"
const BROKER_PORT = "BROKER_PORT"
const BROKER_USERNAME = "BROKER_USERNAME"

const ENDPOINT_NOT_AVAILABLE = "This endpoint is not available"

type CloudflareBroker struct {
	logger        lager.Logger
	Zones         map[string]Zone
	CloudflareAPI api.CloudflareAPIInterface
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

func getZoneKey(instanceID string, bindingID string) string {
	return instanceID + ":" + bindingID
}

func (*CloudflareBroker) Services(context context.Context) []brokerapi.Service {
	return []brokerapi.Service{
		{
			ID:            "31e38e96-df7e-4a38-b3cb-f489fc8ab421",
			Name:          "Cloudflare Performance & Security",
			Description:   "Cloud-based performance and security solution for websites and applications.",
			Bindable:      true,
			PlanUpdatable: false,
			// Tags are "Attributes or names of backing technologies behind the service"
			// http://cloud.spring.io/spring-cloud-connectors/spring-cloud-cloud-foundry-connector.html
			Tags: []string{"Cloudflare"},
			Plans: []brokerapi.ServicePlan{
				{
					ID:          "e5c2ef96-fda2-417a-92af-dee310081600",
					Name:        "cloudflare-free",
					Description: "Cloudflare delivers performance, security, reliability and insights for all websites and applications that join the network. Once your website is on Cloudflare, all traffic will be routed through their intelligent global network of 100+ data centers. Cloudflare’s platform includes a myriad of security features, including DDoS attack mitigation and a web application firewall (WAF) for paid plans.",
					Metadata: &brokerapi.ServicePlanMetadata{
						DisplayName: "Free",
						Bullets: []string{
							"Limited DDoS protection",
							"No bandwidth limits",
							"Global CDN",
							"Shared SSL certificate",
							"I'm Under Attack™ mode",
							"3 Page Rules included",
						},
					},
				},
			},
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "Cloudflare",
				ImageUrl:            "https://www.cloudflare.com/img/logo-cloudflare.svg",
				SupportUrl:          "support.cloudflare.com",
				DocumentationUrl:    "https://github.com/cloudflare/Cloudflare-Pivotal-Cloud-Foundry",
				ProviderDisplayName: "Cloudflare",
				LongDescription:     "Cloudflare delivers performance, security, reliability and insights for all websites and applications that join the network. Once your website is on Cloudflare, all traffic will be routed through their intelligent global network of 100+ data centers. Cloudflare’s platform includes a myriad of security features, including DDoS attack mitigation and a web application firewall (WAF) for paid plans.",
			},
		},
	}
}

func (b *CloudflareBroker) Provision(context context.Context, instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	var authHeaders api.AuthHeaders

	if err := json.Unmarshal(details.RawParameters, &authHeaders); err != nil {
		b.logger.Error("Error decoding details.RawParameters", err)
		return brokerapi.ProvisionedServiceSpec{}, err
	}

	if authHeaders == (api.AuthHeaders{}) {
		return brokerapi.ProvisionedServiceSpec{}, errors.New("Error: Auth parameters are empty")
	}

	// Maybe check if already provision and do not allow unless deprovisioned?
	b.CloudflareAPI.SetAuthHeaders(authHeaders)

	return brokerapi.ProvisionedServiceSpec{}, nil
}

func (b *CloudflareBroker) Deprovision(context context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	// Clear data
	b.CloudflareAPI = &api.CloudflareAPI{}

	return brokerapi.DeprovisionServiceSpec{}, nil
}

func (b *CloudflareBroker) Bind(context context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	paramDomain, ok := details.Parameters["domain"]
	if !ok {
		return brokerapi.Binding{}, errors.New("key 'domain' not found in BindDetails.Parameters.")
	}

	domain, ok := paramDomain.(string)
	if !ok {
		return brokerapi.Binding{}, errors.New("key 'domain' is not a in type string.")
	}

	data, err := b.CloudflareAPI.AddZone(domain)
	if err != nil {
		b.logger.Error("Bind calling api.cloudflare", err)
		return brokerapi.Binding{}, err
	}

	var zoneCreateResponse ZoneCreateResponse
	if err := json.Unmarshal(data, &zoneCreateResponse); err != nil {
		b.logger.Error("Error decoding details.RawParameters", err)
		return brokerapi.Binding{}, err
	}

	if zoneCreateResponse.Success == false {
		b.logger.Error("Error from CloudFlare Client API", errors.New(fmt.Sprintf("%+v", zoneCreateResponse)))
		return brokerapi.Binding{}, errors.New(fmt.Sprintf("%+v", zoneCreateResponse))
	}

	zoneKey := getZoneKey(instanceID, bindingID)
	b.Zones[zoneKey] = zoneCreateResponse.Result

	return brokerapi.Binding{
		Credentials: zoneCreateResponse.Result,
	}, nil
}

func (b *CloudflareBroker) Unbind(context context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	zoneKey := getZoneKey(instanceID, bindingID)
	zone, ok := b.Zones[zoneKey]
	if !ok {
		return errors.New("Zone does not exist")
	}

	// Delete Zone from Cloudflare
	err := b.CloudflareAPI.DeleteZone(zone.ID)
	if err != nil {
		b.logger.Error("Unbind calling api.cloudflare", err)
		return err
	}

	// Remove from local Zone List
	delete(b.Zones, zoneKey)

	return nil
}

func (b *CloudflareBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	b.logger.Debug("LastOperation", lager.Data{"log_message": ENDPOINT_NOT_AVAILABLE})

	return brokerapi.LastOperation{State: "failed", Description: ENDPOINT_NOT_AVAILABLE}, nil
}

func (b *CloudflareBroker) Update(context context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	b.logger.Debug("Update", lager.Data{"log_message": ENDPOINT_NOT_AVAILABLE})

	return brokerapi.UpdateServiceSpec{OperationData: ENDPOINT_NOT_AVAILABLE}, nil
}

func New(logger lager.Logger, zones map[string]Zone) CloudflareBroker {
	cloudflareAPI := &api.CloudflareAPI{}

	return CloudflareBroker{
		Zones:         zones,
		CloudflareAPI: cloudflareAPI,
		logger:        logger,
	}
}
