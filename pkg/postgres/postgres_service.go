package postgres

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/postgres/actions"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/utils/provider"
	v1 "github.com/crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

type PostgresService struct {
	provider.BasicService
	pgCommon.PostgresServiceInformations
}

// Returns the action groups. Part of service.IService interface
func (pg PostgresService) GetActionGroups() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Features",
			Actions: []action.IAction{
				actions.CreateScaleAction(&pg.PostgresServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "User",
			Actions: []action.IAction{
				actions.ShowUserAction(&pg.PostgresServiceInformations),
				actions.CreateUserAction(&pg.PostgresServiceInformations),
				actions.DeleteUserAction(&pg.PostgresServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "Security",
			Actions: []action.IAction{
				actions.CreateGetExposeInformationAction(&pg.PostgresServiceInformations),
				actions.CreateExposeAction(&pg.PostgresServiceInformations),
				actions.DeleteExposeAction(&pg.PostgresServiceInformations),
			},
		},
	}
}

func (pg PostgresService) GetStatus() int {
	status := pg.ClusterInstance.Status.State
	if status == v1.PgclusterStateProcessed ||
		status == v1.PgclusterStateBootstrapping {
		return service.ServiceStatusPending
	} else if status == v1.PgclusterStateInitialized ||
		status == v1.PgclusterStateCreated ||
		status == v1.PgclusterStateBootstrapped {
		return service.ServiceStatusOk
	}

	return service.ServiceStatusError
}
