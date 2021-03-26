package kibana

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/kibana/actions"
	kbCommon "OperatorAutomation/pkg/kibana/common"
	"OperatorAutomation/pkg/utils/provider"
	"OperatorAutomation/pkg/core/action"
	commonV1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1"
)

type KibanaService struct {
	status       commonV1.DeploymentHealth
	provider.BasicService
	kbCommon.KibanaServiceInformations
}

func (kb KibanaService) GetStatus() int {
	if kb.status == commonV1.GreenHealth {
		return service.ServiceStatusOk
	} else if kb.status == commonV1.RedHealth {
		return service.ServiceStatusError
	}

	return service.ServiceStatusPending
}


func (kb KibanaService) GetActionGroups() []action.IActionGroup {
	return []action.IActionGroup{
		action.ActionGroup{
			Name: "User",
			Actions: []action.IAction{
				actions.GetCredentialsAction(&kb.KibanaServiceInformations),
			},
		},
		action.ActionGroup{
			Name: "Security",
			Actions: []action.IAction{
				actions.SetCertificateAction(&kb.KibanaServiceInformations),
				actions.CreateGetExposeInformationAction(&kb.KibanaServiceInformations),
				actions.CreateExposeToggleAction(&kb.KibanaServiceInformations),
			},
		},
	}
}
