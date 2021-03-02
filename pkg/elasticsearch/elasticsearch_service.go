package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

type ElasticSearchService struct {
	providerType   string
	name   string
	status v1.ElasticsearchHealth
	yaml string
	importantSections []string
}

func (es ElasticSearchService) GetType() string {
	return es.providerType
}

func (es ElasticSearchService) GetName() string {
	return es.name
}

func (es ElasticSearchService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{}
}

func (es ElasticSearchService) GetTemplate() service.IServiceTemplate {
	var template service.IServiceTemplate = service.ServiceTemplate{
		ImportantSections: es.importantSections,
		Yaml: es.yaml,
	}

	return template
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
