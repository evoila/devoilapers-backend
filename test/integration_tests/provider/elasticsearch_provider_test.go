package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	"OperatorAutomation/pkg/elasticsearch/dtos/provider_dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func CreateElasticSearchTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	url, err := url.Parse(config.Kubernetes.Server)
	assert.Nil(t, err)

	var esProvider provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		url.Hostname(),
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
		"Hostname",
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
	formData1 := formDataObj1.(provider_dtos.ServiceCreationFormDto)
	formData2 := formDataObj2.(provider_dtos.ServiceCreationFormDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values
	filledForm := provider_dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"
	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := esProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	elasticSearchYaml := yamlTemplate.(provider_dtos.ProviderYamlTemplateDto)
	assert.Equal(t, "MyCluster", elasticSearchYaml.Metadata.Name)
	assert.Equal(t, "MyNamespace", elasticSearchYaml.Metadata.Namespace)
}

func Test_Elasticsearch_Provider_End2End(t *testing.T) {
	esProviderPtr, config := CreateElasticSearchTestProvider(t)
	esProvider := *esProviderPtr

	user := config.Users[0]
	filledForm := provider_dtos.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "es-test-cluster"

	service1Ptr := common_test.CommonProviderStart(t, esProviderPtr, user, filledForm, 3)
	service1 := *service1Ptr

	// Actions
	// --- Scale ---
	// Scale = 1
	// Check
	actionPtr, err := common_test.GetAction(service1Ptr, "Features", "cmd_es_scale")
	assert.Nil(t, err)
	action := *actionPtr
	placeholder := action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	assert.Equal(t, 1, placeholder.NumberOfReplicas)

	// Scale up -> 2
	actionPtr, err = common_test.GetAction(service1Ptr, "Features", "cmd_es_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	placeholder.NumberOfReplicas = 2
	result, err := action.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)
	// Check
	tempServicePtr, err := esProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_es_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	assert.Equal(t, 2, placeholder.NumberOfReplicas)

	// Scale down -> 1
	tempServicePtr, err = esProvider.GetService(user, service1.GetName())
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_es_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	placeholder.NumberOfReplicas = 1
	result, err = action.GetActionExecuteCallback()(placeholder)
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)

	// Check
	tempServicePtr, err = esProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_es_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	assert.Equal(t, 1, placeholder.NumberOfReplicas)


	// --- Exposure ---
	// Check if toggle is correct
	service1Ptr, err = esProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	service1 = *service1Ptr
	toggleActionPtr, err := common_test.GetToggleAction(service1Ptr, "Security", "cmd_es_expose_toggle")
	assert.Nil(t, err)
	toggleAction := *toggleActionPtr
	isSet, err := toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_es_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation := result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)

	// Expose it
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_es_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Set()
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)

	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_es_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.True(t, isSet) // exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_es_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)

	// Hide it again
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_es_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Unset()
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// Check again if it is hidden
	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_es_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_es_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)
	
	
	
	
	// Check whether service is an Elasticsearch service
	service1es, ok := service1.(elasticsearch.ElasticSearchService)
	assert.True(t, ok)

	secret, _ := service1es.K8sApi.GetSecret(user.KubernetesNamespace, service1es.GetName()+"-es-http-certs-internal")

	// Set cert
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_es_set_cert_action")
	assert.Nil(t, err)
	action = *actionPtr
	certDto := action.GetJsonFormResultPlaceholder().(*action_dtos.CertificateDto)
	certDto.CaCrt = base64.StdEncoding.EncodeToString(secret.Data["ca.crt"])
	certDto.TlsCrt = base64.StdEncoding.EncodeToString(secret.Data["tls.crt"])
	certDto.TlsKey = base64.StdEncoding.EncodeToString(secret.Data["tls.key"])
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	assert.Nil(t, result)

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
