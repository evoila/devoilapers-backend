package service

import "OperatorAutomation/pkg/core/action"

var SERVICE_STATUS_OK = 0
var SERVICE_STATUS_WARNING = 1
var SERVICE_STATUS_ERROR = 2

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