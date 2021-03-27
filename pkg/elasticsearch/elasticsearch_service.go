package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/elasticsearch/actions"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/utils/provider"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
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

func (es ElasticSearchService) GetActionGroups() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "User",
			Actions: []action.IAction{
				actions.GetCredentialsAction(&es.ElasticsearchServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "Security",
			Actions: []action.IAction{
				actions.SetCertificateAction(&es.ElasticsearchServiceInformations),
				actions.CreateGetExposeInformationAction(&es.ElasticsearchServiceInformations),
				actions.CreateExposeToggleAction(&es.ElasticsearchServiceInformations),
			},
		},
	}
}
