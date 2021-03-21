package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/pkg/utils/provider"
	"strings"

	"OperatorAutomation/pkg/kubernetes"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

type ElasticSearchService struct {
	status v1.ElasticsearchHealth
	provider.BasicService
	api    *kubernetes.K8sApi
	crdApi *kubernetes.CommonCrdApi
	auth   common.IKubernetesAuthInformation
	host   string
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
			Name: "Elasticsearch Action Group",
			Actions: []action.IAction{
				action.Action{
					Name:          "ExposeThroughIngress",
					UniqueCommand: "expose",
					Placeholder:   &dtos.ExposeInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.ExecuteExposeAction(i.(*dtos.ExposeInformation))
					},
				},
				action.Action{
					Name:          "UnexposeThroughIngress",
					UniqueCommand: "unexpose",
					Placeholder:   &dtos.ExposeInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.ExecuteUnexposeAction(i.(*dtos.ExposeInformation))
					},
				},
				action.Action{
					Name:          "UpdateReplicasCount",
					UniqueCommand: "rescale",
					Placeholder:   &dtos.ScaleInformation{},
					ActionExecuteCallback: func(i interface{}) (interface{}, error) {
						return es.ExecuteRescaleAction(i.(*dtos.ScaleInformation))
					},
				},
			},
		},
	}
}

// ExecuteExposeAction exposes a service through ingress and return error if not successful
func (es *ElasticSearchService) ExecuteExposeAction(dto *dtos.ExposeInformation) (interface{}, error) {
	address := strings.Split(es.host, "/")
	host := strings.Split(address[2], ":")
	return es.api.AddServiceToIngress(es.auth.GetKubernetesNamespace(), dto.IngressName, es.Name+"-es-http", host[0], 9200)
}

// ExecuteUnexposeAction unexposes a service through ingress and return error if not successful
func (es *ElasticSearchService) ExecuteUnexposeAction(dto *dtos.ExposeInformation) (interface{}, error) {

	return nil, es.api.DeleteServiceFromIngress(es.auth.GetKubernetesNamespace(), dto.IngressName, es.Name+"-es-http")
}

// ExecuteRescaleAction rescales es-statefulset and return error if not successful
func (es *ElasticSearchService) ExecuteRescaleAction(dto *dtos.ScaleInformation) (interface{}, error) {
	instance := v1.Elasticsearch{}
	es.crdApi.Get(es.auth.GetKubernetesNamespace(), es.Name, RessourceName, &instance)
	var err error
	for _, nodeset := range instance.Spec.NodeSets {
		name := es.Name + "-es-" + nodeset.Name
		_, err = es.api.UpdateScaleStatefulSet(es.auth.GetKubernetesNamespace(), name, dto.ReplicasCount)
		if err != nil {
			break
		}
	}

	return nil, err
}
