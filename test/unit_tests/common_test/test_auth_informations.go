package common_test

type TestUser struct {
	KubernetesAccessToken string `json:"kubernetes_access_token"`
	KubernetesNamespace   string `json:"kubernetes_namespace"`
}

// Interface IKubernetesAuthInformation
func (user TestUser) GetKubernetesAccessToken() string {
	return user.KubernetesAccessToken
}

// Interface IKubernetesAuthInformation
func (user TestUser) GetKubernetesNamespace() string {
	return user.KubernetesNamespace
}
