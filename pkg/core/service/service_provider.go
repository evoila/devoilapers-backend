package service

import "OperatorAutomation/pkg/core/common"

// Provider of an service like kibana
type IServiceProvider interface {
	// Get Type of the service i.e. postgres
	GetServiceType() string

	// Template to create an instance
	GetTemplate() *IServiceTemplate

	GetServiceDescription() string

	GetServiceImage() string

	GetServices(auth common.IKubernetesAuthInformation) []*IService

	GetService(auth common.IKubernetesAuthInformation, id string) *IService

	CreateService(auth common.IKubernetesAuthInformation, yaml string) error

	DeleteService(auth common.IKubernetesAuthInformation, id string) error
}