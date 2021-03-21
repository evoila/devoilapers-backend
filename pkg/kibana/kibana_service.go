package kibana

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kibana/dtos"
	"OperatorAutomation/pkg/kubernetes"
	"OperatorAutomation/pkg/utils/provider"
	"strings"

	commonV1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1"
)

type KibanaService struct {
	status commonV1.DeploymentHealth
	provider.BasicService
	api    *kubernetes.K8sApi
	crdApi *kubernetes.CommonCrdApi
	auth   common.IKubernetesAuthInformation
	host   string
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
			Name: "Elasticsearch Action Group",
			Actions: []action.IAction{
				action.Action{
					Name:          "ExposeThroughIngress",
					UniqueCommand: "expose",
					Placeholder:   &dtos.ExposeInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.ExecuteExposeAction(i.(*dtos.ExposeInformation))
					},
				},
				action.Action{
					Name:          "UnexposeThroughIngress",
					UniqueCommand: "unexpose",
					Placeholder:   &dtos.ExposeInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.ExecuteUnexposeAction(i.(*dtos.ExposeInformation))
					},
				},
				action.Action{
					Name:          "UpdateReplicasCount",
					UniqueCommand: "rescale",
					Placeholder:   &dtos.ScaleInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return kb.ExecuteRescaleAction(i.(*dtos.ScaleInformation))
					},
				},
			},
		},
	}
}

// ExecuteExposeAction exposes a service through ingress and return error if not successful
func (kb *KibanaService) ExecuteExposeAction(dto *dtos.ExposeInformation) (interface{}, error) {
	address := strings.Split(kb.host, "/")
	host := strings.Split(address[2], ":")
	return kb.api.AddServiceToIngress(kb.auth.GetKubernetesNamespace(), dto.IngressName, kb.Name+"-kb-http", host[0], 5601)
}
func (kb *KibanaService) ExecuteUnexposeAction(dto *dtos.ExposeInformation) (interface{}, error) {

	return nil, kb.api.DeleteServiceFromIngress(kb.auth.GetKubernetesNamespace(), dto.IngressName, kb.Name+"-kb-http")
}

// ExecuteRescaleAction rescales a kb-deployment and return error if not successful
func (kb *KibanaService) ExecuteRescaleAction(dto *dtos.ScaleInformation) (interface{}, error) {

	instance := v1.Kibana{}
	kb.crdApi.Get(kb.auth.GetKubernetesNamespace(), kb.Name, ResourceName, &instance)
	name := kb.Name + "-kb"

	return kb.api.UpdateScaleDeployment(kb.auth.GetKubernetesNamespace(), name, dto.ReplicasCount)
}
