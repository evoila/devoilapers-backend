package utils

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/common"
)

func ToInterface(listOfUsers []config.User) []common.IKubernetesAuthInformation {
	userInterfaceRepresentation := make([]common.IKubernetesAuthInformation, len(listOfUsers))
	for i, user := range listOfUsers {
		userInterfaceRepresentation[i] = user
	}

	return userInterfaceRepresentation
}