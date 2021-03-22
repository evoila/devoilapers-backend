package postgres

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/kubernetes/yaml_types"
	"OperatorAutomation/pkg/postgres/actions"
	common2 "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos"
	"OperatorAutomation/pkg/utils"
	"OperatorAutomation/pkg/utils/logger"
	"OperatorAutomation/pkg/utils/provider"
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"path"
	"strconv"
	"strings"
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

	logger.RTrace("Creating new postgres provider")

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
	logger.RTrace("Going to convert received form data to yaml")

	form := dtos.FormResponseDto{}
	err := json.Unmarshal(jsonFormResult, &form)
	if err != nil {
		logger.RError(err, "Could not unmarshal received form to generate the yaml")
		return nil, err
	}

	// Create form with form default values
	yamlTemplate := dtos.ProviderYamlTemplateDto{}
	err = yaml.Unmarshal([]byte(pg.YamlTemplate), &yamlTemplate)
	if err != nil {
		logger.RError(err, "Could not unmarshal the default postgres yaml template")
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
		logger.RTrace("Tls is requested during yaml conversion")

		if form.TLS.TLSMode == "TlsFromFile" {
			logger.RTrace("Use tls mode from file. Going to construct secrets")

			// Construct secrets accordingly to
			// https://access.crunchydata.com/documentation/postgres-operator/4.6.1/tutorial/tls/
			caSecret := &yaml_types.YamlCaSecret{
				Metadata: yaml_types.Metadata{
					Name:      form.Common.ClusterName + "-ca",
					Namespace: auth.GetKubernetesNamespace(),
				},
				Kind:       "Secret",
				APIVersion: "v1",
				Data: yaml_types.CaData{
					CaCrtBase64: form.TLS.TLSModeFromFile.CaCertBase64,
				},
				Type: "Opaque",
			}

			tlsSecret := &yaml_types.YamlTlsSecret{
				Metadata: yaml_types.Metadata{
					Name:      form.Common.ClusterName + "-tls-keypair",
					Namespace: auth.GetKubernetesNamespace(),
				},
				Kind:       "Secret",
				APIVersion: "v1",
				Data: yaml_types.TlsData{
					TLSCrtBase64: form.TLS.TLSModeFromFile.TlsCertificateBase64,
					TLSKeyBase64: form.TLS.TLSModeFromFile.TlsPrivateKeyBase64,
				},
				Type: "Opaque",
			}

			yamlTemplate.Spec.Tls.CaSecret = caSecret.Metadata.Name
			yamlTemplate.Spec.Tls.TlsSecret = tlsSecret.Metadata.Name

			logger.RTrace("Form to yaml conversion done")
			return []interface{}{caSecret, tlsSecret, yamlTemplate}, nil

		} else if form.TLS.TLSMode == "TlsFromSecret" {
			logger.RTrace("Use tls mode from secret. Going to construct secrets")
			// Use existing secret
			yamlTemplate.Spec.Tls.CaSecret = form.TLS.TLSModeFromSecret.CaSecret
			yamlTemplate.Spec.Tls.TlsSecret = form.TLS.TLSModeFromSecret.TLSSecret
		}
	}

	logger.RTrace("Form to yaml conversion done")
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
	logger.RInfo("Create new postgres service by yaml")

	api, err := kubernetes.GenerateK8sApiFromToken(pg.Host, pg.CaPath, auth.GetKubernetesAccessToken())
	if err != nil {
		logger.RError(err, "Could not generate kubernetes api from token")
		return err
	}

	// Create the objects
	createdObjects, err := api.Apply([]byte(yaml))
	if err != nil {
		logger.RError(err, "Could not apply the yaml")
		return err
	}


	// Since we may have created secrets we have to concat the ownership of the new cluster with
	// the new secrets to ensure that they get deleted if the cluster is deleted.
	// Unfortunately we can not enforce uids, therefore we have to create the objects
	// first and add the ownership later on.

	newClusterName := "";
	newClusterUid := "";
	var newSecretNames []string

	// Loop the new created objects
	for _, createdObject := range createdObjects {
		// Identify the cluster and export the necessary informations
		if createdObject.Object["kind"].(string) == "Pgcluster" {
			metadata := createdObject.Object["metadata"].(map[string]interface{})
			newClusterName = metadata["name"].(string)
			newClusterUid = metadata["uid"].(string)

			logger.RTrace("Identified new cluster " + newClusterName)
		}

		// Identify the secrets and export all
		if createdObject.Object["kind"].(string) == "Secret" {
			metadata := createdObject.Object["metadata"].(map[string]interface{})
			newSecretName := metadata["name"].(string)
			newSecretNames = append(newSecretNames, newSecretName)

			logger.RTrace("Identified new secret " + newSecretName)
		}
	}

	// Construct a patch for the secrets based on the name and the uid that we extracted bevorhead
	logger.RTrace("Preparing ownership patch on new secrets")
	patch := fmt.Sprintf(`
			{
			  "metadata": {
				"ownerReferences": [
				  {
					"apiVersion": "crunchydata.com/v1",
					"controller": true,
					"blockOwnerDeletion": true,
					"kind": "Pgcluster",
					"name": "%s",
					"uid": "%s"
				  }
				]
			  }
			}
	`, newClusterName, newClusterUid)

	// Loop every secret to patch it with the owner reference from above.
	for _, secretName := range newSecretNames {
		logger.RInfo("Apply ownership of cluster " + newClusterName+ " to secret " + secretName)

		// Patch the secret with the owner reference
		_, err = api.ClientSet.CoreV1().Secrets(auth.GetKubernetesNamespace()).Patch(
			context.TODO(),
			secretName,
			types.StrategicMergePatchType,
			[]byte(patch),
			metav1.PatchOptions{})

		if err != nil {
			return err
		}
	}

	logger.RInfo("Postgres service created")
	return nil
}

func (pg PostgresProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	logger.RInfo("Delete postgres service with id " + id)

	postgresCrd, err := pg.createCrdApi(auth)
	if err != nil {
		logger.RError(err, "Service could not be deleted because the crd api could not be created")
		return err
	}

	// Find service which should be deleted
	serviceToDeletePtr, err := pg.GetService(auth, id)
	if err != nil {
		logger.RError(err, "Kubernetes api could not be created")
		return err
	}
	serviceToDelete := (*serviceToDeletePtr).(PostgresService)

	// Revoke the exposure if there is an exposure
	err = actions.Hide(&serviceToDelete.PostgresServiceInformations)
	if err != nil && err.Error() != "service is not exposed" {
		logger.RError(err, "Service could not be deleted because the service is exposed and could not be revoked")
		return err
	}

	// Finally delete the cluster
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
		logger.RError(err, "Could not marshal the kubernetes service struct")
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
