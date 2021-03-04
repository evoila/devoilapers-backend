package kibana

import (
	"OperatorAutomation/pkg/core/service"
	"OperatorAutomation/pkg/utils/provider"
	commonV1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1"
)

type KibanaService struct {
	status commonV1.DeploymentHealth
	provider.BasicService
}

func (kb KibanaService) GetStatus() int {
	if kb.status == commonV1.GreenHealth {
		return service.ServiceStatusOk
	} else if kb.status == commonV1.RedHealth {
		return service.ServiceStatusError
	}

	return service.ServiceStatusPending
}
