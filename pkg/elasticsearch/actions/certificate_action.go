package actions

import (
	"OperatorAutomation/pkg/core/action"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to deliver port informations about the service
func SetCertificateAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Set Certificate",
		UniqueCommand: "cmd_es_set_cert_action",
		Placeholder:   &action_dtos.CertificateDto{},
		ActionExecuteCallback: func(i interface{}) (interface{}, error) {
			return SetCertificateToService(i.(*action_dtos.CertificateDto), service)
		},
	}
}


// Set certificate to elastic search service
// The CertificateDto certDto contains base64 strings
func SetCertificateToService(certDto *action_dtos.CertificateDto, es *esCommon.ElasticsearchServiceInformations) (interface{}, error) {
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
		esCommon.GroupName+"/"+esCommon.GroupVersion,
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

	return nil, es.CrdClient.Update(elasticInstance.Namespace, elasticInstance.Name, esCommon.RessourceName, elasticInstance)
}
