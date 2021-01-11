package kubernetes

import "OperatorAutomation/pkg/core/users"

type KubernetesWrapper struct {

}

func CreateKubernetesWrapper(userInformation users.IUserInformation) *KubernetesWrapper  {
	return &KubernetesWrapper{}
}
