package provider

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
)

// Service base which implements service.IService interface
// Provides simple storage based information delivery
type BasicService struct {
	ProviderType   string
	Name   string
	Yaml string
	ImportantSections []string
	Status int
}

// Returns the service type. Part of service.IService interface
func (cs BasicService) GetType() string {
	return cs.ProviderType
}

// Returns the service name. Part of service.IService interface
func (cs BasicService) GetName() string {
	return cs.Name
}

// Returns the action groups. Part of service.IService interface
func (cs BasicService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{}
}

// Returns the service Template. Part of service.IService interface
func (cs BasicService) GetTemplate() service.IServiceTemplate {
	var template service.IServiceTemplate = service.ServiceTemplate{
		ImportantSections: cs.ImportantSections,
		Yaml: cs.Yaml,
	}

	return template
}

// Returns the service status. Part of service.IService interface
func (cs BasicService) GetStatus() int  {
	return cs.Status
}