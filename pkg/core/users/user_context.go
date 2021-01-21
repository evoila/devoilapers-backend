package users

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
)

// Created by core instance. Do not create by yourself
type UserContext struct {
	serviceProviderRegistry *service.ServiceProviderRegistry
	auth common.IKubernetesAuthInformation
}

// Shut return the service instance with the given id.
// The service instance should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetService(serviceType string, id string) (*service.IService,error) {
	provider,err := ctx.serviceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return nil, err
	}
	return (*provider).GetService(ctx.auth,id),nil
}

// Shut return the all service instances.
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetServices() []*service.IService {
	var services []*service.IService
	for _,provider := range ctx.serviceProviderRegistry.Providers{
		for _,ser := range (*provider).GetServices(ctx.auth){
			services = append(services, ser)
		}
	}
	return services
}

// Shut delete the given service instance from the cluster
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) DeleteService(serviceType string, id string) error {
	provider,err := ctx.serviceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return err
	}
	return (*provider).DeleteService(ctx.auth,id)
}

// Shout create all necessary stuff in the cluster from the given service template.
// Ensure that the type can be identified later on.
func (ctx UserContext) CreateServices(serviceType string, yaml string) error {
	provider,err := ctx.serviceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return err
	}
	return (*provider).CreateService(ctx.auth,yaml)
}
