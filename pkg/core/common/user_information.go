package common

type IKubernetesAuthInformation interface {
	GetKubernetesAccessToken() string
	GetKubernetesNamespace() string
}

