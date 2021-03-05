package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/postgres/actions/dtos"
	pgCommon "OperatorAutomation/pkg/postgres/common"
)

func CreateScaleAction(service *pgCommon.PostgresServiceInformations) action.IAction {
	return action.Action{
		Name: "Scale",
		UniqueCommand: "cmd_pg_scale",
		Placeholder: &dtos.ClusterScaleDto{},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return scaleCluster(service, placeholder.(*dtos.ClusterScaleDto))
		},
	}
}

func scaleCluster(service *pgCommon.PostgresServiceInformations, dto *dtos.ClusterScaleDto) (interface{}, error)  {
	return nil, nil
}
