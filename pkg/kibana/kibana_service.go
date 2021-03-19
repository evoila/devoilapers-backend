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
	Host         string
	K8sApi       *kubernetes.K8sApi
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
		action.ActionGroup{
			Name: "Expose",
			Actions: []action.IAction{
				action.Action{
					Name:          "Expose Via Ingress",
					UniqueCommand: "cmd_expose_action",
					Placeholder:   nil,
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.ExposeService(i)
					},
				},
				action.Action{
					Name:          "Hide",
					UniqueCommand: "cmd_hide_action",
					Placeholder:   nil,
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.HideExposedService(i)
					},
				},
			},
		},
	}
}

func (kb KibanaService) GetName() string {
	return kb.crdInstance.Name
}

func (kb KibanaService) GetType() string {
	return kb.crdInstance.Kind
}

func (kb KibanaService) GetTemplate() service.IServiceTemplate {
	return service.ServiceTemplate{
		Yaml:              kb.Yaml,
		ImportantSections: kb.ImportantSections,
	}
}

// Set certificate to the kibana service
// The CertificateDto certDto contains base64 strings
func (kb KibanaService) SetCertificateToService(certDto *dtos.CertificateDto) (interface{}, error) {
	kibanaInstance := kb.crdInstance
	certDto, err := certDto.EncodeFromBase64ToString()
	if err != nil {
		return nil, err
	}
	tlsCert := map[string][]byte{
		"ca.crt":  []byte(certDto.CaCrt),
		"tls.crt": []byte(certDto.TlsCrt),
		"tls.key": []byte(certDto.TlsKey),
	}
	if secretName, err := kb.K8sApi.CreateTlsSecret(kibanaInstance.Namespace, kibanaInstance.Name, "Kibana", GroupName+"/"+GroupVersion, string(kibanaInstance.UID), tlsCert); err != nil {
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

// ExecuteExposeAction exposes a service through ingress and return error if not successful
func (kb KibanaService) ExposeService(_ interface{}) (interface{}, error) {
	namespace := kb.crdInstance.Namespace

	// Default http port of kibana
	const port int32 = 5601

	// In a namespace, we rule to have only a ingress with convention name: <namespace>-ingress
	return kb.K8sApi.AddServiceToIngress(namespace, namespace+"-ingress", kb.Name+"-kb-http", kb.Host, port)
}

func (kb KibanaService) HideExposedService(_ interface{}) (interface{}, error) {
	namespace := kb.crdInstance.Namespace

	// Check whether ingress exists
	if _, err := kb.K8sApi.GetIngress(namespace, namespace+"-ingress"); err == nil {
		// Delete service from ingress with name convention: <namespace>-ingress
		return nil, kb.K8sApi.DeleteServiceFromIngress(namespace, namespace+"-ingress", kb.Name+"-kb-http")
	} else {
		return nil, nil
	}
}
