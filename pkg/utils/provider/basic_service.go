package provider

import (
	"OperatorAutomation/pkg/core/action"
)

// Service base which implements service.IService interface
// Provides simple storage based information delivery
type BasicService struct {
	ProviderType   string
	Name   string
	Yaml string
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
func (cs BasicService) GetYamlTemplate() string {
	return cs.Yaml
}

// Returns the service status. Part of service.IService interface
func (cs BasicService) GetStatus() int  {
	return cs.Status
}