package common_test

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func CommonProviderStart(t *testing.T, providerPtr *provider.IServiceProvider, user common.IKubernetesAuthInformation, creationForm interface{}, expectedNumberOfActionGroups int) *service.IService {
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}
	provider := *providerPtr

	filledFormBytes, err := json.Marshal(creationForm)
	assert.Nil(t, err)

	yamlObj, err := provider.GetYamlTemplate(user, filledFormBytes)

	yamlBytes, err := yaml.Marshal(yamlObj)
	assert.Nil(t, err)
	yaml := string(yamlBytes)
	assert.True(t, len(yaml) > 10)

	// Check if there is no other service
	services, err := provider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(services))

	// Try create a service with invalid yaml
	err = provider.CreateService(user, "something")
	assert.NotNil(t, err)

	// Try create a service with invalid user data
	err = provider.CreateService(invalidUser, yaml)
	assert.NotNil(t, err)

	// Create a service
	err = provider.CreateService(user, yaml)
	assert.Nil(t, err)

	// Try check if created with invalid user data
	services, err = provider.GetServices(invalidUser)
	assert.NotNil(t, err)

	// Check if created
	services, err = provider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	service0 := *services[0]
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, provider.GetServiceType(), service0.GetType())
	assert.Equal(t, expectedNumberOfActionGroups, len(service0.GetActionGroups()))

	// Try get service with invalid user data
	_, err = provider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok
	service1Ptr, err := WaitForServiceComeUp(provider, user, service0.GetName())
	assert.Nil(t, err)
	service1 := *service1Ptr

	// Ensure service is ok
	assert.Equal(t, service.ServiceStatusOk, service1.GetStatus())

	// Ensure they have the same attributes
	assert.Equal(t, service0.GetName(), service1.GetName())
	assert.Equal(t, service0.GetType(), service1.GetType())

	return service1Ptr
}

func CommonProviderStop(t *testing.T, providerPtr *provider.IServiceProvider, user common.IKubernetesAuthInformation)  {
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}
	provider := *providerPtr

	// Try delete service with invalid id
	err := provider.DeleteService(user, "some-not-existing-id")
	assert.NotNil(t, err)

	services, err := provider.GetServices(user)
	assert.Nil(t, err)
	assert.True(t, len(services) > 0)

	// Try delete service with invalid user
	err = provider.DeleteService(invalidUser, (*services[0]).GetName())
	assert.NotNil(t, err)

	// Delete service
	for _, service := range services {
		err = provider.DeleteService(user, (*service).GetName())
		assert.Nil(t, err)
	}
}