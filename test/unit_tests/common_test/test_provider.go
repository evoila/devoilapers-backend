package common_test

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
)

// Do not initialize use provided CreateTestProvider function
type TestProvider struct {
	GetDescriptionCb          func() string
	GetServiceImageCb         func() string
	GetServiceTypeCb          func() string
	GetTemplateCb             func() *service.IServiceTemplate
	GetServicesCb             func(auth common.IKubernetesAuthInformation) ([]*service.IService, error)
	GetServiceCb              func(auth common.IKubernetesAuthInformation, id string) (*service.IService, error)
	CreateServiceCb           func(auth common.IKubernetesAuthInformation, yaml string) error
	DeleteServiceCb           func(auth common.IKubernetesAuthInformation, id string) error
	SetCertificateToServiceCb func(auth common.IKubernetesAuthInformation, id string, tlsCert map[string][]byte) error
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

func (es TestProvider) GetTemplate(auth common.IKubernetesAuthInformation) *service.IServiceTemplate {
	return es.GetTemplateCb()
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

func (es TestProvider) SetCertificateToService(auth common.IKubernetesAuthInformation, id string, tlsCert map[string][]byte) error {
	return es.SetCertificateToServiceCb(auth, id, tlsCert)
}
