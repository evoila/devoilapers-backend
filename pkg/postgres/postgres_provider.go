package postgres

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	common2 "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/utils/provider"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	"gopkg.in/yaml.v2"
	"path"
)

// Implements IServiceProvider interface
// Use factory method CreatePostgresProvider to create
type PostgresProvider struct {
	provider.BasicProvider
}

// Factory method to create an instance of the PostgresProvider
func CreatePostgresProvider(host string, caPath string, templateDirectoryPath string) PostgresProvider {
	return PostgresProvider{provider.CreateCommonProvider(
		host,
		caPath,
		path.Join(templateDirectoryPath, "postgres.yaml"),
		"Postgres",
		"Postgres is an open source relational database management system.",
		"https://dashboard.snapcraft.io/site_media/appmedia/2018/11/postgresql-icon-256x256.jpg.png",
	)}
}

func (pg PostgresProvider) createCrdApi(auth common.IKubernetesAuthInformation) (*kubernetes.CommonCrdApi, error)  {
	return  kubernetes.CreateCommonCrdApi(pg.Host, pg.CaPath, auth.GetKubernetesAccessToken(), GroupName, GroupVersion)
}

func (pg PostgresProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	postgresCrd,  err := pg.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	kibanaInstances := v1.PgclusterList{}
	err = postgresCrd.List(auth.GetKubernetesNamespace(), ResourceName, &kibanaInstances)
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, kibanaInstance := range kibanaInstances.Items {
		services = append(services, pg.CrdInstanceToServiceInstance(auth, &kibanaInstance))
	}

	return services, nil
}

func (pg PostgresProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	postgresCrd, err := pg.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	kibanaInstance := v1.Pgcluster{}
	err = postgresCrd.Get(auth.GetKubernetesNamespace(), id, ResourceName, &kibanaInstance)
	if err != nil {
		return nil, err
	}

	return pg.CrdInstanceToServiceInstance(auth, &kibanaInstance), nil
}

func (pg PostgresProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		return err
	}

	_, err = api.Apply([]byte(yaml))
	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	postgresCrd, err := pg.createCrdApi(auth)
	if err != nil {
		return err
	}

	//TODO: Check if there is an associated ingress
	return postgresCrd.Delete(auth.GetKubernetesNamespace(), id, ResourceName)
}

// Converts a v1.Pgcluster instance to an service representation
func (pg PostgresProvider) CrdInstanceToServiceInstance(auth common.IKubernetesAuthInformation, crdInstance *v1.Pgcluster) *service.IService {
	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var postgresService service.IService = PostgresService{
		PostgresServiceInformations: common2.PostgresServiceInformations{
			ClusterInstance: crdInstance,
			Auth:            auth,
			Host:            pg.Host,
			CaPath:          pg.CaPath,
		},
		BasicService: provider.BasicService{
			Name:              crdInstance.Name,
			ProviderType:      pg.GetServiceType(),
			Yaml:              string(yamlData),
			ImportantSections: (*pg.Template).GetImportantSections(),
		},
	}

	return &postgresService
}
