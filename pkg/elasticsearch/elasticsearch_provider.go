package elasticsearch

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/provider_dtos"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils"
	"OperatorAutomation/pkg/utils/provider"
	"encoding/json"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	"gopkg.in/yaml.v2"
	"path"
)

// Implements IServiceProvider interface
// Use factory method CreateElasticSearchProvider to create
type ElasticsearchProvider struct {
	kubernetesHostname string
	provider.BasicProvider
}

// Factory method to create an instance of the ElasticsearchProvider
func CreateElasticSearchProvider(
	kubernetesHostname string,
	kuberntesServer string, caPath string, templateDirectoryPath string) ElasticsearchProvider {
	return ElasticsearchProvider{
		kubernetesHostname: kubernetesHostname,
		BasicProvider: provider.CreateCommonProvider(
		kuberntesServer,
		caPath,
		path.Join(templateDirectoryPath, "elasticsearch", "elasticsearch.yaml"),
		path.Join(templateDirectoryPath, "elasticsearch", "create_form.json"),
		"Elasticsearch",
		"Elasticsearch is a distributed, free and open search and analytics engine for all types of data.",
		"https://cdn.iconscout.com/icon/free/png-256/elasticsearch-226094.png",
	)}
}

func (es ElasticsearchProvider) GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
	form := provider_dtos.ServiceCreationFormResponseDto{}
	err := json.Unmarshal(jsonFormResult, &form)
	if err != nil {
		return "", err
	}

	// Create form with form default values
	yamlTemplate := provider_dtos.ProviderYamlTemplateDto{}
	err = yaml.Unmarshal([]byte(es.YamlTemplate), &yamlTemplate)
	if err != nil {
		return "", err
	}

	// Transfer name to the final creation yaml
	yamlTemplate.Metadata.Name = form.Common.ClusterName
	yamlTemplate.Metadata.Namespace = auth.GetKubernetesNamespace()

	return yamlTemplate, nil
}

func (es ElasticsearchProvider) GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error) {
	// Create form with form default values
	formsQuery := provider_dtos.ServiceCreationFormDto{}
	err := json.Unmarshal([]byte(es.FormTemplate), &formsQuery)
	if err != nil {
		return "", err
	}

	// Set a default name
	formsQuery.Properties.Common.Properties.ClusterName.Default = utils.GetRandomKubernetesResourceName("es")
	return formsQuery, nil
}

func (es ElasticsearchProvider) createCrdApi(auth common.IKubernetesAuthInformation) (*kubernetes.CommonCrdApi, error) {
	return kubernetes.CreateCommonCrdApi(es.KubernetsServer, es.CaPath, auth.GetKubernetesAccessToken(), esCommon.GroupName, esCommon.GroupVersion)
}

func (es ElasticsearchProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	elasticSearchCrd, err := es.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	elasticSearchInstances := v1.ElasticsearchList{}
	err = elasticSearchCrd.List(auth.GetKubernetesNamespace(), esCommon.RessourceName, &elasticSearchInstances)
	if err != nil {
		return nil, err
	}

	api, err := kubernetes.GenerateK8sApiFromToken(es.KubernetsServer, es.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, elasticSearchInstanceIterator := range elasticSearchInstances.Items {
		elasticSearchInstance := elasticSearchInstanceIterator
		services = append(services, es.CrdInstanceToServiceInstance(api, elasticSearchCrd, &elasticSearchInstance))
	}

	return services, nil
}

func (es ElasticsearchProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	elasticSearchCrd, err := es.createCrdApi(auth)
	if err != nil {
		return nil, err
	}

	elasticSearchInstance := v1.Elasticsearch{}
	err = elasticSearchCrd.Get(auth.GetKubernetesNamespace(), id, esCommon.RessourceName, &elasticSearchInstance)
	if err != nil {
		return nil, err
	}

	api, err := kubernetes.GenerateK8sApiFromToken(es.KubernetsServer, es.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	return es.CrdInstanceToServiceInstance(api, elasticSearchCrd, &elasticSearchInstance), nil
}

func (es ElasticsearchProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	api, err := kubernetes.GenerateK8sApiFromToken(es.KubernetsServer, es.CaPath, auth.GetKubernetesAccessToken())
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
	return elasticSearchCrd.Delete(auth.GetKubernetesNamespace(), id, esCommon.RessourceName)
}

// Converts a v1.Elasticsearch instance to an service representation
func (es ElasticsearchProvider) CrdInstanceToServiceInstance(api *kubernetes.K8sApi, commonCrdApi *kubernetes.CommonCrdApi, crdInstance *v1.Elasticsearch) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var elasticSearchService service.IService = ElasticSearchService{
		ElasticsearchServiceInformations: esCommon.ElasticsearchServiceInformations{
			Hostname: es.kubernetesHostname,
			K8sApi:       api,
			CrdClient: commonCrdApi,
			ClusterInstance:  crdInstance,
		},
		status:       crdInstance.Status.Health,
		BasicService: provider.BasicService{
			Name:         crdInstance.Name,
			ProviderType: es.GetServiceType(),
			Yaml:         string(yamlData),
		},
	}

	return &elasticSearchService
}
