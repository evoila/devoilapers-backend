package service

// Registry which holds all service providers.
// From outside an app i.e the webserver can register the different providers.
// Therefore the service-store knows which provider it has.
type ServiceProviderRegistry struct {
	// Provider name -> ServiceProvider
	provider map[string]ServiceProvider
}

// Add a new service provider
func (registry ServiceProviderRegistry) RegisterServiceProvider(provider ServiceProvider) {
	// TODO
}

// Find an existing service provider
func (registry ServiceProviderRegistry) GetProviderByName(name string) *ServiceProvider {
	// TODO
	return nil
}

// Provider of an service like kibana
type ServiceProvider interface {
	// Name of the provider/service i.e. postgres
	GetName() string

	// All possible actions. Format: Group-Name -> Action-Name -> Command-Name
	// i.e. Backup and Restore -> [Backup -> cmd_backup_pgdump, Restore -> cmd_restore_pgdump]
	GetActions() map[string]map[string]string

	// Template to create an instance
	GetTemplate() *ServiceTemplate
}
