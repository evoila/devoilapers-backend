package actions

import (
	"OperatorAutomation/pkg/core/action"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	"crypto/rand"
	"encoding/base64"
	msgs "github.com/crunchydata/postgres-operator/pkg/apiservermsgs"
)

// Creates an action to deliver port informations about the service
func ShowUserAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Show",
		UniqueCommand: "cmd_pg_show_users",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return ShowUsers(service)
		},
	}
}

// Creates an action to deliver port informations about the service
func CreateUserAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Add",
		UniqueCommand: "cmd_pg_add_user",
		Placeholder:   &action_dtos.AddUserDto{Password: generateCryptoRandomPassword(16)},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, AddUser(placeholder.(*action_dtos.AddUserDto), service)
		},
	}
}

// Creates an action to remove the exposure
func DeleteUserAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Delete",
		UniqueCommand: "cmd_pg_remove_user",
		Placeholder:   &action_dtos.DeleteUserDto{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, DeleteUser(placeholder.(*action_dtos.DeleteUserDto), service)
		},
	}
}

// Adds an user to the database
func AddUser(dto *action_dtos.AddUserDto, service *pgCommon.PostgresServiceInformations) error {
	_, err := service.PgoApi.CreateUser(&msgs.CreateUserRequest{
		Password:dto.Password,
		Username: dto.Username,
		Namespace: service.ClusterInstance.Namespace,
		Clusters: []string{service.ClusterInstance.Name},
		PasswordType: "md5",
	})

	if err != nil {
		logger.RError(err, "Unable to progress add user action for user " + dto.Username)
	}

	return err
}

// Removes an user from the database
func DeleteUser(dto *action_dtos.DeleteUserDto, service *pgCommon.PostgresServiceInformations) error {
	_, err := service.PgoApi.DeleteUser(&msgs.DeleteUserRequest{
		Username: dto.Username,
		Namespace: service.ClusterInstance.Namespace,
		Clusters: []string{service.ClusterInstance.Name},
	})

	if err != nil {
		logger.RError(err, "Unable to progress delete user action on user " + dto.Username)
	}

	return err
}

// Shows database users
func ShowUsers(service *pgCommon.PostgresServiceInformations) (map[string]string, error) {
	response, err := service.PgoApi.GetUsers(&msgs.ShowUserRequest{
		Namespace: service.ClusterInstance.Namespace,
		Clusters: []string{service.ClusterInstance.Name},
	})

	if err != nil {
		logger.RError(err, "Unable to progress show users action")
		return nil, err
	}

	users := map[string]string{}
	for _, user := range response.Results {
		if user.Password == "" {
			users[user.Username] = "***"
		} else {
			users[user.Username] = user.Password
		}
	}

	return users, err
}

// Generates a random password
func generateCryptoRandomPassword(n int) string {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		logger.RError(err, "Unable generate a password")
		return ""
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
