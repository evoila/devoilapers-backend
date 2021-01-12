package kubernetes

import (
	"OperatorAutomation/pkg/core/common"
)

type KubernetesWrapper struct {
}

func CreateKubernetesWrapper(userInformation common.IKubernetesAuthInformation) *KubernetesWrapper {
	return &KubernetesWrapper{}
}
