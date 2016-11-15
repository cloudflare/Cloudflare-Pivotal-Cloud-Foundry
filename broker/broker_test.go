package broker_test

import (
	"context"
	"testing"

	"code.cloudfoundry.org/lager"

	"github.com/jwineman/Cloudflare-Pivotal-Cloud-Foundry/broker"
	"github.com/pivotal-cf/brokerapi"
)

func TestNewWithEmptyZone(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	zones := map[string]broker.Zone{}
	cloudflarebroker := broker.New(logger, zones)
	if cloudflarebroker.Zones == nil {
		t.Errorf("TestNewWithEmptyZone failed")
	}
}

func TestNewWithValidZone(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")

	zone := broker.Zone{ID: "id", Name: "name", NameServers: nil}
	zones := map[string]broker.Zone{
		"testZone": zone,
	}
	cloudflarebroker := broker.New(logger, zones)
	if len(cloudflarebroker.Zones) != 1 ||
		cloudflarebroker.Zones["testZone"].ID != zone.ID ||
		cloudflarebroker.Zones["testZone"].Name != zone.Name ||
		cloudflarebroker.Zones["testZone"].NameServers != nil {
		t.Errorf("TestNewWithValidZone failed")
	}
}

func TestService(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context

	service := cloudflarebroker.Services(context)
	if service[0].Name != "cloudflare" {
		t.Errorf("Service failed")
	}
}

func TestProvisionEmpty(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"

	_, err := cloudflarebroker.Provision(
		context,
		instanceId,
		brokerapi.ProvisionDetails{
			ServiceID:        "org-guid-here",
			PlanID:           "plan-guid-here",
			OrganizationGUID: "service-guid-here",
			SpaceGUID:        "space-guid-here",
			RawParameters:    nil,
		},
		false,
	)
	if err == nil {
		t.Errorf("Provision failed with RawParameters nil")
	}

	_, err = cloudflarebroker.Provision(
		context,
		instanceId,
		brokerapi.ProvisionDetails{
			ServiceID:        "org-guid-here",
			PlanID:           "plan-guid-here",
			OrganizationGUID: "service-guid-here",
			SpaceGUID:        "space-guid-here",
			RawParameters:    []byte(`{}`),
		},
		false,
	)
	if err == nil {
		t.Errorf("Provision failed with RawParameters empty json")
	}
}

func TestProvisionWithCredentials(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"

	_, err := cloudflarebroker.Provision(
		context,
		instanceId,
		brokerapi.ProvisionDetails{
			ServiceID:        "org-guid-here",
			PlanID:           "plan-guid-here",
			OrganizationGUID: "service-guid-here",
			SpaceGUID:        "space-guid-here",
			RawParameters: []byte(
				`{
					"x-auth-key":   "mykey",
					"x-auth-email": "email@email.com"
				}`),
		},
		false,
	)
	if err != nil {
		t.Errorf("Provision failed with correct credentials")
	}
}

func TestDeprovision(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"

	_, err := cloudflarebroker.Deprovision(
		context,
		instanceId,
		brokerapi.DeprovisionDetails{},
		false,
	)
	if err != nil {
		t.Errorf("Deprovision failed")
	}
}

func TestBindWithEmptyParameters(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"
	bindingId := "2"

	_, err := cloudflarebroker.Bind(
		context,
		instanceId,
		bindingId,
		brokerapi.BindDetails{},
	)
	if err == nil {
		t.Errorf("Bind failed with empty details")
	}
}

func TestBindWithFalseParameters(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"
	bindingId := "2"

	params := map[string]interface{}{
		"domain": 1,
	}
	_, err := cloudflarebroker.Bind(
		context,
		instanceId,
		bindingId,
		brokerapi.BindDetails{Parameters: params},
	)
	if err == nil {
		t.Errorf("Bind failed with false details")
	}
}

func TestBindWithCorrectParameters(t *testing.T) {
	// logger := lager.NewLogger("cloudflare-broker")
	// cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	// var context context.Context
	// instanceId := "1"
	// bindingId := "2"

	// params := map[string]interface{}{
	// 	"domain": "domain.com",
	// }
	// _, err := cloudflarebroker.Bind(
	// 	context,
	// 	instanceId,
	// 	bindingId,
	// 	brokerapi.BindDetails{Parameters: params},
	// )

	// if err != nil {
	// 	t.Errorf("Bind failed with false details")
	// }

	// TODO check binding != (brokerapi.Binding{})
}

func TestUnbind(t *testing.T) {
	// TODO do this after doing bind
}

func TestLastOperation(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"

	operation, err := cloudflarebroker.LastOperation(context, instanceId, "")
	if err != nil ||
		operation.State != "failed" ||
		operation.Description != broker.ENDPOINT_NOT_AVAILABLE {
		t.Errorf("LastOperation failed")
	}
}

func TestUpdate(t *testing.T) {
	logger := lager.NewLogger("cloudflare-broker")
	cloudflarebroker := broker.New(logger, map[string]broker.Zone{})
	var context context.Context
	instanceId := "1"

	operation, err := cloudflarebroker.Update(context, instanceId, brokerapi.UpdateDetails{}, false)
	if err != nil || operation.OperationData != broker.ENDPOINT_NOT_AVAILABLE {
		t.Errorf("LastOperation failed")
	}
}

// TODO getAuthHeaders()
// TODO setAuthHeaders
