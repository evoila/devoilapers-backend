package service

import "OperatorAutomation/pkg/core/common"

// Provider of an service like kibana
type IServiceProvider interface {
	// Get Type of the service i.e. postgres
	GetServiceType() string

	// Template to create an instance
	GetTemplate(auth common.IKubernetesAuthInformation) *IServiceTemplate

	GetServiceDescription() string

	GetServiceImage() string

	GetServices(auth common.IKubernetesAuthInformation) ([]*IService, error)

	GetService(auth common.IKubernetesAuthInformation, id string) (*IService, error)

	CreateService(auth common.IKubernetesAuthInformation, yaml string) error

	DeleteService(auth common.IKubernetesAuthInformation, id string) error
}