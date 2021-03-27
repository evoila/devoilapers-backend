package actions

import (
	"OperatorAutomation/pkg/core/action"
	"OperatorAutomation/pkg/kibana/dtos/action_dtos"
	kbCommon "OperatorAutomation/pkg/kibana/common"
)

func ScaleClusterAction(service *kbCommon.KibanaServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Scale",
		UniqueCommand: "cmd_kb_scale",
		Placeholder: &action_dtos.ClusterScaleDto{
			NumberOfReplicas:  int(service.ClusterInstance.Spec.Count),
		},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, ScaleCluster(service, placeholder.(*action_dtos.ClusterScaleDto))
		},
	}
}

func ScaleCluster(service  *kbCommon.KibanaServiceInformations, dto *action_dtos.ClusterScaleDto) error {
	service.ClusterInstance.Spec.Count = int32(dto.NumberOfReplicas)
	return service.CrdClient.Update(service.ClusterInstance.Namespace, service.ClusterInstance.Name, kbCommon.ResourceName, service.ClusterInstance)
}
