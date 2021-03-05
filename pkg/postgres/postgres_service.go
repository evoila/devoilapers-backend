package postgres

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/postgres/actions"
	pgCommon "OperatorAutomation/pkg/postgres/common"
	"OperatorAutomation/pkg/utils/provider"
	v1 "github.com/Crunchydata/postgres-operator/pkg/apis/crunchydata.com/v1"
)

type PostgresService struct {
	provider.BasicService
	pgCommon.PostgresServiceInformations
}

// Returns the action groups. Part of service.IService interface
func (pg PostgresService) GetActions() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "Informations",
			Actions: []action.IAction{
				actions.CreateGetCredentialsAction(&pg.PostgresServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "Features",
			Actions: []action.IAction{
				actions.CreateScaleAction(&pg.PostgresServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "Security",
			Actions: []action.IAction{
				actions.CreateExposeAction(&pg.PostgresServiceInformations),
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
	} else if status == v1.PgclusterStateBootstrapping {

	}

	return service.ServiceStatusError
}
