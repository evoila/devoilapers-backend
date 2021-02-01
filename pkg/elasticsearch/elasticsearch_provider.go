package elasticsearch

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/crd"
	"OperatorAutomation/pkg/kubernetes"
	"io/ioutil"
)

// Should implement ServiceProvider interface
type ElasticSearchProvider struct {
	template *service.IServiceTemplate
	host string
	caPath string
}

func CreateElasticSearchProvider(host string, caPath string) ElasticSearchProvider {
	data, err := ioutil.ReadFile("configs/yaml_templates/elasticsearch.yaml")
	if err != nil {
		panic("Missing file")
	}

	var template service.IServiceTemplate = service.ServiceTemplate{
		Yaml: string(data),
		ImportantSections: []string{"name"},
	}

	return ElasticSearchProvider{
		template: &template,
		host: host,
		caPath: caPath,
	}
}

func (es ElasticSearchProvider) GetServiceDescription() string {
	return "Elastic search description"
}

func (es ElasticSearchProvider) GetServiceImage() string {
	return "Base64 img"
}

func (es ElasticSearchProvider) GetServiceType() string {
	return "ElasticSearch"
}

func (es ElasticSearchProvider) GetTemplate() *service.IServiceTemplate {
	// Todo: Add namespace/random name based on a struct instead of a file
	return es.template
}

func (es ElasticSearchProvider) GetServices(auth common.IKubernetesAuthInformation) []*service.IService {
	elasticSearchCrd, err := crd.GenerateEsApiBasedOnToken(es.host, es.caPath, auth.GetKubernetesAccessToken())

	if err != nil {
		return nil
	}

	elasticSearchInstances, err := elasticSearchCrd.List(auth.GetKubernetesNamespace())
	if err != nil {
		return nil
	}

	var services []*service.IService
	for _, elasticSearchInstance := range elasticSearchInstances.Items {
		var elasticSearchService service.IService = ElasticSearchService{
			serviceStatus: elasticSearchInstance.Status.Health,
			serviceName : elasticSearchInstance.Name,
			serviceType: es.GetServiceType(),
		}

		services = append(services, &elasticSearchService)
	}

	return services
}

func (es ElasticSearchProvider) GetService(auth common.IKubernetesAuthInformation, id string) *service.IService {
	var els service.IService = ElasticSearchService{serviceType: es.GetServiceType()}
	return &els
}

func (es ElasticSearchProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	api, err := kubernetes.GenerateK8sApiFromToken(es.host, es.caPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	_, err = api.Apply([]byte(yaml))
	if err != nil {
		return err
	}

	return nil
}

func (es ElasticSearchProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	panic("implement me")
}

