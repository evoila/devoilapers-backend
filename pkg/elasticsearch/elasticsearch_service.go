package elasticsearch

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
)

type ElasticSearchService struct {
	providerType   string
	name   string
	status string
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
	if es.status == "green" {
		return service.SERVICE_STATUS_OK
	} else if es.status == "yellow" {
		return service.SERVICE_STATUS_WARNING
	} else if es.status == "red" {
		return service.SERVICE_STATUS_ERROR
	}

	return service.SERVICE_STATUS_PENDING
}
