package provider

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/elasticsearch"
	"OperatorAutomation/test/integration_tests/common_test"

	"testing"
)

func CreateElasticSearchProvider(t *testing.T) (*elasticsearch.ElasticsearchProvider, config.RawConfig) {
	config := common_test.GetConfig(t)

	esProvider := elasticsearch.CreateElasticSearchProvider(
		config.Kubernetes.Server,
		config.Kubernetes.CertificateAuthority,
		config.YamlTemplatePath)

	return &esProvider, config
}

func Test_Elastic_Provider_Expose(t *testing.T) {
	esProvider, config := CreateElasticSearchProvider(t)
	user := config.Users[0]
	url, err := esProvider.ExposeThroughIngress(user, "my-ingress", "elasticsearch-sample-es-http", "myhosst.com")
	t.Error(url, err)

}
