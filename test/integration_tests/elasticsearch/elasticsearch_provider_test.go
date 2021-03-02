package common_test

import (
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/test/integration_tests/common_test"
	"github.com/stretchr/testify/assert"
	"testing"
)


func CreateElasticSearchProvider(t *testing.T) *elasticsearch.ElasticsearchProvider {
	config := common_test.GetConfig(t)

	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

	return &esProvider
}

func Test_Elasticsearch_Provider_GetAttributes(t *testing.T)  {
	esProvider := CreateElasticSearchProvider(t)
	assert.NotEqual(t, "", esProvider.GetServiceImage())
}

