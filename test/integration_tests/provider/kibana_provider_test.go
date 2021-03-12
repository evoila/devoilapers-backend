package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	esDtos "OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/pkg/kibana"
	"OperatorAutomation/pkg/kibana/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	yamlSerializer "gopkg.in/yaml.v2"
	"testing"
	"time"
)

func CreateKibanaTestProvider(t *testing.T) (*kibana.KibanaProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	kbProvider := kibana.CreateKibanaProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath)

	return &kbProvider, config
}

func Test_Kibana_Provider_Create_Panic_Template_Not_Found(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Kibana provider did not panic if the template could not be found")
		}
	}()

	kbProvider := kibana.CreateKibanaProvider(
		"Server",
		"CaPath",
		"NotExistingPath")

	_ = kbProvider
}

func Test_Kibana_Provider_GetAttributes(t *testing.T) {
	kbProvider, _ := CreateKibanaTestProvider(t)

	var dummyEsProvider provider.IServiceProvider = unit_test.TestProvider{
		GetServiceTypeCb: func() string {
			return "Elasticsearch"
		},
		GetServicesCb: func(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
			return []*service.IService{}, nil
		},
	}

	kbProvider.OnCoreInitialized([]*provider.IServiceProvider{
		&dummyEsProvider,
	})

	assert.NotEqual(t, "", kbProvider.GetServiceImage())
	assert.NotEqual(t, "", kbProvider.GetServiceDescription())
	assert.Equal(t, "Kibana", kbProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "MyNamespace",
	}

	// Get json form data
	formDataObj1, err := kbProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj1)
	formDataObj2, err := kbProvider.GetJsonForm(testUser)
	assert.Nil(t, err)
	assert.NotNil(t, formDataObj2)

	// Ensure they are not the same (because of the random name)
	formData1 := formDataObj1.(dtos.FormQueryDto)
	formData2 := formDataObj2.(dtos.FormQueryDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)


	// Generate yaml from form values and ensure it sets the values from form
	filledForm := dtos.FormResponseDto{Common: dtos.FormResponseDtoCommon{
		ClusterName: "MyCluster",
		ElasticSearchInstance: "MyElasticSearchInstance"},
	}

	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := kbProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	kibanaYaml := yamlTemplate.(dtos.ProviderYamlTemplateDto)
	assert.Equal(t, "MyCluster", kibanaYaml.Metadata.Name)
	assert.Equal(t, "MyNamespace", kibanaYaml.Metadata.Namespace)
	assert.Equal(t, "MyElasticSearchInstance", kibanaYaml.Spec.ElasticsearchRef.Name)
}

func Test_Kibana_Provider_End2End(t *testing.T) {
	kbProvider, config := CreateKibanaTestProvider(t)

	user := config.Users[0]
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}

	// Kibana depends on elastic search therefore we need to create it
	var esProvider provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath,
	)

	// Since kibana provider utilizes the es provider we have to give a hint of other existing providers
	kbProvider.OnCoreInitialized([]*provider.IServiceProvider{&esProvider})

	// Create a new es instance
	esFormResponseDto := esDtos.FormResponseDto{Common: esDtos.FormResponseDtoCommon{ClusterName: "kibana-es-test"}}
	esFormResponseDtoBytes, err := json.Marshal(esFormResponseDto)
	assert.Nil(t, err)
	esYamlObj, err := esProvider.GetYamlTemplate(user, esFormResponseDtoBytes)
	assert.Nil(t, err)
	esYaml, err := yamlSerializer.Marshal(esYamlObj)
	assert.Nil(t, err)
	err = esProvider.CreateService(user, string(esYaml))
	assert.Nil(t, err)
	esServices, err := esProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(esServices))
	esService0 := *esServices[0]
	time.Sleep(15 * time.Second)

	// Continue with actual kb provider
	// Generate a form response that would arrive from the frontent
	filledForm := dtos.FormResponseDto{Common: dtos.FormResponseDtoCommon{
		ClusterName: "kibana-test",
		// Reference the elastic search instance
		ElasticSearchInstance: esFormResponseDto.Common.ClusterName,
	}}

	jsonFilledForm, err := json.Marshal(filledForm)
	assert.Nil(t, err)

	// Generate the yaml based on the form value
	yamlObj, err := kbProvider.GetYamlTemplate(user, jsonFilledForm)
	assert.Nil(t, err)
	yamlBytes, err := yamlSerializer.Marshal(yamlObj)
	assert.Nil(t, err)
	yaml := string(yamlBytes)

	// Check if there is no other service
	services, err := kbProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(services))

	// Try create a service with invalid yaml
	err = kbProvider.CreateService(user, "something")
	assert.NotNil(t, err)

	// Try create a service with invalid user data
	err = kbProvider.CreateService(invalidUser, yaml)
	assert.NotNil(t, err)

	// Create a service
	err = kbProvider.CreateService(user, yaml)
	assert.Nil(t, err)

	// Try check if created with invalid user data
	services, err = kbProvider.GetServices(invalidUser)
	assert.NotNil(t, err)

	// Check if created
	services, err = kbProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	service0 := *services[0]
	assert.NotEqual(t, "", service0.GetName())
	assert.Equal(t, kbProvider.GetServiceType(), service0.GetType())
	assert.Equal(t, 0, len(service0.GetActions()))

	// Try get service with invalid user data
	_, err = kbProvider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok. Kibana needs some extra time.
	var service1 service.IService
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Second)

		// Try get service with invalid user data
		service1Ptr, err := kbProvider.GetService(user, service0.GetName())
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
	err = kbProvider.DeleteService(user, "some-not-existing-id")
	assert.NotNil(t, err)

	// Try delete service with invalid user
	err = kbProvider.DeleteService(invalidUser, (*services[0]).GetName())
	assert.NotNil(t, err)

	// Delete service
	err = kbProvider.DeleteService(user, (*services[0]).GetName())
	assert.Nil(t, err)

	// Delete es instance
	err = esProvider.DeleteService(user, esService0.GetName())
	assert.Nil(t, err)
}
