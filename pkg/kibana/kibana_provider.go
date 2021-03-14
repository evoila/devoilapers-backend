package kibana

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"
	"path"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	"gopkg.in/yaml.v2"
)

// Implements IServiceProvider interface
// Use factory method CreateKibanaProvider to create
type KibanaProvider struct {
	provider.BasicProvider
}

// Factory method to create an instance of the KibanaProvider
func CreateKibanaProvider(host string, caPath string, templateDirectoryPath string) KibanaProvider {
	return KibanaProvider{provider.CreateCommonProvider(
		host,
		caPath,
		path.Join(templateDirectoryPath, "kibana.yaml"),
		"Kibana",
		"Kibana is an open source visualization tool mainly used to analyse logs.",
		"https://cdn.iconscout.com/icon/free/png-256/elastic-1-283281.png",
	)}
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
		services = append(services, kb.CrdInstanceToServiceInstance(auth, &kibanaInstance, KibanaCrd))
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

	return kb.CrdInstanceToServiceInstance(auth, &kibanaInstance, KibanaCrd), nil
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
func (kb KibanaProvider) CrdInstanceToServiceInstance(auth common.IKubernetesAuthInformation, crdInstance *v1.Kibana, crdApi *kubernetes.CommonCrdApi) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}
	mApi, _ := kubernetes.GenerateK8sApiFromToken(kb.Host, kb.CaPath, auth.GetKubernetesAccessToken())
	var KibanaService service.IService = KibanaService{
		status: crdInstance.Status.Health,
		BasicService: provider.BasicService{
			Name:              crdInstance.Name,
			ProviderType:      kb.GetServiceType(),
			Yaml:              string(yamlData),
			ImportantSections: (*kb.Template).GetImportantSections(),
		},
		api:    mApi,
		crdApi: crdApi,
		auth:   auth,
	}

	return &KibanaService
}
