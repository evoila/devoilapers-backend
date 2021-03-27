package actions

import (
	"OperatorAutomation/pkg/core/action"
	kbCommon "OperatorAutomation/pkg/kibana/common"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates an action to deliver port informations about the service
func SetCertificateAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Set Certificate",
		UniqueCommand: "cmd_kb_set_cert_action",
		Placeholder:   &action_dtos.CertificateDto{},
		ActionExecuteCallback: func(i interface{}) (interface{}, error) {
			return SetCertificateToService(i.(*action_dtos.CertificateDto), service)
		},
	}
}

// Set certificate to elastic search service
// The CertificateDto certDto contains base64 strings
func SetCertificateToService(certDto *action_dtos.CertificateDto, kb *kbCommon.KibanaServiceInformations) (interface{}, error) {
	kibanaInstance := kb.ClusterInstance
	certDto, err := certDto.EncodeFromBase64ToString()

	if err != nil {
		return nil, err
	}

	tlsCert := map[string][]byte{
		"ca.crt":  []byte(certDto.CaCrt),
		"tls.crt": []byte(certDto.TlsCrt),
		"tls.key": []byte(certDto.TlsKey),
	}

	secretName, err := kb.K8sApi.CreateTlsSecret(
		kibanaInstance.Namespace,
		kibanaInstance.Name,
		"Kibana",
		kbCommon.GroupName+"/"+kbCommon.GroupVersion,
		string(kibanaInstance.UID),
		tlsCert)

	if err != nil {
		return nil, err
	}

	kibanaInstance.Spec.HTTP.TLS.Certificate.SecretName = secretName
	kibanaInstance.ObjectMeta = metav1.ObjectMeta{
		Name:            kibanaInstance.Name,
		Namespace:       kibanaInstance.Namespace,
		ResourceVersion: kibanaInstance.ResourceVersion,
	}

	return nil, kb.CrdClient.Update(kibanaInstance.Namespace, kibanaInstance.Name, kbCommon.ResourceName, kibanaInstance)

}
