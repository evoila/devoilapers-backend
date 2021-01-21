package service

import (
	"golang.org/x/crypto/openpgp/errors"
)

// Registry which holds all service Providers.
// From outside an app i.e the webserver can register the different Providers.
// Therefore the service-store knows which provider it has.
type ServiceProviderRegistry struct {
	// Provider name -> ServiceProvider
	Providers map[string]*IServiceProvider
}

// Find an existing service provider
func (registry ServiceProviderRegistry) GetProviderByName(name string) (*IServiceProvider, error) {
	provider, found := registry.Providers[name]
	if !found {
		return nil, errors.InvalidArgumentError("Provider for service type " + name + " not found")
	}
	return provider, nil
}
