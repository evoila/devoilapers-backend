package common_test

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
)

// Do not initialize use provided CreateTestProvider function
type TestProvider struct {
	GetDescriptionCb  func() string
	GetServiceImageCb func() string
	GetServiceTypeCb  func() string
	GetYamlTemplateCb func(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error)
	GetFormTemplateCb func(auth common.IKubernetesAuthInformation) (interface{}, error)
	GetServicesCb     func(auth common.IKubernetesAuthInformation) ([]*service.IService, error)
	GetServiceCb      func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error)
	CreateServiceCb   func(auth common.IKubernetesAuthInformation, yaml string) error
	DeleteServiceCb   func(auth common.IKubernetesAuthInformation, id string) error
}

func (es TestProvider) OnCoreInitialized(provider []*provider.IServiceProvider) {

}

func (es TestProvider) GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error) {
	return es.GetYamlTemplateCb(auth, jsonFormResult)
}

func (es TestProvider) GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error) {
	return es.GetFormTemplateCb(auth)
}

func (es TestProvider) GetServiceDescription() string {
	return es.GetDescriptionCb()
}

func (es TestProvider) GetServiceImage() string {
	return es.GetServiceImageCb()
}

func (es TestProvider) GetServiceType() string {
	return es.GetServiceTypeCb()
}

func (es TestProvider) GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error) {
	return es.GetServicesCb(auth)
}

func (es TestProvider) GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error) {
	return es.GetServiceCb(auth, id)
}

func (es TestProvider) CreateService(auth common.IKubernetesAuthInformation, yaml string) error {
	return es.CreateServiceCb(auth, yaml)
}

func (es TestProvider) DeleteService(auth common.IKubernetesAuthInformation, id string) error {
	return es.DeleteServiceCb(auth, id)
}
