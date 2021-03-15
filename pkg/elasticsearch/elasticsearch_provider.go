package elasticsearch

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"
	"path"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"

	"gopkg.in/yaml.v2"
)

// Implements IServiceProvider interface
// Use factory method CreateElasticSearchProvider to create
type ElasticsearchProvider struct {
	provider.BasicProvider
}

// Factory method to create an instance of the ElasticsearchProvider
func CreateElasticSearchProvider(host string, caPath string, templateDirectoryPath string) ElasticsearchProvider {
	return ElasticsearchProvider{provider.CreateCommonProvider(
		host,
		caPath,
		path.Join(templateDirectoryPath, "elasticsearch.yaml"),
		"Elasticsearch",
		"Elasticsearch is a distributed, free and open search and analytics engine for all types of data.",
		"https://cdn.iconscout.com/icon/free/png-256/elasticsearch-226094.png",
	)}
}

func (es ElasticsearchProvider) createCrdApi(auth common.IKubernetesAuthInformation) (*kubernetes.CommonCrdApi, error) {
	return kubernetes.CreateCommonCrdApi(es.Host, es.CaPath, auth.GetKubernetesAccessToken(), GroupName, GroupVersion)
}

func (es ElasticsearchProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	elasticSearchCrd, err := es.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	elasticSearchInstances := v1.ElasticsearchList{}
	err = elasticSearchCrd.List(auth.GetKubernetesNamespace(), RessourceName, &elasticSearchInstances)
	if err != nil {
		return nil, err
	}

	api, err := kubernetes.GenerateK8sApiFromToken(es.Host, es.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}
	commonCrdApi, err := es.createCrdApi(auth)
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, elasticSearchInstance := range elasticSearchInstances.Items {
		services = append(services, es.CrdInstanceToServiceInstance(api, commonCrdApi, &elasticSearchInstance))
	}

	return services, nil
}

func (es ElasticsearchProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	elasticSearchCrd, err := es.createCrdApi(auth)
	if err != nil {
		return nil, err
	}

	elasticSearchInstance := v1.Elasticsearch{}
	err = elasticSearchCrd.Get(auth.GetKubernetesNamespace(), id, RessourceName, &elasticSearchInstance)
	if err != nil {
		return nil, err
	}

	api, err := kubernetes.GenerateK8sApiFromToken(es.Host, es.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}
	commonCrdApi, err := es.createCrdApi(auth)
	if err != nil {
		return nil, err
	}

	return es.CrdInstanceToServiceInstance(api, commonCrdApi, &elasticSearchInstance), nil
}

func (es ElasticsearchProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	api, err := kubernetes.GenerateK8sApiFromToken(es.Host, es.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	_, err = api.Apply([]byte(yaml))
	if err != nil {
		return err
	}

	return nil
}

func (es ElasticsearchProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	elasticSearchCrd, err := es.createCrdApi(auth)
	if err != nil {
		return err
	}

	//TODO: Check if there is an associated ingress
	return elasticSearchCrd.Delete(auth.GetKubernetesNamespace(), id, RessourceName)
}

// Converts a v1.Elasticsearch instance to an service representation
func (es ElasticsearchProvider) CrdInstanceToServiceInstance(api *kubernetes.K8sApi, commonCrdApi *kubernetes.CommonCrdApi, crdInstance *v1.Elasticsearch) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var elasticSearchService service.IService = ElasticSearchService{
		K8sApi:       api,
		commonCrdApi: commonCrdApi,
		crdInstance:  crdInstance,
		status:       crdInstance.Status.Health,
		BasicService: provider.BasicService{
			Name:              crdInstance.Name,
			ProviderType:      es.GetServiceType(),
			Yaml:              string(yamlData),
			ImportantSections: (*es.Template).GetImportantSections(),
		},
	}

	return &elasticSearchService
}
