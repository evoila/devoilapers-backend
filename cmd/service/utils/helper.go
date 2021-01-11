package utils

import (
	"OperatorAutomation/cmd/service/config"
	"OperatorAutomation/pkg/core/users"
)

func ToInterface(listOfUsers []config.User) []users.IUserInformation {
	userInterfaceRepresentation := make([]users.IUserInformation, len(listOfUsers))
	for i, user := range listOfUsers {
		userInterfaceRepresentation[i] = user
	}

	return userInterfaceRepresentation
}