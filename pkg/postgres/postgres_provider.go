package postgres

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	common2 "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos"
	"OperatorAutomation/pkg/utils"
	"OperatorAutomation/pkg/utils/provider"
	"encoding/json"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	"gopkg.in/yaml.v2"
	"path"
	"strconv"
)

// Implements IServiceProvider interface
// Use factory method CreatePostgresProvider to create
type PostgresProvider struct {
	provider.BasicProvider
	nginxInformation kubernetes.NginxInformation
}

// Factory method to create an instance of the PostgresProvider
func CreatePostgresProvider(
	host string,
	caPath string,
	templateDirectoryPath string,
	nginxInformation kubernetes.NginxInformation) PostgresProvider {
	return PostgresProvider{
		nginxInformation: nginxInformation,
		BasicProvider: provider.CreateCommonProvider(
			host,
			caPath,
			path.Join(templateDirectoryPath, "postgres", "postgres.yaml"),
			path.Join(templateDirectoryPath, "postgres", "create_form.json"),
			"Postgres",
			"Postgres is an open source relational database management system.",
			"https://dashboard.snapcraft.io/site_media/appmedia/2018/11/postgresql-icon-256x256.jpg.png",
		)}
}

func (pg PostgresProvider) GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
	form := dtos.FormResponseDto{}
	err := json.Unmarshal(jsonFormResult, &form)
	if err != nil {
		return nil, err
	}

	// Create form with form default values
	yamlTemplate := dtos.ProviderYamlTemplateDto{}
	err = yaml.Unmarshal([]byte(pg.YamlTemplate), &yamlTemplate)
	if err != nil {
		return nil, err
	}

	// Transfer namespace
	yamlTemplate.Spec.Namespace = auth.GetKubernetesNamespace()
	yamlTemplate.Metadata.Namespace = auth.GetKubernetesNamespace()

	// Transfer form data to yaml template
	yamlTemplate.Metadata.Annotations.CurrentPrimary = form.Common.ClusterName
	yamlTemplate.Metadata.Labels.CrunchyPghaScope = form.Common.ClusterName
	yamlTemplate.Metadata.Labels.DeploymentName = form.Common.ClusterName
	yamlTemplate.Metadata.Labels.Name = form.Common.ClusterName
	yamlTemplate.Metadata.Labels.PgCluster = form.Common.ClusterName
	yamlTemplate.Metadata.Name = form.Common.ClusterName
	yamlTemplate.Spec.PrimaryStorage.Name = form.Common.ClusterName
	yamlTemplate.Spec.Clustername = form.Common.ClusterName
	yamlTemplate.Spec.Database = form.Common.ClusterName
	yamlTemplate.Spec.Name = form.Common.ClusterName

	yamlTemplate.Spec.User = form.Common.Username
	yamlTemplate.Spec.Port = strconv.Itoa(form.Common.InClusterPort)
	yamlTemplate.Spec.PrimaryStorage.Size = strconv.Itoa(form.Common.ClusterStorageSize) + "G"
	yamlTemplate.Spec.BackrestStorage.Size = yamlTemplate.Spec.PrimaryStorage.Size
	yamlTemplate.Spec.ReplicaStorage.Size = yamlTemplate.Spec.PrimaryStorage.Size

	return yamlTemplate, nil
}

func (pg PostgresProvider) GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error) {
	// Create form with form default values
	formsQuery := dtos.FormQueryDto{}
	//formsQuery := map[string]interface{}{}
	err := json.Unmarshal([]byte(pg.FormTemplate), &formsQuery)
	if err != nil {
		return nil, err
	}

	// Set a default name
	formsQuery.Properties.Common.Properties.ClusterName.Default = utils.GetRandomKubernetesResourceName()

	return formsQuery, nil
}

func (pg PostgresProvider) createCrdApi(auth common.IKubernetesAuthInformation) (*kubernetes.CommonCrdApi, error) {
	return kubernetes.CreateCommonCrdApi(pg.Host, pg.CaPath, auth.GetKubernetesAccessToken(), GroupName, GroupVersion)
}

func (pg PostgresProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	postgresCrd, err := pg.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	postgresInstances := v1.PgclusterList{}
	err = postgresCrd.List(auth.GetKubernetesNamespace(), ResourceName, &postgresInstances)
	if err != nil {
		return nil, err
	}

	var services []*service.IService
	for _, postgresInstanceIterator := range postgresInstances.Items {
		postgresInstance := postgresInstanceIterator
		services = append(services, pg.CrdInstanceToServiceInstance(postgresCrd, auth, &postgresInstance))
	}

	return services, nil
}

func (pg PostgresProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	postgresCrd, err := pg.createCrdApi(auth)

	if err != nil {
		return nil, err
	}

	postgresInstance := v1.Pgcluster{}
	err = postgresCrd.Get(auth.GetKubernetesNamespace(), id, ResourceName, &postgresInstance)
	if err != nil {
		return nil, err
	}

	return pg.CrdInstanceToServiceInstance(postgresCrd, auth, &postgresInstance), nil
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

// Converts a v1.Pgcluster instance to a service representation
func (pg PostgresProvider) CrdInstanceToServiceInstance(
	crdClient *kubernetes.CommonCrdApi,
	auth common.IKubernetesAuthInformation,
	crdInstance *v1.Pgcluster) *service.IService {

	yamlData, err := yaml.Marshal(crdInstance)
	if err != nil {
		yamlData = []byte("Unknown")
	}

	var postgresService service.IService = PostgresService{
		PostgresServiceInformations: common2.PostgresServiceInformations{
			ClusterInstance:  crdInstance,
			Auth:             auth,
			Host:             pg.Host,
			CaPath:           pg.CaPath,
			CrdClient:        crdClient,
			NginxInformation: pg.nginxInformation,
		},
		BasicService: provider.BasicService{
			Name:         crdInstance.Name,
			ProviderType: pg.GetServiceType(),
			Yaml:         string(yamlData),
		},
	}

	return &postgresService
}
