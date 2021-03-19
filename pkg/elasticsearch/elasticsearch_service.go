package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ElasticSearchService struct {
	Host         string
	K8sApi       *kubernetes.K8sApi
	crdInstance  *v1.Elasticsearch
	commonCrdApi *kubernetes.CommonCrdApi
	status       v1.ElasticsearchHealth
	provider.BasicService
}

func (es ElasticSearchService) GetStatus() int {
	if es.status == v1.ElasticsearchGreenHealth {
		return service.ServiceStatusOk
	} else if es.status == v1.ElasticsearchYellowHealth {
		return service.ServiceStatusWarning
	} else if es.status == v1.ElasticsearchRedHealth {
		return service.ServiceStatusError
	}

	return service.ServiceStatusPending
}

func (es ElasticSearchService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Secure",
			Actions: []action.IAction{
				action.Action{
					Name:          "Set Certificate",
					UniqueCommand: "cmd_set_cert_action",
					Placeholder:   &dtos.CertificateDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.SetCertificateToService(i.(*dtos.CertificateDto))
					},
				},
			},
		},
		action.ActionGroup{
			Name: "Backup",
			Actions: []action.IAction{
				action.Action{
					Name:          "Create Repository",
					UniqueCommand: "cmd_create_repo_action",
					Placeholder:   &dtos.RepoDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.CreateRepositoryForBackUp(i.(*dtos.RepoDto))
					},
				},
				action.Action{
					Name:          "Get Repository",
					UniqueCommand: "cmd_get_repo_action",
					Placeholder:   &dtos.RepoDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.GetRepository(i.(*dtos.RepoDto))
					},
				},
				action.Action{
					Name:          "Delete Repository",
					UniqueCommand: "cmd_delete_repo_action",
					Placeholder:   &dtos.RepoDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.DeleteBackupRepository(i.(*dtos.RepoDto))
					},
				},
				action.Action{
					Name:          "Create Backup",
					UniqueCommand: "cmd_create_backup_action",
					Placeholder:   &dtos.SnapshotDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.CreateBackup(i.(*dtos.SnapshotDto))
					},
				},
				action.Action{
					Name:          "Get Backup",
					UniqueCommand: "cmd_get_backup_action",
					Placeholder:   &dtos.SnapshotDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.GetBackup(i.(*dtos.SnapshotDto))
					},
				},
				action.Action{
					Name:          "Delete Backup",
					UniqueCommand: "cmd_delete_backup_action",
					Placeholder:   &dtos.SnapshotDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.DeleteBackup(i.(*dtos.SnapshotDto))
					},
				},
			},
		},
		action.ActionGroup{
			Name: "Expose",
			Actions: []action.IAction{
				action.Action{
					Name:          "Expose Via Ingress",
					UniqueCommand: "cmd_expose_action",
					Placeholder:   nil,
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.ExposeService(i)
					},
				},
				action.Action{
					Name:          "Hide",
					UniqueCommand: "cmd_hide_action",
					Placeholder:   nil,
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.HideExposedService(i)
					},
				},
			},
		},
	}
}

func (es ElasticSearchService) GetName() string {
	return es.crdInstance.Name
}

func (es ElasticSearchService) GetType() string {
	return es.crdInstance.Kind
}

func (es ElasticSearchService) GetTemplate() service.IServiceTemplate {
	return service.ServiceTemplate{
		Yaml:              es.Yaml,
		ImportantSections: es.ImportantSections,
	}
}

// Set certificate to elastic search service
// The CertificateDto certDto contains base64 strings
func (es ElasticSearchService) SetCertificateToService(certDto *dtos.CertificateDto) (interface{}, error) {
	elasticInstance := es.crdInstance
	certDto, err := certDto.EncodeFromBase64ToString()
	if err != nil {
		return nil, err
	}
	tlsCert := map[string][]byte{
		"ca.crt":  []byte(certDto.CaCrt),
		"tls.crt": []byte(certDto.TlsCrt),
		"tls.key": []byte(certDto.TlsKey),
	}
	if secretName, err := es.K8sApi.CreateTlsSecret(elasticInstance.Namespace, elasticInstance.Name, "Elasticsearch", GroupName+"/"+GroupVersion, string(elasticInstance.UID), tlsCert); err != nil {
		return nil, err
	} else {
		elasticInstance.Spec.HTTP.TLS.Certificate.SecretName = secretName
		elasticInstance.ObjectMeta = metav1.ObjectMeta{
			Name:            elasticInstance.Name,
			Namespace:       elasticInstance.Namespace,
			ResourceVersion: elasticInstance.ResourceVersion,
		}
		return nil, es.commonCrdApi.Update(elasticInstance.Namespace, elasticInstance.Name, RessourceName, elasticInstance)
	}
}

// General method to do request (curl -X <method> -u <username:password> -url <url/endpoint> -d data)
func (es ElasticSearchService) curl(endpoint, method string, data []byte) (*http.Response, error) {
	serviceName := es.crdInstance.Name
	serviceNamespace := es.crdInstance.Namespace
	url := fmt.Sprintf("https://%s/%s/%s/%s", es.Host, serviceNamespace, serviceName+"-es-http", endpoint)
	fmt.Println(url)
	req, err := http.NewRequest(
		method,
		url,
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	secret, err := es.K8sApi.GetSecret(serviceNamespace, serviceName+"-es-elastic-user")

	if err != nil {
		return nil, err
	}

	user := "elastic"
	password := secret.Data[user]

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password))))

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return client.Do(req)
}

// Create the repository for saving backup (a.k.a snapshot)
func (es ElasticSearchService) CreateRepositoryForBackUp(request *dtos.RepoDto) (interface{}, error) {
	if data, err := json.Marshal(request.Body); err != nil {
		return nil, err
	} else {
		endpoint := fmt.Sprintf("_snapshot/%s", request.Repository)
		method := "PUT"
		return es.curl(endpoint, method, data)
	}
}

// Get the repository
func (es ElasticSearchService) GetRepository(request *dtos.RepoDto) (interface{}, error) {
	endpoint := fmt.Sprintf("_snapshot/%s", request.Repository)
	method := "GET"
	return es.curl(endpoint, method, nil)
}

// Delete the repository to save the backup (a.k.a snapshot)
func (es ElasticSearchService) DeleteBackupRepository(request *dtos.RepoDto) (interface{}, error) {
	endpoint := fmt.Sprintf("_snapshot/%s", request.Repository)
	method := "DELETE"
	return es.curl(endpoint, method, nil)
}

// Create the service backup
func (es ElasticSearchService) CreateBackup(request *dtos.SnapshotDto) (interface{}, error) {
	if data, err := json.Marshal(request.Body); err != nil {
		return nil, err
	} else {
		endpoint := fmt.Sprintf("_snapshot/%s/%s", request.Repository, request.Snapshot)
		method := "POST"
		return es.curl(endpoint, method, data)
	}
}

// Get the service backup
func (es ElasticSearchService) GetBackup(request *dtos.SnapshotDto) (interface{}, error) {
	endpoint := fmt.Sprintf("_snapshot/%s/%s", request.Repository, request.Snapshot)
	method := "GET"
	return es.curl(endpoint, method, nil)
}

// Delete the service backup
func (es ElasticSearchService) DeleteBackup(request *dtos.SnapshotDto) (interface{}, error) {
	endpoint := fmt.Sprintf("_snapshot/%s/%s", request.Repository, request.Snapshot)
	method := "DELETE"
	return es.curl(endpoint, method, nil)
}

// Expose funtion exposes a service through ingress.
// Return url for exposed service if success; otherwise, return error
func (es ElasticSearchService) ExposeService(_ interface{}) (interface{}, error) {

	namespace := es.crdInstance.Namespace

	// Default http port for elasticsearch
	const port int32 = 9200

	// In a namespace, we rule to have only a ingress with convention name: <namespace>-ingress
	return es.K8sApi.AddServiceToIngress(namespace, namespace+"-ingress", es.Name+"-es-http", es.Host, port)
}

// Stop exposing the service through ingress
// Return error if not success
func (es ElasticSearchService) HideExposedService(_ interface{}) (interface{}, error) {
	namespace := es.crdInstance.Namespace

	// Check whether ingress exists
	if _, err := es.K8sApi.GetIngress(namespace, namespace+"-ingress"); err == nil {
		// Delete service from ingress with name convention: <namespace>-ingress
		return nil, es.K8sApi.DeleteServiceFromIngress(namespace, namespace+"-ingress", es.Name+"-es-http")
	} else {
		return nil, nil
	}
}
