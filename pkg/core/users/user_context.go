package users

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
)

type IUserInformation interface {
	GetName() string
	GetPassword() string
	GetKubernetesAccessToken() string
	GetKubernetesNamespace() string
	GetRole() string
}

type UserContext struct {
	IUserInformation
	*kubernetes.KubernetesWrapper
}

func (ctx UserContext) GetService(id string) *service.Service {
	return nil
}

func (ctx UserContext) GetServices() []*service.Service {
	return nil
}

func (ctx UserContext) DeleteService(service *service.Service) error {
	return nil
}

func (ctx UserContext) CreateServices(serviceTemplate *service.ServiceTemplate) error {
	return nil
}