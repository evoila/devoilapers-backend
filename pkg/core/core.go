package core

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/users"
)

// Use CreateCore to create an instance.
// Holds all references for interacting with the system from a higher level
type Core struct {
	//UserContextManagement *users.UserContextManagement
	provider.ServiceProviderRegistry
}

// Creates an instance of the core struct holding the references for
// the user-context-management and services
func CreateCore(providers []*provider.IServiceProvider) *Core {
	core := Core{provider.ServiceProviderRegistry{Providers: map[string]*provider.IServiceProvider{}}}
	for _, provider := range providers {
		providerType := (*provider).GetServiceType()

		if _, ok := core.Providers[providerType]; ok {
			panic("Duplicate provider type during core initialization")
		}

		core.Providers[providerType] = provider
	}

	for _, provider := range providers {
		(*provider).OnCoreInitialized(providers)
	}

	return &core
}

// Creates a new user context based on the given authentifcation data
func (core Core) CrateUserContext(userInformation common.IKubernetesAuthInformation) *users.UserContext {
	return &users.UserContext{ServiceProviderRegistry: &core.ServiceProviderRegistry, Auth: userInformation}
}
