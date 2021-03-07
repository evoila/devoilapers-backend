package elasticsearch

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/utils/provider"

	v1 "github.com/elastic/cloud-on-k8s/pkg/apis/elasticsearch/v1"
)

type ElasticSearchService struct {
	status v1.ElasticsearchHealth
	provider.BasicService
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
