package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/test/integration_tests/common_test"
	unit_test "OperatorAutomation/test/unit_tests/common_test"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func CreateElasticSearchTestProvider(t *testing.T) (*elasticsearch.ElasticsearchProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

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
		KubernetesNamespace: "A_LONG_NAMESPACE",
	}

	template := *esProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template.GetImportantSections()))
	assert.Equal(t, "metadata.name", template.GetImportantSections()[0])

	template2 := *esProvider.GetTemplate(testUser)
	assert.True(t, strings.Contains(template2.GetYAML(), "namespace: "+testUser.KubernetesNamespace))
	assert.Equal(t, 1, len(template2.GetImportantSections()))
	assert.Equal(t, "metadata.name", template2.GetImportantSections()[0])
	assert.NotEqual(t, template2.GetYAML(), template.GetYAML())
}

func Test_Elasticsearch_Provider_End2End(t *testing.T) {
	esProvider, config := CreateElasticSearchTestProvider(t)

	user := config.Users[0]
	invalidUser := unit_test.TestUser{KubernetesNamespace: "namespace", KubernetesAccessToken: "InvalidToken"}
	yaml := (*esProvider.GetTemplate(user)).GetYAML()

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
	assert.Equal(t, 1, len(service0.GetActions()))
	assert.True(t,
		service.ServiceStatusPending == service0.GetStatus() ||
			service.ServiceStatusOk == service0.GetStatus(),
	)

	// Try get service with invalid user data
	_, err = esProvider.GetService(invalidUser, service0.GetName())
	assert.NotNil(t, err)

	// Wait for service to become ok
	var service1 service.IService
	for i := 0; i < 12; i++ {
		time.Sleep(5 * time.Second)

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
	assert.Equal(t, service0.GetTemplate().GetImportantSections(), service1.GetTemplate().GetImportantSections())

	// Check whether service is an Elasticsearch service
	service2, ok := service1.(elasticsearch.ElasticSearchService)
	assert.True(t, ok)

	secret, _ := service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-es-http-certs-internal")

	// Test set certificate to service
	certDto := &dtos.CertificateDto{
		CaCrt:  base64.StdEncoding.EncodeToString(secret.Data["ca.crt"]),
		TlsCrt: base64.StdEncoding.EncodeToString(secret.Data["tls.crt"]),
		TlsKey: base64.StdEncoding.EncodeToString(secret.Data["tls.key"]),
	}

	_, err = service2.SetCertificateToService(certDto)
	assert.Nil(t, err)

	// Check status of service after setting the certificate
	var service3 service.IService
	for i := 0; i < 10; i++ {
		tmpService, err := esProvider.GetService(user, service0.GetName())
		assert.Nil(t, err)
		assert.NotNil(t, tmpService)
		if (*tmpService).GetStatus() == service.ServiceStatusOk {
			service3 = *tmpService
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	assert.NotNil(t, service3)
	assert.True(t, service.ServiceStatusOk == service3.GetStatus())

	// Try delete service with invalid id
	err = esProvider.DeleteService(user, "some-not-existing-id")
	assert.NotNil(t, err)

	// Try delete service with invalid user
	err = esProvider.DeleteService(invalidUser, (*services[0]).GetName())
	assert.NotNil(t, err)

	// Delete service
	err = esProvider.DeleteService(user, (*services[0]).GetName())
	assert.Nil(t, err)

	// Wait till delete service is done
	time.Sleep(5 * time.Second)
	// Check whether the secret with associated certificate is also deleted
	secret, err = service2.K8sApi.GetSecret(user.KubernetesNamespace, service2.GetName()+"-tls-cert")
	assert.NotNil(t, err)
}
