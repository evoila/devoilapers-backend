package kibana

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"
	"path"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (es KibanaProvider) SetCertificateToService(auth common.IKubernetesAuthInformation, id string, tlcCert map[string][]byte) error {
	if api, err := kubernetes.GenerateK8sApiFromToken(es.Host, es.CaPath, auth.GetKubernetesAccessToken()); err != nil {
		return err
	} else {
		if elasticSearchCrd, err := es.createCrdApi(auth); err != nil {
			return err
		} else {
			kibanaInstance := v1.Kibana{}
			err = elasticSearchCrd.Get(auth.GetKubernetesNamespace(), id, ResourceName, &kibanaInstance)
			if err != nil {
				return err
			}
			if secretName, err := api.CreateTlsSecret(auth.GetKubernetesNamespace(), id, "Kibana", GroupName+"/"+GroupVersion, string(kibanaInstance.UID), tlcCert); err != nil {
				return err
			} else {
				kibanaInstance.Spec.HTTP.TLS.Certificate.SecretName = secretName
				kibanaInstance.ObjectMeta = metav1.ObjectMeta{
					Name:            kibanaInstance.Name,
					Namespace:       kibanaInstance.Namespace,
					ResourceVersion: kibanaInstance.ResourceVersion,
				}
				return elasticSearchCrd.Update(auth.GetKubernetesNamespace(), id, ResourceName, &kibanaInstance)
			}
		}
	}
}

// Converts a v1.Kibana instance to an service reprkbentation
func (kb KibanaProvider) CrdInstanceToServiceInstance(crdInstance *v1.Kibana) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var KibanaService service.IService = KibanaService{
		status: crdInstance.Status.Health,
		BasicService: provider.BasicService{
			Name:              crdInstance.Name,
			ProviderType:      kb.GetServiceType(),
			Yaml:              string(yamlData),
			ImportantSections: (*kb.Template).GetImportantSections(),
		},
	}

	return &KibanaService
}
