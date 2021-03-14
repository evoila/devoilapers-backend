package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

func CreateElasticSearchTestProvider(t *testing.T) (*elasticsearch.ElasticsearchProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath)

	return &esProvider, config
}

func Test_Create_Panic_Template_Not_Found(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Elasticsearch provider did not panic if the template could not be found")
		}
	}()

	esProvider := elasticsearch.CreateElasticSearchProvider(
		"Server",
		"CaPath",
		"NotExistingPath")

	_ = esProvider
}

func Test_Elasticsearch_Provider_GetAttributes(t *testing.T) {
	esProvider, _ := CreateElasticSearchTestProvider(t)

	assert.NotEqual(t, "", esProvider.GetServiceImage())
	assert.NotEqual(t, "", esProvider.GetServiceDescription())
	assert.Equal(t, "Elasticsearch", esProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "MyNamespace",
	}

	// Get json form data
	formDataObj1, err := esProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj1)
	formDataObj2, err := esProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj2)

	// Ensure they are not the same (because of the random name)
	formData1 := formDataObj1.(dtos.FormQueryDto)
	formData2 := formDataObj2.(dtos.FormQueryDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values
	filledForm := dtos.FormResponseDto{Common: dtos.FormResponseDtoCommon{ClusterName: "MyCluster"}}
	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := esProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	elasticSearchYaml := yamlTemplate.(dtos.ProviderYamlTemplateDto)
	assert.Equal(t, "MyCluster", elasticSearchYaml.Metadata.Name)
	assert.Equal(t, "MyNamespace", elasticSearchYaml.Metadata.Namespace)
}

func Test_Elasticsearch_Provider_End2End(t *testing.T) {
	esProvider, config := CreateElasticSearchTestProvider(t)

	user := config.Users[0]
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}

	filledForm := dtos.FormResponseDto{Common: dtos.FormResponseDtoCommon{ClusterName: "my-es-test-cluster"}}
	filledFormBytes, err := json.Marshal(filledForm)
	assert.Nil(t, err)

	yamlObj, err := esProvider.GetYamlTemplate(user, filledFormBytes)

	yamlBytes, err := yaml.Marshal(yamlObj)
	assert.Nil(t, err)
	yaml := string(yamlBytes)
	assert.True(t, len(yaml) > 10)

	// Check if there is no other service
	services, err := esProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(services))

	// Try create a service with invalid yaml
	err = esProvider.CreateService(user, "something")
	assert.NotNil(t, err)

	// Try create a service with invalid user data
	err = esProvider.CreateService(invalidUser, yaml)
	assert.NotNil(t, err)

	// Create a service
	err = esProvider.CreateService(user, yaml)
	assert.Nil(t, err)

	// Try check if created with invalid user data
	services, err = esProvider.GetServices(invalidUser)
	assert.NotNil(t, err)

	// Check if created
	services, err = esProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	service0 := *services[0]
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, esProvider.GetServiceType(), service0.GetType())
	assert.Equal(t, 0, len(service0.GetActions()))
	assert.True(t,
		service.ServiceStatusPending == service0.GetStatus() ||
			service.ServiceStatusOk == service0.GetStatus(),
	)

	// Try get service with invalid user data
	_, err = esProvider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok
	var service1 service.IService
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Second)

		// Try get service with invalid user data
		service1Ptr, err := esProvider.GetService(user, service0.GetName())
		assert.Nil(t, err)
		service1 = *service1Ptr

		if service1.GetStatus() == service.ServiceStatusOk {
			break
		}
	}

	// Ensure service is ok
	assert.Equal(t, service.ServiceStatusOk, service1.GetStatus())

	// Ensure they have the same attributes
	assert.Equal(t, service0.GetName(), service1.GetName())
	assert.Equal(t, service0.GetType(), service1.GetType())

	// Try delete service with invalid id
	err = esProvider.DeleteService(user, "some-not-existing-id")
	assert.NotNil(t, err)

	// Try delete service with invalid user
	err = esProvider.DeleteService(invalidUser, (*services[0]).GetName())
	assert.NotNil(t, err)

	// Delete service
	err = esProvider.DeleteService(user, (*services[0]).GetName())
	assert.Nil(t, err)
}
