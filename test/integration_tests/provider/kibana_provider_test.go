package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/kibana"
	"OperatorAutomation/pkg/kibana/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateKibanaTestProvider(t *testing.T) (*kibana.KibanaProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	kbProvider := kibana.CreateKibanaProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

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

	assert.NotEqual(t, "", kbProvider.GetServiceImage())
	assert.NotEqual(t, "", kbProvider.GetServiceDescription())
	assert.Equal(t, "Kibana", kbProvider.GetServiceType())

	testUser := unit_test.TestUser{
		KubernetesNamespace: "A_LONG_NAMESPACE",
	}

	template := *kbProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template.GetImportantSections()))
	assert.Equal(t, "metadata.name", template.GetImportantSections()[0])

	template2 := *kbProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template2.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template2.GetImportantSections()))
	assert.Equal(t, "metadata.name", template2.GetImportantSections()[0])
	assert.NotEqual(t, template2.GetYAML(), template.GetYAML())
}

func Test_Kibana_Provider_End2End(t *testing.T) {
	kbProvider, config := CreateKibanaTestProvider(t)

	user := config.Users[0]
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}

	// Kibana depends on elastic search therefore we need to create it
	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath,
	)

	// Create a new es instance
	esYaml := (*esProvider.GetTemplate(user)).GetYAML()
	err := esProvider.CreateService(user, esYaml)
	assert.Nil(t, err)
	esServices, err := esProvider.GetServices(user)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(esServices))
	esService0 := *esServices[0]
	time.Sleep(10 * time.Second)

	// Continue with actual kb provider
	yaml := (*kbProvider.GetTemplate(user)).GetYAML()
	// Fill in the name of the es instance. Since the yaml needs a reference to come up ok
	yaml = strings.Replace(yaml, "YOUR_ELASTICSEARCH_INSTANCE_NAME", esService0.GetName(), 1)

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
	assert.Equal(t, 1, len(service0.GetActions()))

	// Try get service with invalid user data
	_, err = kbProvider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok. Kibana needs some extra time.
	var service1 service.IService
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

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
	assert.Equal(t, service0.GetTemplate().GetImportantSections(), service1.GetTemplate().GetImportantSections())

	// Check whether service is an Kibana service
	service2, ok := service1.(kibana.KibanaService)
	assert.True(t, ok)

	secret, _ := service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-kb-http-certs-internal")

	// Test set certificate to service
	certDto := &dtos.CertificateDto{
		CaCrt:  string(secret.Data["ca.crt"]),
		TlsCrt: string(secret.Data["tls.crt"]),
		TlsKey: string(secret.Data["tls.key"]),
	}

	_, err = service2.SetCertificateToService(certDto)
	assert.Nil(t, err)

	// Check status of service after setting the certificate
	var service3 *service.IService
	for i := 0; i < 5; i++ {
		tmpService, err := kbProvider.GetService(user, service2.GetName())
		assert.Nil(t, err)
		assert.NotNil(t, tmpService)
		if (*tmpService).GetStatus() == service.ServiceStatusOk {
			service3 = tmpService
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	// Check whether service is not nil
	assert.NotNil(t, service3)

	// Check whether status of service is ok after setting the certificate
	assert.True(t, service.ServiceStatusOk == (*service3).GetStatus())

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

	// Wait till delete service is done
	time.Sleep(5 * time.Second)
	// Check whether the secret with associated certificate is also deleted
	secret, err = service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-tls-cert")
	assert.NotNil(t, err)
}
