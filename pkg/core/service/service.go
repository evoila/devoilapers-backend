package service

import "OperatorAutomation/pkg/core/action"

// Defines an abstraction for an service-instance
type IService interface {
	// Get actions provided by this service
	GetActions() []action.IActionGroup
	// Get the template on which the service depends
	GetTemplate() IServiceTemplate
	// Get IService Status
    GetStatus() int

	GetType() string

	GetName() string
}