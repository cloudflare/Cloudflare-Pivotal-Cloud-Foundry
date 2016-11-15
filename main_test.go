package main_test

// import (
// 	"strconv"
// 	"testing"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"

// 	"github.com/pivotal-cf/brokerapi"
// )

// var _ = Describe(".RedisService", func() {
// 	var (
// 		redisBroker broker.Broker
// 	)

// 	JustBeforeEach(func() {
// 		redisBroker = broker.New(broker.DatabaseIDs{}, "0.0.0.0")
// 	})

// 	Describe(".Services", func() {
// 		It("works", func() {
// 			Expect(redisBroker.Services()[0].Name).To(Equal("Shared Redis"))
// 		})
// 	})

// 	Describe(".Provision", func() {
// 		Context("if there are databases available", func() {
// 			It("should return a nil error", func() {
// 				_, err := redisBroker.Provision(
// 					"Pikachu",
// 					brokerapi.ProvisionDetails{},
// 					false,
// 				)
// 				Expect(err).To(BeNil())
// 			})
// 		})

// 		Context("if there are no databases available", func() {
// 			JustBeforeEach(func() {
// 				var fullArray broker.DatabaseIDs
// 				for i := range fullArray {
// 					fullArray[i] = strconv.Itoa(i)
// 				}
// 				redisBroker = broker.New(fullArray, "0.0.0.0")
// 			})

// 			It("should return an error if it can't provision a redis db", func() {
// 				_, err := redisBroker.Provision(
// 					"Pikachu",
// 					brokerapi.ProvisionDetails{},
// 					false,
// 				)
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})
// 	})

// 	Describe(".Deprovision", func() {
// 		Context("when a service does not exist", func() {
// 			It("should return an error", func() {
// 				_, err := redisBroker.Deprovision(
// 					"Pikachu",
// 					brokerapi.DeprovisionDetails{},
// 					false,
// 				)
// 				Expect(err).ToNot(BeNil())
// 			})
// 		})

// 		Context("given a service has been provisioned", func() {
// 			JustBeforeEach(func() {
// 				redisBroker = broker.New(broker.DatabaseIDs{"Pikachu"}, "0.0.0.0")
// 			})

// 			It("should return no error", func() {
// 				_, err := redisBroker.Deprovision(
// 					"Pikachu",
// 					brokerapi.DeprovisionDetails{},
// 					false,
// 				)
// 				Expect(err).To(BeNil())
// 			})
// 		})
// 	})

// 	Describe(".Bind", func() {
// 		Context("given a service has been provisioned", func() {
// 			JustBeforeEach(func() {
// 				redisBroker.Provision("Pikachu", brokerapi.ProvisionDetails{}, false)
// 			})

// 			It("should return no errors when sucessfully bound", func() {
// 				_, err := redisBroker.Bind("Pikachu", "test", brokerapi.BindDetails{})
// 				Expect(err).To(BeNil())
// 			})

// 			It("should return a binding when successfully bound", func() {
// 				binding, err := redisBroker.Bind("Pikachu", "test", brokerapi.BindDetails{})
// 				Expect(err).To(BeNil())
// 				if credentials, ok := binding.Credentials.(broker.Credentials); ok {
// 					Expect(credentials.Database).To(Equal(1))
// 				} else {
// 					Fail("Could not cast to Credentials object.")
// 				}
// 			})
// 		})
// 	})
// })

// var _ = Describe("Functional test", func() {
// 	It("should provision and then deprovision", func() {
// 		redisBroker := broker.New(broker.DatabaseIDs{}, "0.0.0.0")

// 		_, err := redisBroker.Provision(
// 			"Pikachu",
// 			brokerapi.ProvisionDetails{},
// 			false,
// 		)
// 		Expect(err).To(BeNil())

// 		_, err = redisBroker.Deprovision(
// 			"Pikachu",
// 			brokerapi.DeprovisionDetails{},
// 			false,
// 		)
// 		Expect(err).To(BeNil())

// 		_, err = redisBroker.Deprovision(
// 			"Pikachu",
// 			brokerapi.DeprovisionDetails{},
// 			false,
// 		)
// 		Expect(err).ToNot(BeNil())
// 	})
// })

// func TestRedisBroker(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "RedisBroker Suite")
// }

// // func TestProvision(t *testing.T) {

// // 	details := brokerapi.ProvisionDetails{
// // 		ServiceID:        "org-guid-here",
// // 		PlanID:           "plan-guid-here",
// // 		OrganizationGUID: "service-guid-here",
// // 		SpaceGUID:        "space-guid-here",
// // 		RawParameters: []byte(`{
// // 			"x-auth-key":   "mykey",
// // 			"x-auth-email": "email@email.com",
// // 		}`),
// // 	}

// // 	result, err := Provision(context.Context{}, "1", details, false)
// // 	fmt.Println(result)
// // 	fmt.Println(err)

// // 	if true != true {
// // 		t.Error("everything I know is wrong")
// // 	}
// // }
