package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/actions"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/provider"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ElasticSearchService struct {
	status v1.ElasticsearchHealth
	provider.BasicService
	esCommon.ElasticsearchServiceInformations
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

// Set certificate to elastic search service
// The CertificateDto certDto contains base64 strings
func (es ElasticSearchService) SetCertificateToService(certDto *action_dtos.CertificateDto) (interface{}, error) {
	elasticInstance := es.ClusterInstance
	certDto, err := certDto.EncodeFromBase64ToString()

	if err != nil {
		return nil, err
	}

	tlsCert := map[string][]byte{
		"ca.crt":  []byte(certDto.CaCrt),
		"tls.crt": []byte(certDto.TlsCrt),
		"tls.key": []byte(certDto.TlsKey),
	}

	secretName, err := es.K8sApi.CreateTlsSecret(
		elasticInstance.Namespace,
		elasticInstance.Name,
		"Elasticsearch",
		GroupName+"/"+GroupVersion,
		string(elasticInstance.UID),
		tlsCert)

	if err != nil {
		return nil, err
	}

	elasticInstance.Spec.HTTP.TLS.Certificate.SecretName = secretName
	elasticInstance.ObjectMeta = metav1.ObjectMeta{
		Name:            elasticInstance.Name,
		Namespace:       elasticInstance.Namespace,
		ResourceVersion: elasticInstance.ResourceVersion,
	}

	return nil, es.CrdClient.Update(elasticInstance.Namespace, elasticInstance.Name, RessourceName, elasticInstance)
}

func (es ElasticSearchService) GetActionGroups() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Secure",
			Actions: []action.IAction{
				action.FormAction{
					Name:          "Set Certificate",
					UniqueCommand: "cmd_es_set_cert_action",
					Placeholder:   &action_dtos.CertificateDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.SetCertificateToService(i.(*action_dtos.CertificateDto))
					},
				},
			},
		},
		action.ActionGroup{
			Name: "Exposure",
			Actions: []action.IAction{
				actions.CreateGetExposeInformationAction(&es.ElasticsearchServiceInformations),
				actions.CreateExposeAction(&es.ElasticsearchServiceInformations),
				actions.DeleteExposeAction(&es.ElasticsearchServiceInformations),
			},
		},
	}
}
