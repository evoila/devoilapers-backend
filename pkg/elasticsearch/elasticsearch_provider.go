package elasticsearch

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/crd"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

// Should implement ServiceProvider interface
type ElasticSearchProvider struct {
	template *service.IServiceTemplate
	host string
	caPath string
}

func CreateElasticSearchProvider(host string, caPath string, templatePath string) ElasticSearchProvider {
	yamlTemplatePath := path.Join(templatePath, "elasticsearch.yaml")
	templateData, err := ioutil.ReadFile(yamlTemplatePath)
	if err != nil {
		panic("Yaml template could not be found under path: " + yamlTemplatePath)
	}

	var template service.IServiceTemplate = service.ServiceTemplate{
		Yaml: string(templateData),
		ImportantSections: []string{"name"},
	}

	return ElasticSearchProvider{
		template: &template,
		host: host,
		caPath: caPath,
	}
}

func (es ElasticSearchProvider) GetServiceDescription() string {
	return "Elasticsearch is a distributed, free and open search and analytics engine for all types of data."
}

func (es ElasticSearchProvider) GetServiceImage() string {
	return "https://cdn.iconscout.com/icon/free/png-256/elasticsearch-226094.png"
}

func (es ElasticSearchProvider) GetServiceType() string {
	return "Elasticsearch"
}

func (es ElasticSearchProvider) GetTemplate(auth common.IKubernetesAuthInformation) *service.IServiceTemplate {
	originalTemplate := (*es.template)
	yamlTemplate := originalTemplate.GetYAML()
	yamlTemplate = utils.FillWithData(auth, yamlTemplate)

	var template service.IServiceTemplate = service.ServiceTemplate{
		Yaml: yamlTemplate,
		ImportantSections: originalTemplate.GetImportantSections(),
	}

	return &template
}

func (es ElasticSearchProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	elasticSearchCrd, err := crd.GenerateEsApiBasedOnToken(es.host, es.caPath, auth.GetKubernetesAccessToken())

	if err != nil {
		return nil, err
	}

	elasticSearchInstances, err := elasticSearchCrd.List(auth.GetKubernetesNamespace())
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, elasticSearchInstance := range elasticSearchInstances.Items {
		services = append(services, es.CrdInstanceToServiceInstance(&elasticSearchInstance))
	}

	return services, nil
}

func (es ElasticSearchProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	elasticSearchCrd, err := crd.GenerateEsApiBasedOnToken(es.host, es.caPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return nil, err
	}

	crdInstance, err := elasticSearchCrd.Get(auth.GetKubernetesNamespace(), id)
	if err != nil {
		return nil, err
	}

	return es.CrdInstanceToServiceInstance(crdInstance), nil
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
	elasticSearchCrd, err := crd.GenerateEsApiBasedOnToken(es.host, es.caPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	//TODO: Check if there is an associated ingress

	return elasticSearchCrd.Delete(auth.GetKubernetesNamespace(), id)
}

// Converts a crd.Elasticsearch instance to an service representation
func (es ElasticSearchProvider) CrdInstanceToServiceInstance(crdInstance *v1.Elasticsearch) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var elasticSearchService service.IService = ElasticSearchService{
		status: crdInstance.Status.Health,
		name : crdInstance.Name,
		providerType: es.GetServiceType(),
		yaml: string(yamlData),
		importantSections: (*es.template).GetImportantSections(),
	}

	return &elasticSearchService
}

