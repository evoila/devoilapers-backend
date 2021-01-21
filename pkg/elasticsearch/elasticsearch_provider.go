package elasticsearch

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
)

// Should implement ServiceProvider interface
type ElasticSearchProvider struct {

}

func (es ElasticSearchProvider) GetServiceType() string {
	return "ElasticSearch"
}

func (es ElasticSearchProvider) GetTemplate() *service.IServiceTemplate {
	panic("implement me")
}

func (es ElasticSearchProvider) GetServices(auth common.IKubernetesAuthInformation) []*service.IService {
	panic("implement me")
}

func (es ElasticSearchProvider) GetService(auth common.IKubernetesAuthInformation, id string) *service.IService {
	var els service.IService = ElasticSearchService{auth: auth}
	return &els
}

func (es ElasticSearchProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	panic("implement me")
}

func (es ElasticSearchProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	panic("implement me")
}

