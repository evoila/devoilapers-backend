package actions

import (
	"OperatorAutomation/pkg/core/action"
	esCommon "OperatorAutomation/pkg/elasticsearch/common"
	"OperatorAutomation/pkg/elasticsearch/dtos/action_dtos"
)

func ScaleClusterAction(service *esCommon.ElasticsearchServiceInformations) action.IAction {
	return action.FormAction{
		Name:          "Scale",
		UniqueCommand: "cmd_es_scale",
		Placeholder: &action_dtos.ClusterScaleDto{
			NumberOfReplicas: scalePlaceholder(service),
		},
		ActionExecuteCallback: func(placeholder interface{}) (interface{}, error) {
			return nil, ScaleCluster(service, placeholder.(*action_dtos.ClusterScaleDto))
		},
	}
}

func scalePlaceholder(service *esCommon.ElasticsearchServiceInformations) int {
	if len(service.ClusterInstance.Spec.NodeSets) == 0 {
		return 0
	}

	return int(service.ClusterInstance.Spec.NodeSets[0].Count)
}

func ScaleCluster(service *esCommon.ElasticsearchServiceInformations, dto *action_dtos.ClusterScaleDto) error {
	service.ClusterInstance.Spec.NodeSets[0].Count = int32(dto.NumberOfReplicas)
	return service.CrdClient.Update(service.ClusterInstance.Namespace, service.ClusterInstance.Name, esCommon.RessourceName, service.ClusterInstance)
}
