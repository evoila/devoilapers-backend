package dummy

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/dummy"
	"OperatorAutomation/test/unit_tests/common_test"
	"github.com/stretchr/testify/assert"
	yamlSerializer "gopkg.in/yaml.v2"
	"testing"
)

func Test_DummyProvider_Get_Attributes(t *testing.T) {
	dummyProvider := dummy.CreateDummyProvider()
	user := common_test.TestUser{}

	// Get json form data
	formDataObj1, err := dummyProvider.GetJsonForm(user)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj1)

	// Generate yaml
	yamlTemplate, err := dummyProvider.GetYamlTemplate(user, nil)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)
	yaml1, err := yamlSerializer.Marshal(yamlTemplate)
	assert.Nil(t, err)
	assert.True(t, len(yaml1) > 0)

	assert.Equal(t, "DummyService", dummyProvider.GetServiceType())
	assert.True(t, len(dummyProvider.GetServiceImage()) > 0)
	assert.True(t, len(dummyProvider.GetServiceDescription()) > 0)

	dummyProvider2 := dummy.CreateDummyProvider()
	yamlTemplate2, err := dummyProvider2.GetYamlTemplate(user, nil)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)
	yaml2, err := yamlSerializer.Marshal(yamlTemplate2)
	assert.Nil(t, err)

	assert.Equal(t, yaml1, yaml2)
}

func Test_DummyProvider_Create_Services(t *testing.T) {
	dummyProvider := dummy.CreateDummyProvider()
	assert.Nil(t, dummyProvider.CreateService(common_test.TestUser{}, "test: yaml"))
}

func Test_DummyProvider_Get_Services(t *testing.T) {
	auth := common_test.TestUser{}
	dummyProvider := dummy.CreateDummyProvider()

	err := dummyProvider.CreateService(auth, "1")
	assert.Nil(t, err)
	services, err := dummyProvider.GetServices(auth)
	services0 := *services[0]
	assert.Equal(t, "1", services0.GetYamlTemplate())

	err = dummyProvider.CreateService(auth, "2")
	assert.Nil(t, err)

	// Get all services
	services, err = dummyProvider.GetServices(auth)
	assert.NotNil(t, services)
	assert.Equal(t, 2, len(services))
	services1 := *services[1]
	assert.Equal(t, 1, len(services1.GetActions()))
	// Ensure values refer to other objects
	assert.NotEqual(t, services1.GetActions(), services1.GetActions())

	// Get single service
	servicePtr, err := dummyProvider.GetService(auth, services1.GetName())
	assert.Nil(t, err)
	service0 := *servicePtr
	assert.NotNil(t, service0)

	// Ensure they are two different instances
	assert.False(t, &services1 == &service0)

	// Compare values of GetServices and GetService
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, service.ServiceStatusOk, service0.GetStatus())
	assert.Equal(t, 1, len(service0.GetActions()))
	assert.NotEqual(t, "", service0.GetActions()[0].GetName())
	assert.Equal(t, 2, len(service0.GetActions()[0].GetActions()))

	assert.Equal(t, service0.GetName(), services1.GetName())
	assert.Equal(t, service0.GetStatus(), services1.GetStatus())
	assert.Equal(t, service0.GetType(), services1.GetType())
	assert.Equal(t, service0.GetYamlTemplate(), services1.GetYamlTemplate())

	_, err = dummyProvider.GetService(auth, "Not existing")
	assert.NotNil(t, err)
}

func Test_DummyProvider_Delete_Services(t *testing.T) {
	auth := common_test.TestUser{}
	dummyProvider := dummy.CreateDummyProvider()

	// Delete not existing
	err := dummyProvider.DeleteService(auth, "x")
	assert.NotNil(t, err)

	// Create 2
	assert.Nil(t, dummyProvider.CreateService(auth, "1"))
	assert.Nil(t, dummyProvider.CreateService(auth, "2"))

	// Get both
	services, err := dummyProvider.GetServices(auth)
	assert.Equal(t, 2, len(services))
	createdService0 := *services[0]

	// Delete the second one
	err = dummyProvider.DeleteService(auth, createdService0.GetName())
	assert.Nil(t, err)
	services, err = dummyProvider.GetServices(auth)
	assert.Equal(t, 1, len(services))

	// Try delete the second one again
	err = dummyProvider.DeleteService(auth, createdService0.GetName())
	assert.NotNil(t, err)

	// Get and delete remaining
	createdService1 := *services[0]
	err = dummyProvider.DeleteService(auth, createdService1.GetName())
	assert.Nil(t, err)
}
