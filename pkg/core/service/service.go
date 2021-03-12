package service

import "OperatorAutomation/pkg/core/action"

const ServiceStatusOk = 0
const ServiceStatusWarning = 1
const ServiceStatusError = 2
const ServiceStatusPending = 3

// Defines an abstraction for an service-instance
type IService interface {
	// Get actions provided by this service
	GetActions() []action.IActionGroup
	// Get the template on which the service depends
	GetYamlTemplate() string
	// Get service Status. See defined consts.
    GetStatus() int
	// Get type of service i.e. postgres/elasticsearch
	GetType() string
	// Name of the deployed service
	GetName() string
}