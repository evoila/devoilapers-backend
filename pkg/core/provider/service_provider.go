package provider

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/service"
)

// Provider of an service like kibana
type IServiceProvider interface {
	// Will be called by the core if it is initialized
	OnCoreInitialized(provider []*IServiceProvider)

	// Get Type of the service i.e. postgres
	GetServiceType() string

	// Template to create an instance
	GetYamlTemplate(auth common.IKubernetesAuthInformation, jsonFormResult []byte) (interface{}, error)

	// Delivers a json schema to generate the required form to satisfy the yaml generation
	GetJsonForm(auth common.IKubernetesAuthInformation) (interface{}, error)

	// Get a short description of this service type
	GetServiceDescription() string

	// Get an base64 or url to an image
	GetServiceImage() string

	// Get services of this type from the kubernetes cluster
	GetServices(auth common.IKubernetesAuthInformation) ([]*service.IService, error)

	// Get a specific service of this type from the kubernetes cluster
	GetService(auth common.IKubernetesAuthInformation, id string) (*service.IService, error)

	// Create a new service
	CreateService(auth common.IKubernetesAuthInformation, yaml string) error

	// Delete the service
	DeleteService(auth common.IKubernetesAuthInformation, id string) error
}
