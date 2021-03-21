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
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	yamlSerializer "gopkg.in/yaml.v2"
	"strings"
	"testing"
	"time"
)

func CreateKibanaTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	var kbProvider provider.IServiceProvider = kibana.CreateKibanaProvider(
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
	kbProviderPtr, _ := CreateKibanaTestProvider(t)
	kbProvider := *kbProviderPtr

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
	formData1 := formDataObj1.(dtos.ServiceCreationFormDto)
	formData2 := formDataObj2.(dtos.ServiceCreationFormDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values and ensure it sets the values from form
	filledForm := dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"
	filledForm.Common.ElasticSearchInstance = "MyElasticSearchInstance"

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
	kbProviderPtr, config := CreateKibanaTestProvider(t)
	kbProvider := *kbProviderPtr
	user := config.Users[0]

	// Kibana depends on elastic search therefore we need to create it
	var esProvider provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath,
	)

	// Since kibana provider utilizes the es provider we have to give a hint of other existing providers
	kbProvider.OnCoreInitialized([]*provider.IServiceProvider{&esProvider})

	// Create a new es instance
	esFormResponseDto := esDtos.ServiceCreationFormResponseDto{}
	esFormResponseDto.Common.ClusterName = "kibana-es-test"
	esFormResponseDtoBytes, err := json.Marshal(esFormResponseDto)
	assert.Nil(t, err)
	esYamlObj, err := esProvider.GetYamlTemplate(user, esFormResponseDtoBytes)
	assert.Nil(t, err)
	esYaml, err := yamlSerializer.Marshal(esYamlObj)
	assert.Nil(t, err)
	err = esProvider.CreateService(user, string(esYaml))
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)
	esServices, err := esProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(esServices))

	// Check if kibana form offers elasticsearch
	form, err := kbProvider.GetJsonForm(user)
	assert.Nil(t, err)
	formJsonBytes, err := json.Marshal(form)
	assert.True(t, strings.Contains(string(formJsonBytes), "kibana-es-test"))

	// Continue with actual kb provider
	// Generate a form response that would arrive from the frontent
	filledForm := dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "kibana-test"
	filledForm.Common.ElasticSearchInstance = esFormResponseDto.Common.ClusterName

	service1Ptr := common_test.CommonProviderStart(t, kbProviderPtr, user, filledForm, 1)
	service1 := *service1Ptr

	// Check whether service is an Kibana service
	service2, ok := service1.(kibana.KibanaService)
	assert.True(t, ok)

	secret, _ := service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-kb-http-certs-internal")

	// Test set certificate to service
	certDto := &dtos.CertificateDto{
		CaCrt:  base64.StdEncoding.EncodeToString(secret.Data["ca.crt"]),
		TlsCrt: base64.StdEncoding.EncodeToString(secret.Data["tls.crt"]),
		TlsKey: base64.StdEncoding.EncodeToString(secret.Data["tls.key"]),
	}

	_, err = service2.SetCertificateToService(certDto)
	assert.Nil(t, err)

	// Check status of service after setting the certificate
	service3Ptr, err := common_test.WaitForServiceComeUp(kbProvider, user, service1.GetName())
	assert.Nil(t, err)
	service3 := *service3Ptr

	// Check whether status of service is ok after setting the certificate
	assert.True(t, service.ServiceStatusOk == service3.GetStatus())

	common_test.CommonProviderStop(t, kbProviderPtr, user)

	// Wait till delete service is done
	time.Sleep(5 * time.Second)
	// Check whether the secret with associated certificate is also deleted
	secret, err = service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-tls-cert")
	assert.NotNil(t, err)
}
