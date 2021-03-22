package users

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
)

// Created by core instance. Do not create by yourself
type UserContext struct {
	ServiceProviderRegistry *provider.ServiceProviderRegistry
	Auth                    common.IKubernetesAuthInformation
}

// Shut return the service instance with the given id.
// The service instance should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetService(serviceType string, id string) (*service.IService, error) {
	provider, err := ctx.ServiceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return nil, err
	}

	return (*provider).GetService(ctx.Auth, id)
}

// Shut return the all service instances.
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) GetServices() ([]*service.IService, error) {
	var services []*service.IService
	var err error = nil

	for _, provider := range ctx.ServiceProviderRegistry.Providers {
		runningServices, providerErr := (*provider).GetServices(ctx.Auth)
		if providerErr != nil {
			err = providerErr
		}

		for _, ser := range runningServices {
			services = append(services, ser)
		}
	}

	return services, err
}

// Shut delete the given service instance from the cluster
// The service instances should be a concrete implementation by i.e. elasticsearch.
// To ensure this, we probably need a factory pattern.
// Try not to hardcode if trees if possible.
func (ctx UserContext) DeleteService(serviceType string, id string) error {
	provider, err := ctx.ServiceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return err
	}
	return (*provider).DeleteService(ctx.Auth, id)
}

// Shout create all necessary stuff in the cluster from the given service template.
// Ensure that the type can be identified later on.
func (ctx UserContext) CreateServices(serviceType string, yaml string) error {
	provider, err := ctx.ServiceProviderRegistry.GetProviderByName(serviceType)
	if err != nil {
		return err
	}
	return (*provider).CreateService(ctx.Auth, yaml)
}
