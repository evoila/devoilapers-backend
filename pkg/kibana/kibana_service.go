package kibana

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"

	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kibana/dtos"

	commonV1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KibanaService struct {
	k8sApi       *kubernetes.K8sApi
	crdInstance  *v1.Kibana
	commonCrdApi *kubernetes.CommonCrdApi
	status       commonV1.DeploymentHealth
	provider.BasicService
}

func (kb KibanaService) GetStatus() int {
	if kb.status == commonV1.GreenHealth {
		return service.ServiceStatusOk
	} else if kb.status == commonV1.RedHealth {
		return service.ServiceStatusError
	}

	return service.ServiceStatusPending
}

func (kb KibanaService) SetCertificateToService(certDto *dtos.CertificateDto) (interface{}, error) {
	kibanaInstance := kb.crdInstance
	tlsCert := map[string][]byte{
		"ca.crt":  []byte(certDto.CaCrt),
		"tls.crt": []byte(certDto.TlsCrt),
		"tls.key": []byte(certDto.TlsKey),
	}
	if secretName, err := kb.k8sApi.CreateTlsSecret(kibanaInstance.Namespace, kibanaInstance.Name, "Kibana", GroupName+"/"+GroupVersion, string(kibanaInstance.UID), tlsCert); err != nil {
		return nil, err
	} else {
		kibanaInstance.Spec.HTTP.TLS.Certificate.SecretName = secretName
		kibanaInstance.ObjectMeta = metav1.ObjectMeta{
			Name:            kibanaInstance.Name,
			Namespace:       kibanaInstance.Namespace,
			ResourceVersion: kibanaInstance.ResourceVersion,
		}
		return nil, kb.commonCrdApi.Update(kibanaInstance.Namespace, kibanaInstance.Name, ResourceName, kibanaInstance)
	}
}

func (kb KibanaService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Secure",
			Actions: []action.IAction{
				action.Action{
					Name:          "Set Certificate",
					UniqueCommand: "cmd_set_cert_action",
					Placeholder:   &dtos.CertificateDto{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.SetCertificateToService(i.(*dtos.CertificateDto))
					},
				},
			},
		},
	}
}
