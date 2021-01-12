package common

// Kubernetes authentication information
type IKubernetesAuthInformation interface {
	// Returns the kubernetes token
	GetKubernetesAccessToken() string
	// Returns the kubernetes namespace
	GetKubernetesNamespace() string
}
