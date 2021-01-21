package service

import "OperatorAutomation/pkg/core/common"

// Provider of an service like kibana
type IServiceProvider interface {
	// Name of the provider/service i.e. postgres
	GetName() string

	// Template to create an instance
	GetTemplate() *IServiceTemplate

	GetServiceDescription() string

	GetServiceImage() string

	GetServices(auth common.IKubernetesAuthInformation) []*IService

	GetService(auth common.IKubernetesAuthInformation, id string) *IService

	CreateService(auth common.IKubernetesAuthInformation, yaml string) error

	DeleteService(auth common.IKubernetesAuthInformation, id string) error
}