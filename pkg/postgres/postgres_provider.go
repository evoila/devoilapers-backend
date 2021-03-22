package postgres

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/kubernetes/yaml_types"
	common2 "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos"
	"OperatorAutomation/pkg/utils"
	"OperatorAutomation/pkg/utils/provider"
	"encoding/json"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"path"
	"strconv"
	"strings"
)

const PostgresProviderLogPrefix = "File: postgres_provider.go: "

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

	log.Info(PostgresProviderLogPrefix + "Creating new postgres provider.")

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
	log.Trace(PostgresProviderLogPrefix + "Going to convert received form data to yaml.")

	form := dtos.FormResponseDto{}
	err := json.Unmarshal(jsonFormResult, &form)
	if err != nil {
		log.Error(PostgresProviderLogPrefix + "Could not unmarshal received form to generate the yaml.")
		return nil, err
	}

	// Create form with form default values
	yamlTemplate := dtos.ProviderYamlTemplateDto{}
	err = yaml.Unmarshal([]byte(pg.YamlTemplate), &yamlTemplate)
	if err != nil {
		log.Error(PostgresProviderLogPrefix + "Could not unmarshal the default postgres yaml template.")
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

	yamlTemplate.Spec.User = strings.ToLower(form.Common.Username)
	yamlTemplate.Spec.Port = strconv.Itoa(form.Common.InClusterPort)
	yamlTemplate.Spec.PrimaryStorage.Size = strconv.Itoa(form.Common.ClusterStorageSize) + "G"
	yamlTemplate.Spec.BackrestStorage.Size = yamlTemplate.Spec.PrimaryStorage.Size
	yamlTemplate.Spec.ReplicaStorage.Size = yamlTemplate.Spec.PrimaryStorage.Size

	// Check if tls should be used
	if form.TLS.UseTLS {
		log.Trace(PostgresProviderLogPrefix + "Tls is requested during yaml conversion")

		if form.TLS.TLSMode == "TlsFromFile" {
			log.Trace(PostgresProviderLogPrefix + "Use tls mode from file. Going to construct secrets.")

			// Construct secrets accordingly to
			// https://access.crunchydata.com/documentation/postgres-operator/4.6.1/tutorial/tls/
			caSecret := &yaml_types.YamlCaSecret{
				Metadata: yaml_types.Metadata{
					Name: form.Common.ClusterName + "-ca",
					Namespace: auth.GetKubernetesNamespace(),
				},
				Kind: "Secret",
				APIVersion: "v1",
				Data: yaml_types.CaData{
					CaCrtBase64: form.TLS.TLSModeFromFile.CaCertBase64,
				},
				Type: "Opaque",
			}

			tlsSecret :=  &yaml_types.YamlTlsSecret{
				Metadata: yaml_types.Metadata{
					Name: form.Common.ClusterName + "-tls-keypair",
					Namespace: auth.GetKubernetesNamespace(),
				},
				Kind: "Secret",
				APIVersion: "v1",
				Data: yaml_types.TlsData{
					TLSCrtBase64: form.TLS.TLSModeFromFile.TlsCertificateBase64,
					TLSKeyBase64: form.TLS.TLSModeFromFile.TlsPrivateKeyBase64,
				},
				Type: "Opaque",
			}

			yamlTemplate.Spec.Tls.CaSecret = caSecret.Metadata.Name
			yamlTemplate.Spec.Tls.TlsSecret = tlsSecret.Metadata.Name

			log.Trace(PostgresProviderLogPrefix + "Form to yaml conversion done")
			return []interface{}{caSecret, tlsSecret, yamlTemplate}, nil

		} else if form.TLS.TLSMode == "TlsFromSecret" {
			// Use existing secret
			yamlTemplate.Spec.Tls.CaSecret = form.TLS.TLSModeFromSecret.CaSecret
			yamlTemplate.Spec.Tls.TlsSecret = form.TLS.TLSModeFromSecret.TLSSecret
		}

	}

	log.Trace(PostgresProviderLogPrefix + "Form to yaml conversion done")
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

	err = api.Apply([]byte(yaml))
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
