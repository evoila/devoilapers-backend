package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/dtos"
	"OperatorAutomation/pkg/utils/provider"
	"fmt"

	"OperatorAutomation/pkg/kubernetes"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

type ElasticSearchService struct {
	status v1.ElasticsearchHealth
	provider.BasicService
	api      *kubernetes.K8sApi
	esCrdApi *kubernetes.CommonCrdApi
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
			},
		},
	}
}

// ExecuteExposeAction exposes a service through ingress and return error if not successful
func (es ElasticSearchService) ExecuteExposeAction(dto *dtos.ExposeInformation) (interface{}, error) {

	fmt.Errorf("es_provider_in namespace: ", "default"+" service:"+es.Name, "\n")
	return es.api.AddServiceToIngress("default", dto.IngressName, es.Name+"-es-http", "myhosst.com", 9200)
}
