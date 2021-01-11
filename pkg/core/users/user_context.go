package users

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kubernetes"
)

// Created by core instance. Do not create by yourself
type UserContext struct {
	common.IKubernetesAuthInformation
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