package actions

import (
	"OperatorAutomation/pkg/core/action"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	msgs "github.com/crunchydata/postgres-operator/pkg/apiservermsgs"
)

// Creates an action to deliver port informations about the service
func CreateBackupAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Backup",
		UniqueCommand: "cmd_pg_create_backup",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, CreateBackup(service)
		},
	}
}

// Creates an action to deliver port informations about the service
func RestoreBackupAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Restore",
		UniqueCommand: "cmd_pg_restore_backup",
		Placeholder:   &action_dtos.AddUserDto{Password: generateCryptoRandomPassword(16)},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, AddUser(placeholder.(*action_dtos.AddUserDto), service)
		},
	}
}

// Adds an user to the database
func CreateBackup(service *pgCommon.PostgresServiceInformations) error {
	_, err := service.PgoApi.CreateBackup(&msgs.CreateBackrestBackupRequest{
		Namespace: service.ClusterInstance.Namespace,
		Selector: service.ClusterInstance.Name,
	})

	if err != nil {
		logger.RError(err, "Unable to progress create backup action for database " + service.ClusterInstance.Name)
	}

	return err
}

// Removes an user from the database
func RestoreBackup(dto *action_dtos.DeleteUserDto, service *pgCommon.PostgresServiceInformations) error {
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