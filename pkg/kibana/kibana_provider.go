package kibana

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kibana/dtos"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils"
	providerUtils "OperatorAutomation/pkg/utils/provider"
	"encoding/json"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	"gopkg.in/yaml.v2"
	"path"
	"strings"
)

// Implements IServiceProvider interface
// Use factory method CreateKibanaProvider to create
type KibanaProvider struct {
	sharedKibanaData *sharedKibanaData
	providerUtils.BasicProvider
}

// Data holder class to overcome interface pointer problems
type sharedKibanaData struct {
	esProvider *provider.IServiceProvider
}

// Factory method to create an instance of the KibanaProvider
func CreateKibanaProvider(host string, caPath string, templateDirectoryPath string) KibanaProvider {
	return KibanaProvider{sharedKibanaData: &sharedKibanaData{}, BasicProvider: providerUtils.CreateCommonProvider(
		host,
		caPath,
		path.Join(templateDirectoryPath, "kibana", "kibana.yaml"),
		path.Join(templateDirectoryPath, "kibana", "create_form.json"),
		"Kibana",
		"Kibana is an open source visualization tool mainly used to analyse logs.",
		"https://cdn.iconscout.com/icon/free/png-256/elastic-1-283281.png",
	)}
}

func (kb KibanaProvider) OnCoreInitialized(providers []*provider.IServiceProvider) {
	// Safe elasticsearch provider to satisfy form later on
	for idx, provider := range providers {
		if strings.ToLower((*provider).GetServiceType()) == "elasticsearch" {
			kb.sharedKibanaData.esProvider = providers[idx]
		}
	}

	if kb.sharedKibanaData.esProvider == nil {
		panic("Elasticsearch provider could not be resolved but is necessary for kibana")
	}
}

func (kb KibanaProvider) GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
	form := dtos.FormResponseDto{}
	err := json.Unmarshal(jsonFormResult, &form)
	if err != nil {
		return nil, err
	}

	// Create form with form default values
	yamlTemplate := dtos.ProviderYamlTemplateDto{}
	err = yaml.Unmarshal([]byte(kb.YamlTemplate), &yamlTemplate)
	if err != nil {
		return nil, err
	}

	// Transfer name to the final creation yaml
	yamlTemplate.Metadata.Name = form.Common.ClusterName
	yamlTemplate.Metadata.Namespace = auth.GetKubernetesNamespace()
	yamlTemplate.Spec.ElasticsearchRef.Name = form.Common.ElasticSearchInstance

	return yamlTemplate, nil
}

func (kb KibanaProvider) GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error) {
	// Create form with form default values
	formsQuery := dtos.FormQueryDto{}
	err := json.Unmarshal([]byte(kb.FormTemplate), &formsQuery)
	if err != nil {
		return nil, err
	}

	// Query elastic search instances
	esServices, err := (*kb.sharedKibanaData.esProvider).GetServices(auth)
	if err != nil {
		return nil, err
	}

	// Set a default name
	formsQuery.Properties.Common.Properties.ClusterName.Default = utils.GetRandomKubernetesResourceName()

	// Append elastic search instances as choosable reference
	for _, esServicePtr := range esServices {
		esService := *esServicePtr
		formsQuery.Properties.Common.Properties.ElasticSearchInstance.OneOf =
			append(formsQuery.Properties.Common.Properties.ElasticSearchInstance.OneOf,
				dtos.OneOf{
					Description: esService.GetName(),
					Enum:        []string{esService.GetName()},
				},
			)
	}

	return formsQuery, nil
}

func (kb KibanaProvider) createCrdApi(auth common.IKubernetesAuthInformation) (*kubernetes.CommonCrdApi, error) {
	return kubernetes.CreateCommonCrdApi(kb.Host, kb.CaPath, auth.GetKubernetesAccessToken(), GroupName, GroupVersion)
}

func (kb KibanaProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	KibanaCrd, err := kb.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	kibanaInstances := v1.KibanaList{}
	err = KibanaCrd.List(auth.GetKubernetesNamespace(), ResourceName, &kibanaInstances)
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, kibanaInstance := range kibanaInstances.Items {
		services = append(services, kb.CrdInstanceToServiceInstance(&kibanaInstance))
	}

	return services, nil
}

func (kb KibanaProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	KibanaCrd, err := kb.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	kibanaInstance := v1.Kibana{}
	err = KibanaCrd.Get(auth.GetKubernetesNamespace(), id, ResourceName, &kibanaInstance)
	if err != nil {
		return nil, err
	}

	return kb.CrdInstanceToServiceInstance(&kibanaInstance), nil
}

func (kb KibanaProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	api, err := kubernetes.GenerateK8sApiFromToken(kb.Host, kb.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	_, err = api.Apply([]byte(yaml))
	if err != nil {
		return err
	}

	return nil
}

func (kb KibanaProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	KibanaCrd, err := kb.createCrdApi(auth)
	if err != nil {
		return err
	}

	//TODO: Check if there is an associated ingress
	return KibanaCrd.Delete(auth.GetKubernetesNamespace(), id, ResourceName)
}

// Converts a v1.Kibana instance to an service reprkbentation
func (kb KibanaProvider) CrdInstanceToServiceInstance(crdInstance *v1.Kibana) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var KibanaService service.IService = KibanaService{
		status: crdInstance.Status.Health,
		BasicService: providerUtils.BasicService{
			Name:         crdInstance.Name,
			ProviderType: kb.GetServiceType(),
			Yaml:         string(yamlData),
		},
	}

	return &KibanaService
}
