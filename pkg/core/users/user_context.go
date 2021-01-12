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

// Shut return the service instance with the given id.
// The service instance should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetService(id string) *service.Service {
	// TODO
	return nil
}

// Shut return the all service instances.
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetServices() []*service.Service {
	// TODO
	return nil
}

// Shut delete the given service instance from the cluster
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) DeleteService(service *service.Service) error {
	// TODO
	return nil
}

// Shout create all necessary stuff in the cluster from the given service template.
// Ensure that the type can be identified later on.
func (ctx UserContext) CreateServices(serviceTemplate *service.ServiceTemplate) error {
	// TODO
	return nil
}
