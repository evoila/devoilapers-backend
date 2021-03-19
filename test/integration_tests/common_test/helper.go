package common_test

import (
	"OperatorAutomation/pkg/core/common"
	"OperatorAutomation/pkg/core/provider"
	"OperatorAutomation/pkg/core/service"
	"errors"
	"time"
)

func WaitForServiceComeUp(provider provider.IServiceProvider, user common.IKubernetesAuthInformation, serviceId string) (*service.IService, error)  {
	for i := 0; i < 100; i++ {
		time.Sleep(5 * time.Second)

		// Try get service with invalid user data
		servicePtr, err := provider.GetService(user, serviceId)
		if err != nil {
			return nil, err
		}

		if (*servicePtr).GetStatus() == service.ServiceStatusOk {
			return servicePtr, nil
		}
	}

	return nil, errors.New("Service with id " + serviceId + " does not become ready.")
}
