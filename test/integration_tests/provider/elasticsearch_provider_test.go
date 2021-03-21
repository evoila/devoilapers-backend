package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateElasticSearchTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	var esProvider provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
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
	esProviderPtr, _ := CreateElasticSearchTestProvider(t)
	esProvider := *esProviderPtr

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
	formData1 := formDataObj1.(dtos.ServiceCreationFormDto)
	formData2 := formDataObj2.(dtos.ServiceCreationFormDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values
	filledForm := dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"
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
	esProviderPtr, config := CreateElasticSearchTestProvider(t)
	esProvider := *esProviderPtr

	user := config.Users[0]
	filledForm := dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "es-test-cluster"

	service1Ptr := common_test.CommonProviderStart(t, esProviderPtr, user, filledForm, 1)
	service1 := *service1Ptr

	// Check whether service is an Elasticsearch service
	service1es, ok := service1.(elasticsearch.ElasticSearchService)
	assert.True(t, ok)

	secret, _ := service1es.K8sApi.GetSecret(user.KubernetesNamespace, service1es.GetName()+"-es-http-certs-internal")

	// Test set certificate to service
	certDto := &dtos.CertificateDto{
		CaCrt:  base64.StdEncoding.EncodeToString(secret.Data["ca.crt"]),
		TlsCrt: base64.StdEncoding.EncodeToString(secret.Data["tls.crt"]),
		TlsKey: base64.StdEncoding.EncodeToString(secret.Data["tls.key"]),
	}

	_, err := service1es.SetCertificateToService(certDto)
	assert.Nil(t, err)

	// Check status of service after setting the certificate
	service3ptr, err := common_test.WaitForServiceComeUp(esProvider, user, service1.GetName())
	assert.Nil(t, err)
	service3 := *service3ptr
	assert.NotNil(t, service3)
	assert.True(t, service.ServiceStatusOk == service3.GetStatus())

	// Delete everything
	common_test.CommonProviderStop(t, esProviderPtr, user)

	// Check whether the secret with associated certificate is also deleted
	secret, err = service1es.K8sApi.GetSecret(user.KubernetesNamespace, service1es.GetName()+"-tls-cert")
	assert.NotNil(t, err)
}
