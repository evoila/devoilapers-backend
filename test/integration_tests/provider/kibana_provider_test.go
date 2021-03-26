package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos/provider_dtos"
	"OperatorAutomation/pkg/kibana"
	provider_dtos2 "OperatorAutomation/pkg/kibana/dtos/provider_dtos"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	yamlSerializer "gopkg.in/yaml.v2"
	"net/url"
	"strings"
	"testing"
	"time"
)

func CreateKibanaTestProvider(t *testing.T) (*provider.IServiceProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	url, err := url.Parse(config.Kubernetes.Server)
	assert.Nil(t, err)

	var kbProvider provider.IServiceProvider = kibana.CreateKibanaProvider(
		url.Hostname(),
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
		"Hostname",
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
	formData1 := formDataObj1.(provider_dtos2.ServiceCreationFormDto)
	formData2 := formDataObj2.(provider_dtos2.ServiceCreationFormDto)

	assert.NotEqual(t,
		formData1.Properties.Common.Properties.ClusterName.Default,
		formData2.Properties.Common.Properties.ClusterName.Default)

	// Generate yaml from form values and ensure it sets the values from form
	filledForm := provider_dtos2.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "MyCluster"
	filledForm.Common.ElasticSearchInstance = "MyElasticSearchInstance"

	filledFormData, err := json.Marshal(filledForm)
	assert.Nil(t, err)
	yamlTemplate, err := kbProvider.GetYamlTemplate(testUser, filledFormData)
	assert.Nil(t, err)
	assert.NotNil(t, yamlTemplate)

	kibanaYaml := yamlTemplate.(provider_dtos2.ProviderYamlTemplateDto)
	assert.Equal(t, "MyCluster", kibanaYaml.Metadata.Name)
	assert.Equal(t, "MyNamespace", kibanaYaml.Metadata.Namespace)
	assert.Equal(t, "MyElasticSearchInstance", kibanaYaml.Spec.ElasticsearchRef.Name)
}

func Test_Kibana_Provider_End2End(t *testing.T) {
	kbProviderPtr, config := CreateKibanaTestProvider(t)
	kbProvider := *kbProviderPtr
	user := config.Users[0]

	url, err := url.Parse(config.Kubernetes.Server)
	assert.Nil(t, err)

	// Kibana depends on elastic search therefore we need to create it
	var esProvider provider.IServiceProvider = elasticsearch.CreateElasticSearchProvider(
		url.Hostname(),
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.ResourcesTemplatesPath,
	)

	// Since kibana provider utilizes the es provider we have to give a hint of other existing providers
	kbProvider.OnCoreInitialized([]*provider.IServiceProvider{&esProvider})

	// Create a new es instance
	esFormResponseDto := provider_dtos.ServiceCreationFormResponseDto{}
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
	filledForm := provider_dtos2.ServiceCreationFormResponseDto{}
	filledForm.Common.ClusterName = "kibana-test"
	filledForm.Common.ElasticSearchInstance = esFormResponseDto.Common.ClusterName

	service1Ptr := common_test.CommonProviderStart(t, kbProviderPtr, user, filledForm, 2)
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
	actionPtr, err = common_test.GetAction(service1Ptr, "Features", "cmd_kb_scale")
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
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_kb_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	assert.Equal(t, 2, placeholder.NumberOfReplicas)

	// Scale down -> 1
	tempServicePtr, err = esProvider.GetService(user, service1.GetName())
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_kb_scale")
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
	actionPtr, err = common_test.GetAction(tempServicePtr, "Features", "cmd_kb_scale")
	assert.Nil(t, err)
	action = *actionPtr
	placeholder = action.GetJsonFormResultPlaceholder().(*action_dtos.ClusterScaleDto)
	assert.Equal(t, 1, placeholder.NumberOfReplicas)

	// --- Exposure ---
	// Check if toggle is correct
	service1Ptr, err = esProvider.GetService(user, service1.GetName())
	assert.Nil(t, err)
	service1 = *service1Ptr
	toggleActionPtr, err := common_test.GetToggleAction(service1Ptr, "Security", "cmd_kb_expose_toggle")
	assert.Nil(t, err)
	toggleAction := *toggleActionPtr
	isSet, err := toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_kb_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation := result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)

	// Expose it
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_kb_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Set()
	assert.Nil(t, err)
	assert.Nil(t, result)
	time.Sleep(5 * time.Second)

	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_kb_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.True(t, isSet) // exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_kb_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)

	// Hide it again
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_kb_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	result, err = toggleAction.Unset()
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)

	// Check again if it is hidden
	// Check if toggle is correct
	toggleActionPtr, err = common_test.GetToggleAction(service1Ptr, "Security", "cmd_kb_expose_toggle")
	assert.Nil(t, err)
	toggleAction = *toggleActionPtr
	isSet, err = toggleAction.Get()
	assert.Nil(t, err)
	assert.False(t, isSet) // Not exposed

	// Check expose details
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_kb_get_expose_info")
	assert.Nil(t, err)
	action = *actionPtr
	result, err = action.GetActionExecuteCallback()(action.GetJsonFormResultPlaceholder())
	assert.Nil(t, err)
	clusterExposeInformation = result.(*action_dtos.ExposeInformations)
	assert.True(t, len(clusterExposeInformation.Host) > 0)

	// --- Certs ---
	// Check whether service is an Kibana service
	service2, ok := service1.(kibana.KibanaService)
	assert.True(t, ok)

	secret, _ := service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-kb-http-certs-internal")

	// Set cert
	actionPtr, err = common_test.GetAction(service1Ptr, "Security", "cmd_kb_set_cert_action")
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
	service3Ptr, err := common_test.WaitForServiceComeUp(kbProvider, user, service1.GetName())
	assert.Nil(t, err)
	service3 := *service3Ptr

	// Check whether status of service is ok after setting the certificate
	assert.True(t, service.ServiceStatusOk == service3.GetStatus())

	common_test.CommonProviderStop(t, kbProviderPtr, user)
	common_test.CommonProviderStop(t, &esProvider, user)

	// Wait till delete service is done
	time.Sleep(5 * time.Second)
	// Check whether the secret with associated certificate is also deleted
	secret, err = service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-tls-cert")
	assert.NotNil(t, err)
}
