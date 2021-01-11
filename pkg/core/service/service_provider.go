package service

type ServiceProviderRegistry struct {
	// Provdider name -> ServiceProvider
	provider map[string]ServiceProvider
}

// Add a new service provider
func (registry ServiceProviderRegistry) RegisterServiceProvider(provider ServiceProvider)  {

}

// Find an existing service provider
func (registry ServiceProviderRegistry) GetProviderByName(name string) *ServiceProvider  {
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