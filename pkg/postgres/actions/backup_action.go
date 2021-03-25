package actions

import (
	"OperatorAutomation/pkg/core/action"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/postgres/dtos/action_dtos"
	"OperatorAutomation/pkg/utils/logger"
	msgs "github.com/crunchydata/postgres-operator/pkg/apiservermsgs"
	"strconv"
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
		Placeholder:   &action_dtos.RestoreBackupDto{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, RestoreBackup(placeholder.(*action_dtos.RestoreBackupDto), service)
		},
	}
}

// Creates an action to deliver port informations about the service
func ShowBackupsAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Show",
		UniqueCommand: "cmd_pg_show_backup",
		Placeholder:   nil,
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return ShowBackups(service)
		},
	}
}

// Adds an user to the database
func CreateBackup(service *pgCommon.PostgresServiceInformations) error {
	resp, err := service.PgoApi.CreateBackup(&msgs.CreateBackrestBackupRequest{
		Namespace: service.ClusterInstance.Namespace,
		Selector:  "name=" + service.ClusterInstance.Name,
	})

	if err != nil {
		logger.RError(err, "Unable to progress create backup action for database "+service.ClusterInstance.Name)
	}

	_ = resp
	return err
}

// Removes an user from the database
func RestoreBackup(dto *action_dtos.RestoreBackupDto, service *pgCommon.PostgresServiceInformations) error {
	_, err := service.PgoApi.RestoreBackup(&msgs.RestoreRequest{
		BackrestStorageType: "s3",
		Namespace:           service.ClusterInstance.Namespace,
		FromCluster:         dto.OldClusterName,
	})

	if err != nil {
		logger.RError(err, "Unable to progress restore backup action for "+service.ClusterInstance.Name)
	}

	return err
}

// Removes an user from the database
func ShowBackups(service *pgCommon.PostgresServiceInformations) (interface{}, error) {
	response, err := service.PgoApi.ShowBackups(service.ClusterInstance.Namespace, service.ClusterInstance.Name)

	if err != nil {
		logger.RError(err, "Unable to progress show backup action for "+service.ClusterInstance.Name)
		return nil, err
	}

	backups := map[string]string{}
	backup := response.Items[0]

	backups["Backup typ"] = backup.StorageType
	backups["Total backups"] = strconv.Itoa(len(backup.Info[0].Backups))
	backups["Status"] = backup.Info[0].Status.Message

	return backups, err
}
